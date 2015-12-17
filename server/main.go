package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/onebook/captcha"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	//	"net/url"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"strings"
)

var db *sql.DB

const DSN = "yi.kangsoo:@tcp(zabbix:3306)/dev_yi.kangsoo_service"

var cachedCaptchaQuene []cachedCaptcha

var port = flag.String("port", "", "http service port")
var stype = flag.String("stype", "", "store type: could be redis, memcache")
var saddress = flag.String("saddress", "", "store address: ip + port")
var expire = flag.Int("expire", 300, "store expire time")
var domain = flag.String("domain", "", "cookie domain")
var cache = flag.Int("cache", 0, "the number of cached captcha")

func init() {
	rand.Seed(time.Now().UnixNano())
}

type CacheValue struct {
	Key       string `json:"key"`
	Result    string `json:"result"`
	Ip        string `json:"ip"`
	CreatedAt string `json:"created_at"`
}

func (p *CacheValue) ToJson() string {
	s, err := json.Marshal(p)
	if err != nil {
		return ""
	}

	return string(s)
}

func GetDb() *sql.DB {
	if db == nil {
		_db, err := sql.Open("mysql", DSN)
		if err != nil {
			panic(err.Error())
		}
		//	defer db.Close()
		_db.SetMaxIdleConns(10)

		db = _db
		fmt.Println("-- DB_NEW: --")
	}
	fmt.Printf("-- DB_PING: %x -- \n", db.Ping)
	return db
}

func main() {
	parseArgs()

	var store captcha.Store

	if *stype == "redis" {
		store = captcha.NewRedisStore(*saddress, *expire)
	} else {
		store = captcha.NewMemcacheStore(*saddress, *expire)
	}

	_, e := store.Get("test-store-connection")
	if e != nil && e.Error() != "memcache: cache miss" {
		fmt.Printf("store error: %s \n", e.Error())
		os.Exit(1)
	}

	http.HandleFunc("/captcha", func(res http.ResponseWriter, req *http.Request) {
		p_key := req.URL.Query().Get("key")
		db := GetDb()
		if p_key != "" {
			var key string
			var secret string
			if err := db.QueryRow("SELECT `key`, `secret` FROM `auth` WHERE `key` = ?", p_key).Scan(&key, &secret); err != nil {
				res.WriteHeader(500)
				res.Write([]byte("UNKNOWN ERROR, "))
				res.Write([]byte(err.Error()))
				return
			}
		}

		uid, result, png := getCaptcha()

		cValue := CacheValue{
			Key:       p_key,
			Result:    result,
			Ip:        strings.Split(req.RemoteAddr, ":")[0],
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		//key := []byte(time.Now().Format(time.StampNano) + "#" + strconv.Itoa(x))

		// set store (TODO: error handle)
		store.Set(uid, cValue.ToJson())

		// close connection
		res.Header().Set("Connection", "close")
		res.Header().Set("Content-Type", "image/png")
		// &http.Cookie{...}.String() will remove the dot prefix from domain
		res.Header().Set("Set-Cookie", "captchaId="+uid+"; Path=/; Domain="+*domain+"; Max-Age="+strconv.Itoa(*expire)+"; HttpOnly")
		res.WriteHeader(200)
		res.Write(png)
	})

	http.HandleFunc("/captcha/verify", func(res http.ResponseWriter, req *http.Request) {
		p_key := req.URL.Query().Get("key")
		p_secret := req.URL.Query().Get("md")
		cid := req.URL.Query().Get("cid")
		cval := req.URL.Query().Get("value")

		res.Header().Set("Connection", "close")

		db := GetDb()
		var key string
		var secret string
		if p_key != "" {
			if err := db.QueryRow("SELECT `key`, `secret` FROM `auth` WHERE `key` = ?", p_key).Scan(&key, &secret); err != nil {
				if err == sql.ErrNoRows {
					res.WriteHeader(400)
					res.Write([]byte("NOT EXISTS SERVICE, "))
					res.Write([]byte(err.Error()))
					return
				} else {
					res.WriteHeader(500)
					res.Write([]byte("UNKNOWN ERROR, "))
					res.Write([]byte(err.Error()))
					return
				}
			}
		}

		// API 권한 인증
		if p_key != key || p_secret != secret {
			res.WriteHeader(400)
			res.Write([]byte("UNKNOWN SERVICE"))
			return
		}

		// Store를 통해 CaptchaID의 구조체 정보 조회
		x, err := store.Get(cid)
		if err != nil {
			//			panic(err)
			// CaptchaID가 존재하지 않는 경우
			res.WriteHeader(400)
			res.Write([]byte(fmt.Sprintf("\"%s\" NOT EXISTS", cid)))
			return
		}

		// Captcha 구조체 Unmarshal
		var cv CacheValue
		if err := json.Unmarshal([]byte(x), &cv); err != nil {
			//panic(err)
			res.WriteHeader(500)
			res.Write([]byte(fmt.Sprint(err)))
			return
		}

		if cv.Result == cval {
			store.Del(cid)
			res.WriteHeader(200)
			res.Write([]byte(x))
		} else {
			res.WriteHeader(400)
			res.Write([]byte("INVALID"))
		}
		return

		// close connection
		res.Header().Set("Connection", "close")
		//res.Header().Set("Content-Type", "image/png")
		// &http.Cookie{...}.String() will remove the dot prefix from domain
		//res.Header().Set("Set-Cookie", "captchaId="+uid+"; Path=/; Domain="+*domain+"; Max-Age="+strconv.Itoa(*expire)+"; HttpOnly")
		res.WriteHeader(200)
		res.Write([]byte(x))
		//res.Write(png)
	})

	fmt.Printf("captcha serve on port: %s \n", *port)
	http.ListenAndServe(":"+*port, nil)
}

// parse args
func parseArgs() {
	flag.Parse()

	assert(*stype == "redis" || *stype == "memcache", "stype must be redis or memcache")
	assert(*saddress != "", "saddress required")

	p, e := strconv.Atoi(*port)
	assert(e == nil && p > 0, "port required and must above 0")

	if *cache > 0 {
		initCache(*cache)
	}

	return
}

func assert(ok bool, msg string) {
	if !ok {
		fmt.Println(msg)
		os.Exit(1)
	}
}

// init cache
func initCache(num int) {
	cachedCaptchaQuene = make([]cachedCaptcha, num)

	for i := 0; i < num; i++ {
		exp, result := captcha.RandomExp()
		img := captcha.NewImage("", exp, 150, 50)

		newCaptcha := cachedCaptcha{
			Result: result,
			Img:    img.EncodedPNG(),
		}

		cachedCaptchaQuene[i] = newCaptcha
	}

	fmt.Println("captcha cache finished")
}

type cachedCaptcha struct {
	Result string
	Img    []byte
}

func getCaptcha() (uid string, result string, png []byte) {
	var num = len(cachedCaptchaQuene)
	num = 0

	uid, _ = captcha.NewUID()

	if num > 0 {
		// from cache
		index := rand.Int() % num

		result = cachedCaptchaQuene[index].Result
		png = cachedCaptchaQuene[index].Img
	} else {
		var exp []byte
		exp, result = captcha.RandomExp()
		img := captcha.NewImage("", exp, 150, 50)
		png = img.EncodedPNG()
	}

	return
}

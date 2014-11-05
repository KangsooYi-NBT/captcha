package main

import "github.com/onebook/go-captcha"
import "math/rand"
import "net/http"
import "strconv"
import "time"
import "fmt"
import "os"

var cachedCaptchaQuene []cachedCaptcha

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	port, stype, saddress, expire, domain := parseArgs(os.Args)

	var store captcha.Store

	if stype == "redis" {
		store = captcha.NewRedisStore(saddress, expire)
	} else {
		store = captcha.NewMemcacheStore(saddress, expire)
	}

	http.HandleFunc("/captcha", func(res http.ResponseWriter, req *http.Request) {
		uid, result, png := getCaptcha()

		// set store
		store.Set(uid, result)

		res.Header().Set("Content-Type", "image/png")
		// &http.Cookie{...}.String() will remove the dot prefix from domain
		res.Header().Set("Set-Cookie", "captchaId="+uid+"; Path=/; Domain="+domain+"; Max-Age="+strconv.Itoa(expire)+"; HttpOnly")
		res.WriteHeader(200)
		res.Write(png)
	})

	fmt.Println("captcha serve on port" + port)
	http.ListenAndServe(port, nil)
}

// parse args
func parseArgs(args []string) (port, stype, saddress string, expire int, domain string) {
	var err error
	var cache int

	for k, v := range args {
		switch v {
		case "--port":
			p, e := strconv.Atoi(args[k+1])
			assert(e == nil, "port must be a number")
			assert(p > 0, "port must above 0")
			port = ":" + args[k+1]
		case "--stype":
			stype = args[k+1]
		case "--saddress":
			saddress = args[k+1]
		case "--expire":
			expire, err = strconv.Atoi(args[k+1])
			assert(err == nil, "expire must be an integer")
		case "--domain":
			domain = args[k+1]
		case "--cache":
			cache, err = strconv.Atoi(args[k+1])
			assert(err == nil, "cache must be an integer")
		}
	}

	assert(stype == "redis" || stype == "memcache", "stype must be redis or memcache")
	assert(saddress != "", "saddress required")
	assert(len(port) > 1, "port required")

	if cache > 0 {
		initCache(cache)
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

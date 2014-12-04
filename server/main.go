package main

import "github.com/onebook/captcha"
import "math/rand"
import "net/http"
import "strconv"
import "flag"
import "time"
import "fmt"
import "os"

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
		uid, result, png := getCaptcha()

		// set store (TODO: error handle)
		store.Set(uid, result)

		// close connection
		res.Header().Set("Connection", "close")
		res.Header().Set("Content-Type", "image/png")
		// &http.Cookie{...}.String() will remove the dot prefix from domain
		res.Header().Set("Set-Cookie", "captchaId="+uid+"; Path=/; Domain="+*domain+"; Max-Age="+strconv.Itoa(*expire)+"; HttpOnly")
		res.WriteHeader(200)
		res.Write(png)
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

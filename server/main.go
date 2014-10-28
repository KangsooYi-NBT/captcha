package main

import "github.com/onebook/go-captcha"
import "net/http"
import "strconv"
import "fmt"
import "os"

func main() {
	stype, saddress, expire, domain := parseArgs(os.Args)

	var store captcha.Store

	if stype == "redis" {
		store = captcha.NewRedisStore(saddress, expire)
	} else {
		store = captcha.NewMemcacheStore(saddress, expire)
	}

	http.HandleFunc("/captcha", func(res http.ResponseWriter, req *http.Request) {
		exp, result := captcha.RandomExp()
		img := captcha.NewImage("", exp, 150, 50)

		uid, _ := captcha.NewUID()

		// set store
		store.Set(uid, result)
		// set cookie
		http.SetCookie(res, &http.Cookie{
			Name:     "captchaId",
			Value:    uid,
			Path:     "/",
			Domain:   domain,
			MaxAge:   expire,
			Secure:   false,
			HttpOnly: true,
		})

		res.WriteHeader(200)
		res.Header().Set("Content-Type", "image/png")
		res.Write(img.EncodedPNG())
	})

	fmt.Println("captcha serve on port: 3000")
	http.ListenAndServe(":3000", nil)
}

func parseArgs(args []string) (stype, saddress string, expire int, domain string) {
	var err error

	for k, v := range args {
		switch v {
		case "--stype":
			stype = args[k+1]
		case "--saddress":
			saddress = args[k+1]
		case "--expire":
			expire, err = strconv.Atoi(args[k+1])
		case "--domain":
			domain = args[k+1]
		}
	}

	if err != nil {
		fmt.Println("expire must be an integer")
		os.Exit(1)
	}

	if stype != "redis" && stype != "memcache" {
		fmt.Println("stype must be redis or memcache")
		os.Exit(1)
	}

	if saddress == "" {
		fmt.Println("saddress required")
		os.Exit(1)
	}

	return
}

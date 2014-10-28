package main

import "../../go-captcha"
import "net/http"
import "fmt"

func main() {
	http.HandleFunc("/captcha", func(res http.ResponseWriter, req *http.Request) {
		exp, result := captcha.RandomExp()
		img := captcha.NewImage("", exp, 150, 50)

		fmt.Println("result:", result)

		res.WriteHeader(200)
		res.Header().Set("Content-Type", "image/png")
		res.Write(img.EncodedPNG())
	})

	http.ListenAndServe(":3000", nil)
}

package main

import "../../go-captcha"
import "os"
import "io"

func main() {
	file, _ := os.Create("test.png")
	defer file.Close()

	var w io.WriterTo
	d := captcha.RandomDigits(6)

	w = captcha.NewImage("", d, 150, 50)

	_, _ = w.WriteTo(file)
}

package main

import "github.com/onebook/go-captcha"
import "fmt"
import "os"
import "io"

func main() {
	file1, _ := os.Create("test1.png")
	file2, _ := os.Create("test2.png")
	defer file1.Close()
	defer file2.Close()

	var w io.WriterTo

	digits := captcha.RandomDigits(6)
	exp, result := captcha.RandomExp()

	w = captcha.NewImage("", digits, 150, 50)
	_, _ = w.WriteTo(file1)
	w = captcha.NewImage("", exp, 150, 50)
	_, _ = w.WriteTo(file2)

	fmt.Printf("result: %s", result)
}

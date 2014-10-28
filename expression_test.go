package captcha

import "testing"

func TestRandomNum(t *testing.T) {
	for i := 0; i < 10; i++ {
		num := randomNum()

		t.Log(num)

		if num < min || num > max {
			t.Fatal(num)
		}
	}
}

func TestRandomExp(t *testing.T) {
	for i := 0; i < 10; i++ {
		exp, result := RandomExp()
		t.Log(exp)
		t.Log(result)
	}
}

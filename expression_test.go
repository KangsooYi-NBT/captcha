package captcha

import "testing"

func TestRandomNum(t *testing.T) {
	for i := 0; i < 10; i++ {
		d := randomNum()

		t.Log(d)

		if d < min || d > max {
			t.Fatal(d)
		}
	}
}

func TestRandomExp(t *testing.T) {
	for i := 0; i < 10; i++ {
		e, r := RandomExp()
		t.Log(e)
		t.Log(r)
	}
}

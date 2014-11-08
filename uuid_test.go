package captcha

import "testing"

func TestNewUID(t *testing.T) {
	uid, e := NewUID()

	if e != nil {
		t.Fatal(e)
	}

	t.Logf("length: %d, uuid: %s", len(uid), uid)

	for i := 0; i < 100000; i++ {
		u, _ := NewUID()

		if u == uid {
			t.Fatal("same uuid")
		}
	}
}

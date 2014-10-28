package captcha

import "testing"

func TestNewUID(t *testing.T) {
	u4, e := NewUID()

	if e != nil {
		t.Fatal(e)
	}

	t.Logf("length: %d, uuid: %s", len(u4), u4)
}

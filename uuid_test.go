package captcha

import "testing"

func TestUuidV4(t *testing.T) {
	u4, e := uuidV4()

	if e != nil {
		t.Fatal(e)
	}

	t.Logf("length: %d, uuid: %s", len(u4), u4)
}

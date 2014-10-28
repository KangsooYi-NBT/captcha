package captcha

import "testing"
import "time"

func TestRedisStore(t *testing.T) {
	s := NewRedisStore("localhost:6379", 1)

	key, value := "10086", "sb SB 移动"

	err := s.Set(key, value)
	assertTestError(err, t)

	result, err := s.Get(key)
	assertTestError(err, t)

	if result != value {
		t.Fatalf("result: %s not correct", result)
	}

	time.Sleep(time.Second * 1)

	noResult, err := s.Get(key)
	assertTestError(err, t)

	if noResult != "" {
		t.Fatal("noResult: %s", noResult)
	}
}

func TestMemcacheStore(t *testing.T) {
	s := NewMemcacheStore("localhost:11211", 1)

	key, value := "10086", "sb SB 移动"

	err := s.Set(key, value)
	assertTestError(err, t)

	result, err := s.Get(key)
	assertTestError(err, t)

	if result != value {
		t.Fatalf("result: %s not correct", result)
	}

	time.Sleep(time.Second * 1)

	noResult, err := s.Get(key)

	if err.Error() != "memcache: cache miss" {
		t.Fatal("cache not miss")
	}

	if noResult != "" {
		t.Fatal("noResult: %s", noResult)
	}
}

func assertTestError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

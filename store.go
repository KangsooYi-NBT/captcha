package captcha

import "github.com/bradfitz/gomemcache/memcache"
import "github.com/garyburd/redigo/redis"
import "time"

type Store interface {
	Set(id string, result string) error
	Get(id string) (string, error)
}

// redis store
type redisStore struct {
	Address string
	Expire  int
	Pool    *redis.Pool
}

func NewRedisStore(address string, expire int) Store {
	s := new(redisStore)
	s.Address = address
	s.Expire = expire

	pool := &redis.Pool{
		MaxIdle:     20,
		MaxActive:   100,
		IdleTimeout: 10 * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", s.Address)

			if err != nil {
				panic(err)
			}

			return
		},
	}

	s.Pool = pool

	return s
}

func (s *redisStore) Set(id, result string) error {
	_, err := s.Pool.Get().Do("SETEX", id, s.Expire, result)
	return err
}

func (s *redisStore) Get(id string) (result string, err error) {
	data, err := s.Pool.Get().Do("GET", id)
	if data == nil {
		result = ""
	} else {
		result = string(data.([]byte)[:])
	}

	return
}

// memcache store
type memcacheStore struct {
	Address string
	Expire  int32
	Client  *memcache.Client
}

func NewMemcacheStore(address string, expire int32) Store {
	s := new(memcacheStore)
	s.Address = address
	s.Expire = expire

	client := memcache.New(address)
	s.Client = client

	return s
}

func (s *memcacheStore) Set(id, result string) error {
	return s.Client.Set(&memcache.Item{
		Key:        id,
		Value:      []byte(result),
		Expiration: s.Expire,
	})
}

func (s *memcacheStore) Get(id string) (result string, err error) {
	item, err := s.Client.Get(id)

	if err != nil {
		return "", err
	}

	result = string(item.Value)
	return
}

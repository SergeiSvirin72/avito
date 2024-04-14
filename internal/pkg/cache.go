package pkg

import "github.com/redis/go-redis/v9"

func NewCache(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

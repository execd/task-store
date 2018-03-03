package store

import "github.com/go-redis/redis"

func NewClient(address string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

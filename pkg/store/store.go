package store

import "github.com/go-redis/redis"

// NewClient creates a redis client
func NewClient(address string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

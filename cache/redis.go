package cache

import (
	"github.com/go-redis/redis/v8"
)

// Redis cache object shared across the application
var Redis *redis.Client

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

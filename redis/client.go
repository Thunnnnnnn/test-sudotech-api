package redisclient

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	RedisEndpoint := os.Getenv("REDIS_ENDPOINT")

	if RedisEndpoint == "" {
		panic("REDIS_ENDPOINT is empty")
	}

	return redis.NewClient(&redis.Options{
		Addr:     RedisEndpoint,
		Password: "",
		DB:       0,
		// TLSConfig: &tls.Config{},
	})
}

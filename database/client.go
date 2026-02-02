package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx context.Context
	RDB *redis.Client
)

func InitRedis() {
	Ctx = context.Background()

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
		// TLSConfig: &tls.Config{},
	})

	// optional: test connection
	if err := RDB.Ping(Ctx).Err(); err != nil {
		panic(err)
	}
}

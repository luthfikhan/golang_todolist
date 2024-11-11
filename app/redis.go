package app

import (
	"os"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

func InitiateRedis() *redis.Client {
	if os.Getenv("GO_ENV") == "test" {
		rdb, _ := redismock.NewClientMock()

		return rdb
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return rdb
}

package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *Configuration) *redis.Client {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Url,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis/Valkey:", err)
	}
	return rdb
}

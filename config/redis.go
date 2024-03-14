package config

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *RedisOption) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Host,
		Password:     cfg.Password,
		DB:           0, // use default DB
		ReadTimeout:  time.Duration(cfg.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeoutMs) * time.Millisecond,
	})

	return rdb
}

package config

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	apperr "github.com/rahmatrdn/go-skeleton/error"
)

func NewRedis(config *Config) *redis.Client {
	redisStore := redis.NewClient(&redis.Options{
		Addr:         config.RedisOption.Host,
		ReadTimeout:  time.Duration(config.RedisOption.ReadTimeoutMs),
		WriteTimeout: time.Duration(config.RedisOption.ReadTimeoutMs),
	})
	return redisStore
}

func SetCache[T any](cacheManager *redis.Client, ctx context.Context, prefix string, key string, executeData func(context.Context, string) (T, error)) *T {
	var data []byte
	var object T
	if err := cacheManager.Get(ctx, prefix+"_"+key).Scan(&data); err == nil {
		err := json.Unmarshal(data, &object)
		apperr.PanicLogging(err)

		return &object
	}
	value, err := executeData(ctx, key)
	if err != nil {
		panic(apperr.NotFoundError{
			Message: err.Error(),
		})
	}
	cacheValue, err := json.Marshal(value)
	apperr.PanicLogging(err)

	if err := cacheManager.Set(ctx, prefix+"_"+key, cacheValue, -1).Err(); err != nil {
		apperr.PanicLogging(err)
	}
	return &value
}

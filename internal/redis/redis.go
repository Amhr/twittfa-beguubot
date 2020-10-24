package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedis(v *viper.Viper) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	return &RedisCache{rdb: rdb}
}

func (r *RedisCache) Get(key string, def string, ctx context.Context) string {
	item := r.rdb.Get(ctx, key)
	if item.Err() != nil {
		return def
	}
	return item.Val()
}

func (r *RedisCache) Set(key string, val string, ttl time.Duration, ctx context.Context) *redis.StatusCmd {
	return r.rdb.Set(ctx, key, val, ttl)
}

func (r *RedisCache) Key(items ...string) string {
	return strings.Join(items, ":")
}

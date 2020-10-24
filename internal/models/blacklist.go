package models

import (
	"context"
	"github.com/amhr/begubot/internal/redis"
	"strconv"
	"time"
)

func Block(f int, t int, c *redis.RedisCache) {
	c.Set(c.Key("blacklist", strconv.Itoa(f), strconv.Itoa(t)), "true", time.Duration(30*24)*time.Hour, context.Background())
}

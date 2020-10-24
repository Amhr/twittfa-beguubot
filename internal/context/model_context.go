package context

import (
	"github.com/amhr/begubot/internal/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

type ModelContext struct {
	DB    *gorm.DB
	Redis *redis.RedisCache
	Bot   *tgbotapi.BotAPI
}

func NewContextModel(db *gorm.DB, redis *redis.RedisCache, b *tgbotapi.BotAPI) *ModelContext {
	return &ModelContext{
		DB:    db,
		Redis: redis,
		Bot:   b,
	}
}

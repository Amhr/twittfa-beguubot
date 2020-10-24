package internal

import (
	epimetheus2 "github.com/amhr/begubot/internal/epimetheus"
	"github.com/amhr/begubot/internal/location"
	"github.com/amhr/begubot/internal/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

/**
Beguu bot is the main type of the application.
this holds other dependencies that application needs
*/
type BeguuBot struct {
	Config    *viper.Viper
	Metrics   *epimetheus2.MetricsManager
	Dev       bool
	Bot       *tgbotapi.BotAPI
	Exit      chan bool
	Updates   chan *tgbotapi.Update
	Cache     *redis.RedisCache
	Locations []location.Locationer
	DB        *gorm.DB
}

func NewBeguu(
	c *viper.Viper,
	ep *epimetheus2.MetricsManager,
	b *tgbotapi.BotAPI,
	rc *redis.RedisCache,
	db *gorm.DB,
) *BeguuBot {
	dev := c.Get("dev") == true

	if dev {
		b.Debug = false
	}

	locations := make([]location.Locationer, 0)
	//static important
	locations = append(locations, location.NewStartLocation(ep, b))
	locations = append(locations, location.NewCancelLocation(ep, b))

	// other routes
	locations = append(locations, location.NewHomeLocation(ep, b))
	locations = append(locations, location.NewSendAnnmsgLocation(ep, b))
	locations = append(locations, location.NewMyLinkLocation(ep, b))

	// and 404
	locations = append(locations, location.New404Location(ep, b))

	return &BeguuBot{
		Config:    c,
		Metrics:   ep,
		Dev:       dev,
		Bot:       b,
		Exit:      make(chan bool, 1),
		Updates:   make(chan *tgbotapi.Update, 10),
		Cache:     rc,
		Locations: locations,
		DB:        db,
	}
}

func (b *BeguuBot) Run() {

	// running workers
	totalWorkers := b.Config.GetInt("WORKERS")
	for i := 0; i < totalWorkers; i++ {
		go b.Worker(int32(i))
	}

	// fetching updates ..
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.Bot.GetUpdatesChan(u)

	if err != nil {
		logrus.Errorf("error while trying to start get update %w", err)
	}

	for update := range updates {
		if update.Message != nil || update.CallbackQuery != nil { // ignore any non-Message Updates
			b.Updates <- &update
		}

	}

	<-b.Exit

}

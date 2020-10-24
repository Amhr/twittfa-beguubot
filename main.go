package main

import (
	"fmt"
	"github.com/amhr/begubot/internal"
	epimetheus2 "github.com/amhr/begubot/internal/epimetheus"
	"github.com/amhr/begubot/internal/logrus"
	"github.com/amhr/begubot/internal/models"
	"github.com/amhr/begubot/internal/redis"
	viper2 "github.com/amhr/begubot/internal/viper"
	"github.com/cafebazaar/epimetheus"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
)

func main() {

	// Setup Config
	v, err := viper2.NewViper()
	if err != nil {
		fmt.Errorf("error %w", err)
	}

	// Setup Logrus
	logrus.MakeLogrus(v)

	//setup Epimetheus Server
	epserv := epimetheus.NewEpimetheusServer(v)
	go epserv.Listen()
	ep := epimetheus2.NewMetricsManager(epserv)

	// telegram bot

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	// Redis Cache
	rc := redis.NewRedis(v)

	// maria db
	db, err := models.NewGorm()
	if err != nil {
		panic(err)
	}

	//make beguu
	b := internal.NewBeguu(v, ep, bot, rc, db)
	b.Run()

}

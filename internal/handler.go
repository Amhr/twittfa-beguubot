package internal

import (
	"github.com/amhr/begubot/internal/callbacks"
	"github.com/amhr/begubot/internal/context"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func (b *BeguuBot) HandleMessage(u *tgbotapi.Update) {
	t := b.Metrics.Requests.Start()
	defer b.Metrics.Requests.Done(t, "Message", "main", "ok")

	currentUser := models.NewUser(u.Message.From, b.Cache, b.DB, b.Metrics)
	currentUser.Load()

	for _, loc := range b.Locations {
		loc.ForceLocation(currentUser, u)
	}
	for _, loc := range b.Locations {
		if loc.IsValid(currentUser, u) {
			loc.Run(currentUser, u)
			b.Metrics.Requests.Done(t, "Message", loc.GetName(), "ok")
			break
		}
	}

}

func (b *BeguuBot) HandleCallback(u *tgbotapi.Update) {
	t := b.Metrics.Requests.Start()

	currentUser := models.NewUser(u.CallbackQuery.From, b.Cache, b.DB, b.Metrics)
	currentUser.Load()

	msg := u.CallbackQuery.Data
	exp := strings.Split(msg, "-")
	defer b.Metrics.Requests.Done(t, "Callback", exp[0], "ok")
	switch exp[0] {
	case "open":
		callbacks.OpenCallback(currentUser, u, context.NewContextModel(b.DB, b.Cache, b.Bot), exp)
	case "reply":
		callbacks.ReplyCallback(currentUser, u, context.NewContextModel(b.DB, b.Cache, b.Bot), exp)
	case "block":
		callbacks.BlockCallback(currentUser, u, context.NewContextModel(b.DB, b.Cache, b.Bot), exp)

	}

}

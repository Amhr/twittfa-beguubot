package location

import (
	"github.com/amhr/begubot/internal/epimetheus"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Location404 struct {
	Metrics *epimetheus.MetricsManager
	bot     *tgbotapi.BotAPI
}

func New404Location(m *epimetheus.MetricsManager, b *tgbotapi.BotAPI) *Location404 {
	return &Location404{Metrics: m, bot: b}
}

func (l Location404) IsValid(u *models.UserManager, up *tgbotapi.Update) bool {
	return true
}

func (l Location404) Run(u *models.UserManager, up *tgbotapi.Update) {
	c := tgbotapi.NewMessage(u.ID64(), "اشتباه است")
	c.ReplyMarkup = keyboards.HomeKeyboard()
	l.bot.Send(c)
}

func (l Location404) GetName() string {
	return "404"
}

func (l Location404) ForceLocation(u *models.UserManager, up *tgbotapi.Update) {
	//
}

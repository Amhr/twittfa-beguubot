package location

import (
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Locationer interface {
	IsValid(u *models.UserManager, up *tgbotapi.Update) bool
	Run(u *models.UserManager, up *tgbotapi.Update)
	GetName() string
	ForceLocation(u *models.UserManager, up *tgbotapi.Update)
}

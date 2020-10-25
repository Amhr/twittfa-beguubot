package location

import (
	"github.com/amhr/begubot/internal/epimetheus"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type LocationCancel struct {
	Metrics *epimetheus.MetricsManager
	bot     *tgbotapi.BotAPI
}

func NewCancelLocation(m *epimetheus.MetricsManager, b *tgbotapi.BotAPI) *LocationCancel {
	return &LocationCancel{Metrics: m, bot: b}
}

func (l LocationCancel) IsValid(u *models.UserManager, up *tgbotapi.Update) bool {
	return u.GetLocation() == "cancel"
}

func (l LocationCancel) Run(u *models.UserManager, up *tgbotapi.Update) {
}

func (l LocationCancel) GetName() string {
	return "cancel"
}

func (l LocationCancel) ForceLocation(u *models.UserManager, up *tgbotapi.Update) {
	s := strings.Split(up.Message.Text, " ")
	if s[0] == "/cancel" || up.Message.Text == keyboards.TXT_CANCEL {
		msgs := u.GetWaitingMsgs()
		for _, msgId := range msgs {
			msg := models.GetMessage(msgId, u.ContextModel)
			if msg != nil {
				msg.Cancel(u)
			}
		}
		u.SetWaitingMsgs([]int{})
		u.ClearCache()
	}
}

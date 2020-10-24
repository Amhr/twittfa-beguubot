package location

import (
	"github.com/amhr/begubot/internal/epimetheus"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type LocationStart struct {
	Metrics *epimetheus.MetricsManager
	bot     *tgbotapi.BotAPI
}

func NewStartLocation(m *epimetheus.MetricsManager, b *tgbotapi.BotAPI) *LocationStart {
	return &LocationStart{Metrics: m, bot: b}
}

func (l LocationStart) IsValid(u *models.UserManager, up *tgbotapi.Update) bool {
	return u.GetLocation() == "start"
}

func (l LocationStart) Run(u *models.UserManager, up *tgbotapi.Update) {
	cmd := strings.Split(up.Message.Text, " ")
	// check if its send message
	if len(cmd) == 2 {

		identifier := cmd[1]
		expIdent := strings.Split(identifier, "_")
		if len(expIdent) == 2 {
			u.ClearCache()
			u.SetCache("annmsg_id", expIdent[0])
			u.SetLocation("annmsg")
			u.SetStep("1")
			r := NewSendAnnmsgLocation(l.Metrics, l.bot)
			r.Run(u, up)
			return
		}

		c := tgbotapi.NewMessage(u.ID64(), `⚠️ لینک وارد شده اشتباه است.
لینکی که از طریقش وارد ربات شدید اشتباهه. `)
		l.bot.Send(c)
		return
	} else {
		c := tgbotapi.NewMessage(u.ID64(), `📨 به ربات پیام ناشناس خوش اومدید

اگه میخوای بقیه بهت پیام بدن روی لینک ناشناس من کلیک کن ، لینک ناشناست رو توی توئیتر یا اینستاگرام بذار تا بقیه بهت پیام بدن.`)
		c.ReplyMarkup = keyboards.HomeKeyboard()
		l.bot.Send(c)
		return
	}

}

func (l LocationStart) GetName() string {
	return "Start"
}

func (l LocationStart) ForceLocation(u *models.UserManager, up *tgbotapi.Update) {
	s := strings.Split(up.Message.Text, " ")
	if s[0] == "/start" {
		u.SetLocation("start")
	}
}

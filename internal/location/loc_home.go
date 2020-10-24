package location

import (
	"fmt"
	"github.com/amhr/begubot/internal/epimetheus"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type LocationHome struct {
	Metrics *epimetheus.MetricsManager
	bot     *tgbotapi.BotAPI
}

func NewHomeLocation(m *epimetheus.MetricsManager, b *tgbotapi.BotAPI) *LocationHome {
	return &LocationHome{Metrics: m, bot: b}
}

func (l LocationHome) IsValid(u *models.UserManager, up *tgbotapi.Update) bool {
	return u.GetLocation() == "home"
}

func (l LocationHome) Run(u *models.UserManager, up *tgbotapi.Update) {
	c := tgbotapi.NewMessage(u.ID64(), fmt.Sprintf(`📨 به ربات پیام ناشناس خوش اومدید

اگه میخوای بقیه بهت پیام بدن روی %s کلیک کن ، لینک ناشناست رو توی توئیتر یا اینستاگرام بذار تا بقیه بهت پیام بدن.`, keyboards.TXT_MY_LINK))
	c.ReplyMarkup = keyboards.HomeKeyboard()
	l.bot.Send(c)
}

func (l LocationHome) GetName() string {
	return "Home"
}

func (l LocationHome) ForceLocation(u *models.UserManager, up *tgbotapi.Update) {
	if up.Message.Text == keyboards.TXT_HOME {
		u.SetLocation("home")
	}
}

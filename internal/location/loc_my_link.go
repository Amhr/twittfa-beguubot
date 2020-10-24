package location

import (
	"fmt"
	"github.com/amhr/begubot/internal/epimetheus"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type LocationMyLink struct {
	Metrics *epimetheus.MetricsManager
	bot     *tgbotapi.BotAPI
}

func NewMyLinkLocation(m *epimetheus.MetricsManager, b *tgbotapi.BotAPI) *LocationMyLink {
	return &LocationMyLink{Metrics: m, bot: b}
}

func (l LocationMyLink) IsValid(u *models.UserManager, up *tgbotapi.Update) bool {
	return u.GetLocation() == "my_link"
}

func (l LocationMyLink) Run(u *models.UserManager, up *tgbotapi.Update) {
	link := fmt.Sprintf("https://t.me/%s?start=%s", "BeguuBot", u.MyLinkIdentifier())
	c := tgbotapi.NewMessage(u.ID64(), fmt.Sprintf(`📨 لینک صندوق پیام شما

این لینک رو هر جا میخوای بذار ، از طریق این اینک هر پیامی دریافت کنید بهت اطلاع میدم.

%s`, link))
	c.ReplyMarkup = keyboards.HomeKeyboard()
	l.bot.Send(c)
	u.SetLocation("home")
}

func (l LocationMyLink) GetName() string {
	return "Mylink"
}

func (l LocationMyLink) ForceLocation(u *models.UserManager, up *tgbotapi.Update) {
	if up.Message.Text == keyboards.TXT_MY_LINK {
		u.SetLocation("my_link")
	}
}

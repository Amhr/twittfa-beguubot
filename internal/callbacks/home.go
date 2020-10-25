package callbacks

import (
	"github.com/amhr/begubot/internal/context"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HomeCallback(u *models.UserManager, update *tgbotapi.Update, c *context.ModelContext, data []string) {
	u.ClearCache()
	sendAble := tgbotapi.NewMessage(u.ID64(), `حله.

چه کاری برات انحام بدم؟`)
	sendAble.ReplyMarkup = keyboards.HomeKeyboard()
	c.Bot.Send(sendAble)
	c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "درحال انجام ..."))
}

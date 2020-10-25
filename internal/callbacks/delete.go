package callbacks

import (
	"github.com/amhr/begubot/internal/context"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func DeleteCallback(u *models.UserManager, update *tgbotapi.Update, c *context.ModelContext, data []string) {
	nf := func() {
		c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "یافت نشد"))
	}
	if len(data) != 2 {
		nf()
		return
	}
	msgId := data[1]
	msgIdInt, _ := strconv.Atoi(msgId)
	msg := models.GetMessage(msgIdInt, c)
	if msg.FromId != u.UserMessage.DatabaseID {
		nf()
		return
	}

	if msg.Status != 2 {
		c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "این پیام قبلا ارسال شده است"))
		return
	}

	otherUser := u.GetUserBy("id", msg.ToId)
	if otherUser == nil {
		c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "⚠️ مشکلی پیش آمد! لینک ناشناس اشتباه میباشد"))
		return
	}
	msg.Cancel(u)
	c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "انجام شد"))
}

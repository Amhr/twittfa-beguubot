package callbacks

import (
	"github.com/amhr/begubot/internal/context"
	"github.com/amhr/begubot/internal/location"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func ReplyCallback(u *models.UserManager, update *tgbotapi.Update, c *context.ModelContext, data []string) {
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
	if msg.ToId != u.UserMessage.DatabaseID {
		nf()
		return
	}
	go c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "درحال ارسال ..."))
	sendable := tgbotapi.NewMessage(u.ID64(), "در حال پاسخ به این پیام هستید")
	sendable.ReplyToMessageID = update.CallbackQuery.Message.MessageID
	c.Bot.Send(sendable)
	u.SetCache("annmsg_id", strconv.Itoa(msg.FromId))
	u.SetCache("annmsg_reply", strconv.Itoa(msg.ID))
	u.SetStep("1")
	u.SetLocation("annmsg")
	e := location.NewSendAnnmsgLocation(u.Metrics, c.Bot)
	e.Run(u, update)

}

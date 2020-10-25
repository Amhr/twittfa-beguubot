package callbacks

import (
	"github.com/amhr/begubot/internal/context"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func OpenCallback(u *models.UserManager, update *tgbotapi.Update, c *context.ModelContext, data []string) {
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
	// check if msg has been already opened
	if msg.Status != 0 {
		c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "این پیام رو قبلا باز کرده بودید!"))
		go c.Bot.Send(tgbotapi.NewDeleteMessage(u.ID64(), update.CallbackQuery.Message.MessageID))
		return
	}
	// update message seen
	msg.Status = 1
	msg.SaveCache(c.Redis)

	// send message
	replyTo := models.GetMessage(msg.ReplyTo, c)
	d := keyboards.MessageDetailKeyboard(msgIdInt, u.IsBlocked(msg.FromId))
	c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "درحال ارسال ..."))
	go c.Bot.Send(tgbotapi.NewDeleteMessage(u.ID64(), update.CallbackQuery.Message.MessageID))
	replyMessageId := 0
	if replyTo.FromId == u.UserMessage.DatabaseID {
		r, err := c.Bot.Send(models.SendMessage(replyTo, u.ID64(), nil, 0))
		if err == nil {
			replyMessageId = r.MessageID
		}
	}
	sendableMsg := models.SendMessage(msg, u.ID64(), &d, replyMessageId)
	c.Bot.Send(sendableMsg)

	// send message seen feedback
	otherUser := u.GetUserBy("db", msg.FromId)
	feedbackSendable := tgbotapi.NewMessage(int64(otherUser.TelegramID), "👀 این پیام ات رو دید.")
	feedbackSendable.ReplyToMessageID = msg.SenderMessageID
	c.Bot.Send(feedbackSendable)

}

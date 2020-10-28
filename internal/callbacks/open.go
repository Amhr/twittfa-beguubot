package callbacks

import (
	"fmt"
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
	msgHolderIdInt, _ := strconv.Atoi(msgId)
	msgHolder := models.GetMessage(msgHolderIdInt, c)
	if msgHolder.ToId != u.UserMessage.DatabaseID {
		nf()
		return
	}
	// check if msg has been already opened
	if msgHolder.Status != 0 {
		c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "این پیام رو قبلا باز کرده بودید!"))
		go c.Bot.Send(tgbotapi.NewDeleteMessage(u.ID64(), update.CallbackQuery.Message.MessageID))
		return
	}
	// update message seen
	msgHolder.Status = 1
	msgHolder.SaveCache(c.Redis)

	var proccableMsgs []int

	if msgHolder.Type == "GROUP" {
		proccableMsgs = msgHolder.Msgs()
	} else {
		proccableMsgs = []int{msgHolderIdInt}
	}

	c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "درحال ارسال ..."))

	// send messages
	var lastMessageId int

	for _, msgIdInt := range proccableMsgs {
		msg := models.GetMessage(msgIdInt, u.ContextModel)
		d := keyboards.MessageDetailKeyboard(msgIdInt, u.IsBlocked(msg.FromId))

		replyMessageId := 0
		replyTo := models.GetMessage(msg.ReplyTo, c)
		if replyTo.FromId == u.UserMessage.DatabaseID {
			if replyTo.Type == "Text" {
				replyTo.Data = fmt.Sprintf(`پیامی که شما ارسال کردید:

%s`, replyTo.Data)
			} else {
				replyTo.Caption = fmt.Sprintf(`پیامی که شما ارسال کردید

%s`, replyTo.Caption)
			}
			r, err := c.Bot.Send(models.SendMessage(replyTo, u.ID64(), nil, 0))
			if err == nil {
				replyMessageId = r.MessageID
			}
		}

		sendableMsg := models.SendMessage(msg, u.ID64(), &d, replyMessageId)
		c.Bot.Send(sendableMsg)
		msg.Status = 1
		msg.SaveCache(c.Redis)
		// send message seen feedback
		if msgHolder.Type != "GROUP" {
			otherUser := u.GetUserBy("db", msg.FromId)
			feedbackSendable := tgbotapi.NewMessage(int64(otherUser.TelegramID), "👀 این پیام ات رو دید.")
			feedbackSendable.ReplyToMessageID = msg.SenderMessageID
			c.Bot.Send(feedbackSendable)
		} else {
			lastMessageId = msg.SenderMessageID
		}
	}

	go c.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(u.ID64(), update.CallbackQuery.Message.MessageID, tgbotapi.NewInlineKeyboardMarkup()))
	go c.Bot.Send(tgbotapi.NewEditMessageText(u.ID64(), update.CallbackQuery.Message.MessageID, `📩 نمایش پیام های جدید :`))

	if msgHolder.Type == "GROUP" {
		fmt.Println(lastMessageId)
		otherUser := u.GetUserBy("db", msgHolder.FromId)
		fmt.Println(otherUser)
		feedbackSendable := tgbotapi.NewMessage(int64(otherUser.TelegramID), "👀 این چند تا پیامی که فرستاده بودی رو دید")
		feedbackSendable.ReplyToMessageID = lastMessageId
		c.Bot.Send(feedbackSendable)
	}

}

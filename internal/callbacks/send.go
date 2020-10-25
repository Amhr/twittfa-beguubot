package callbacks

import (
	"github.com/amhr/begubot/internal/context"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func SendCallback(u *models.UserManager, update *tgbotapi.Update, c *context.ModelContext, data []string) {
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

	msg.Status = 0
	msg.SaveCache(c.Redis)

	otherUser := u.GetUserBy("id", msg.ToId)
	if otherUser == nil {
		c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "⚠️ مشکلی پیش آمد! لینک ناشناس اشتباه میباشد"))
		return
	}
	otherUserSend := tgbotapi.NewMessage(int64(otherUser.TelegramID), `💌 یک پیام جدید دریافت کردید!
برای نمایش پیام روی دکمه نمایش کلیک کنید.`)
	otherUserSend.ReplyMarkup = keyboards.ShowMessageKeyboard(msg.ID)
	go c.Bot.Send(otherUserSend)
	done := tgbotapi.NewMessage(u.ID64(), `👆 پیام بالا برای مخاطبتون ارسال شد.
هر وقت پیام رو ببینه بهت اطلاع میدم.`)
	done.ReplyMarkup = keyboards.CancelKeyboard()
	done.ReplyToMessageID = msg.SenderMessageID
	go c.Bot.Send(done)
	go c.Bot.Send(tgbotapi.NewDeleteMessage(u.ID64(), msg.BotPreviewMessageID))
	u.UnsetFromWaitingMsgs(msg.ID)

}

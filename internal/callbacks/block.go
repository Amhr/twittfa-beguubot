package callbacks

import (
	"github.com/amhr/begubot/internal/context"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func BlockCallback(u *models.UserManager, update *tgbotapi.Update, c *context.ModelContext, data []string) {
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

	if u.IsBlocked(msg.FromId) {
		u.Unblock(msg.FromId)
	} else {
		u.Block(msg.FromId)
	}

	go c.Bot.Send(tgbotapi.NewEditMessageReplyMarkup(u.ID64(), update.CallbackQuery.Message.MessageID, keyboards.MessageDetailKeyboard(msgIdInt, u.IsBlocked(msg.FromId))))
	go c.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "انجام شد"))
}

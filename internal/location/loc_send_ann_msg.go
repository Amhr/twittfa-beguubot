package location

import (
	"fmt"
	"github.com/amhr/begubot/internal/epimetheus"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"strconv"
)

type LocationSendAnnmsg struct {
	Metrics *epimetheus.MetricsManager
	bot     *tgbotapi.BotAPI
}

func NewSendAnnmsgLocation(m *epimetheus.MetricsManager, b *tgbotapi.BotAPI) *LocationSendAnnmsg {
	return &LocationSendAnnmsg{Metrics: m, bot: b}
}

func (l LocationSendAnnmsg) IsValid(u *models.UserManager, up *tgbotapi.Update) bool {
	return u.GetLocation() == "annmsg"
}

func (l LocationSendAnnmsg) Run(u *models.UserManager, up *tgbotapi.Update) {

	step := u.GetStep()

	c := u.GetCache("annmsg_id")
	reply_to := u.GetCache("annmsg_reply")
	if c == "" {
		u.Error("⚠️ مشکلی پیش آمد! لینک ناشناس اشتباه میباشد", l.bot)
		return
	}
	id, _ := strconv.Atoi(c)
	usr := u.GetUserBy("db", id)
	if usr == nil {
		u.Error("⚠️ مشکلی پیش آمد! لینک ناشناس اشتباه میباشد", l.bot)
		return
	}
	if u.ImBlocked(usr.DatabaseID) {
		u.Error("🔒 متاسفانه شما توسط این کاربر بلاک شده اید", l.bot)
		u.ClearCache()
		return
	}

	switch step {
	case "1":

		msg := tgbotapi.NewMessage(u.ID64(), `📩 درحال ارسال پیام ناشناس میباشد:
هر حرفی تو دلت هست بگو. هیچ اسمی از تو ذخیره نمیشه.`)
		if reply_to != "" {
			msg = tgbotapi.NewMessage(u.ID64(), `📩 هر جوابی که میخوای میتونی ارسال کنی:`)
		}
		msg.ReplyMarkup = keyboards.CancelKeyboard()
		go l.bot.Send(msg)
		u.SetStep("2")
	case "2":
		msg := models.ConvertUpdateToAnnmsg(up)
		if msg == nil {
			c := tgbotapi.NewMessage(u.ID64(), `🚫 پیامی که ارسال کردید توسط ربات پشتیبانی نمیشود!

لطفا یک پیام جدید ارسال کنید`)
			c.ReplyMarkup = keyboards.CancelKeyboard()
			go l.bot.Send(c)
			return
		}
		msg.FromId = u.UserMessage.DatabaseID
		msg.ToId = usr.DatabaseID
		if reply_to != "" {
			msg.ReplyTo, _ = strconv.Atoi(reply_to)
		}
		_, err := msg.Save(u.DB, u.Cache)
		if err != nil {
			u.Error(`مشکلی پیش آمد.
لطفا چند دقیقه دیگه دوباره تلاش کنید`, l.bot)
			logrus.WithField("action", "SendAnnMsgSave").Error(err)
		}

		otherUser := u.GetUserBy("id", msg.ToId)
		if otherUser == nil {
			u.Error("⚠️ مشکلی پیش آمد! لینک ناشناس اشتباه میباشد", l.bot)
			return
		}
		send := models.SendMessage(msg, u.ID64(), nil, 0)
		if send != nil {
			otherUserSend := tgbotapi.NewMessage(int64(otherUser.TelegramID), `💌 یک پیام جدید دریافت کردید!
برای نمایش پیام روی دکمه نمایش کلیک کنید.`)
			otherUserSend.ReplyMarkup = keyboards.ShowMessageKeyboard(msg.ID)
			go l.bot.Send(otherUserSend)
			done := tgbotapi.NewMessage(u.ID64(), `👆 پیام بالا برای مخاطبتون ارسال شد.
هر وقت پیام رو ببینه بهت اطلاع میدم.`)
			done.ReplyMarkup = keyboards.HomeKeyboard()
			l.bot.Send(done)

			// finish message sending ids
			msg.SenderMessageID = up.Message.MessageID
			msg.SaveCache(u.Cache)
			u.ClearCache()

		} else {
			c := tgbotapi.NewMessage(u.ID64(), `🚫 پیامی که ارسال کردید توسط ربات پشتیبانی نمیشود!

لطفا یک پیام جدید ارسال کنید`)
			c.ReplyMarkup = keyboards.CancelKeyboard()
			go l.bot.Send(c)
			return
		}

		fmt.Println(msg)

	}

}

func (l LocationSendAnnmsg) GetName() string {
	return "annmsg"
}

func (l LocationSendAnnmsg) ForceLocation(u *models.UserManager, up *tgbotapi.Update) {

}

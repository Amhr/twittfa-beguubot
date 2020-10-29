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
هر حرفی تو دلت هست بگو. هیچ اسمی از تو ذخیره نمیشه.

میتونی پیام هات رو با هم فورووارد کنی یا یکی یکی ارسال کنی. قبل از ارسالشون بهت نشون میدیم پیامی که میخوای بفرستی رو و اگه تایید کردی ارسال میکنیم`)
		if reply_to != "" {
			msg = tgbotapi.NewMessage(u.ID64(), `📩 هر جوابی که میخوای میتونی ارسال کنی:`)
		}
		msg.ReplyMarkup = keyboards.SendAnnmsgKeyboard()
		go l.bot.Send(msg)
		u.SetStep("2")
	case "2":

		if up.Message.Text == keyboards.TXT_SEND {
			if len(u.GetWaitingMsgs()) > 0 {
				l.FinishSendMessage(u, up)
				return
			} else {
				c := tgbotapi.NewMessage(u.ID64(), `🚫 هنوز پیامی ارسال نکردی!`)
				c.ReplyMarkup = keyboards.SendAnnmsgKeyboard()
				go l.bot.Send(c)
				return
			}
		}

		msg := models.ConvertUpdateToAnnmsg(up)
		if msg == nil {
			c := tgbotapi.NewMessage(u.ID64(), `🚫 پیامی که ارسال کردید توسط ربات پشتیبانی نمیشود!

لطفا یک پیام جدید ارسال کنید`)
			c.ReplyMarkup = keyboards.SendAnnmsgKeyboard()
			go l.bot.Send(c)
			return
		}
		if len(u.GetWaitingMsgs()) > 4 {
			c := tgbotapi.NewMessage(u.ID64(), `🚫 بیشتر از ۵ تا پیام همزمان نمیتونی ارسال کنی!
یا روی ارسال کلیک کن یا یکی از پیام های قبلی رو حذف کن`)
			c.ReplyMarkup = keyboards.SendAnnmsgKeyboard()
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
		// send multiple message
		msg.Status = 2
		msg.SenderMessageID = up.Message.MessageID
		k := keyboards.FinishSendMessageKeyboard(msg.ID)
		send := tgbotapi.NewMessage(u.ID64(), `⚙️ مدیریت پیام`)
		send.ReplyMarkup = k
		send.ReplyToMessageID = up.Message.MessageID
		r, e := l.bot.Send(send)
		if e == nil {
			msg.BotPreviewMessageID = r.MessageID
			msg.SaveCache(u.Cache)
		}
		u.AddWaitingMsg(msg.ID)
		manageMsg := tgbotapi.NewMessage(u.ID64(), fmt.Sprintf(`👍خب %d پیام آماده ارساله.
ادامه پیام ها رو ارسال کن یا اگر خواستی قبلی ها رو حذف کن.

در انتها فراموش نکن روی %s کلیک کنی.`, len(u.GetWaitingMsgs()), keyboards.TXT_SEND))
		manageMsg.ReplyMarkup = keyboards.SendAnnmsgKeyboard()
		r, err = l.bot.Send(manageMsg)
		if err == nil {
			u.DelDeletableMsgs(l.bot)
			u.AddDeletableMsg(r.MessageID)
		}
		return

		// old school send message
		if true {
			otherUserSend := tgbotapi.NewMessage(int64(otherUser.TelegramID), `💌 یک پیام جدید دریافت کردید!
برای نمایش پیام روی دکمه نمایش کلیک کنید.`)
			otherUserSend.ReplyMarkup = keyboards.ShowMessageKeyboard(msg.ID)
			go l.bot.Send(otherUserSend)
			done := tgbotapi.NewMessage(u.ID64(), `👆 پیام بالا برای مخاطبتون ارسال شد.
هر وقت پیام رو ببینه بهت اطلاع میدم.`)
			done.ReplyMarkup = keyboards.HomeKeyboard()
			l.bot.Send(done)

			// finish message sending ids
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

func (l LocationSendAnnmsg) FinishSendMessage(u *models.UserManager, up *tgbotapi.Update) {
	c := u.GetCache("annmsg_id")
	if c == "" {
		u.Error("⚠️ مشکلی پیش آمد! لینک ناشناس اشتباه میباشد", l.bot)
		return
	}
	id, _ := strconv.Atoi(c)
	otherUser := u.GetUserBy("db", id)
	msgIds := u.GetWaitingMsgs()
	sendableMsg := &models.Annmsg{
		Type:                "GROUP",
		Data:                "",
		Caption:             "",
		FromId:              u.UserMessage.DatabaseID,
		ToId:                otherUser.DatabaseID,
		ID:                  -1,
		ReplyTo:             -1,
		Status:              0,
		SenderMessageID:     0,
		RecieverMessageID:   0,
		BotPreviewMessageID: 0,
		Group:               msgIds,
	}
	dbannmsg, err := sendableMsg.Save(u.DB, u.Cache)
	if err != nil {
		u.Error("مشکلی پیش آمد !!", l.bot)
		return
	}
	sendableMsg.ID = int(dbannmsg.ID)
	sendableMsg.SaveCache(u.Cache)

	// message is created. sending created message

	otherUserSend := tgbotapi.NewMessage(int64(otherUser.TelegramID), `💌 یک پیام جدید دریافت کردید!
برای نمایش پیام روی دکمه نمایش کلیک کنید.`)
	otherUserSend.ReplyMarkup = keyboards.ShowMessageKeyboard(sendableMsg.ID)
	_, e := l.bot.Send(otherUserSend)
	if e != nil {
		u.Error(`متاسفانه این کاربر ربات رو خاموش کرده :(`, l.bot)
		u.ClearCache()
		u.DelDeletableMsgs(l.bot)
		return
	}

	u.AddDeletableMsg(up.Message.MessageID)
	msgs := sendableMsg.Msgs()
	for _, msgId := range msgs {
		msg := models.GetMessage(msgId, u.ContextModel)
		l.bot.Send(tgbotapi.NewDeleteMessage(u.ID64(), msg.BotPreviewMessageID))
		u.UnsetFromWaitingMsgs(msg.ID)
	}
	u.DelDeletableMsgs(l.bot)
	sendedMsg := tgbotapi.NewMessage(u.ID64(), `✅ با موفقیت ارسال شد. اگر مخاطبتون پیام رو ببینه بهت اطلاع میدم

چه کاری برات انجام بدم؟`)
	sendedMsg.ReplyMarkup = keyboards.HomeKeyboard()
	r, e := l.bot.Send(sendedMsg)
	if e == nil {
		sendableMsg.SenderMessageID = r.MessageID
	}
	u.ClearCache()

}

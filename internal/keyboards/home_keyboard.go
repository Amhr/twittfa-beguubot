package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HomeKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(TXT_MY_LINK),
		),
	)
}

func CancelKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(TXT_CANCEL),
		),
	)
}

func ShowMessageKeyboard(msgId int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(TXT_OPEN, fmt.Sprintf("open-%d", msgId)),
		),
	)
}

func MessageDetailKeyboard(msgId int, isBlocked bool) tgbotapi.InlineKeyboardMarkup {
	txt := "🔓بلاک کردن"
	if isBlocked {
		txt = "🔐‌آزاد کردن"
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✍️ پاسخ دادن", fmt.Sprintf("reply-%d", msgId)),
			tgbotapi.NewInlineKeyboardButtonData(txt, fmt.Sprintf("block-%d", msgId)),
		),
	)
}

func FinishSendMessageKeyboard(msgId int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✔️ ارسال", fmt.Sprintf("send-%d", msgId)),
			tgbotapi.NewInlineKeyboardButtonData("❌ انصراف", fmt.Sprintf("delete-%d", msgId)),
		),
	)
}

package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *BeguuBot) Worker(id int32, updates chan tgbotapi.Update) {
	for update := range updates {
		if update.Message != nil {
			b.HandleMessage(&update)
		} else if update.CallbackQuery != nil {
			b.HandleCallback(&update)
		}
	}
}

package internal

func (b *BeguuBot) Worker(id int32) {
	for update := range b.Updates {
		if update.Message != nil {
			b.HandleMessage(update)
		} else if update.CallbackQuery != nil {
			b.HandleCallback(update)
		}
	}
}

package models

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UserMessage struct {
	FirstName  string
	Username   string
	TelegramID int
	DatabaseID int
}

func (m *UserMessage) Byte() []byte {
	b, _ := json.Marshal(m)
	return b
}

func UserMessageFromDBUser(dbuser *DBUser) *UserMessage {
	return &UserMessage{
		FirstName:  dbuser.FirstName,
		Username:   dbuser.Username,
		TelegramID: dbuser.TelegramID,
		DatabaseID: int(dbuser.ID),
	}
}

func UserMessageFromTelegramUser(dbuser *tgbotapi.User, DbId int) *UserMessage {
	return &UserMessage{
		FirstName:  dbuser.FirstName,
		Username:   dbuser.UserName,
		TelegramID: dbuser.ID,
		DatabaseID: DbId,
	}
}

func UserMessageFromByte(b []byte) *UserMessage {
	um := &UserMessage{}
	json.Unmarshal(b, um)
	return um
}

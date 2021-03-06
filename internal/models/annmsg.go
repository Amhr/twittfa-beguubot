package models

import (
	"context"
	"encoding/json"
	context2 "github.com/amhr/begubot/internal/context"
	"github.com/amhr/begubot/internal/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Annmsg struct {
	Type                string
	Data                string
	Caption             string
	FromId              int
	ToId                int
	ID                  int
	ReplyTo             int
	Status              int
	SenderMessageID     int
	RecieverMessageID   int
	BotPreviewMessageID int
	Group               []int
	GroupRaw            string
}

type DBAnnmsg struct {
	gorm.Model
	Type    string
	Data    string
	Caption string
	ReplyTo int
	FromId  int    `gorm:"index"`
	ToId    int    `gorm:"index"`
	Status  int    `gorm:"default:0;type:int(1)"`
	Group   string `gorm:"default:null"`
}

func (d *DBAnnmsg) ToMessage() *Annmsg {
	return &Annmsg{
		Type:     d.Type,
		Data:     d.Data,
		Caption:  d.Caption,
		FromId:   d.FromId,
		ToId:     d.ToId,
		ID:       int(d.ID),
		ReplyTo:  d.ReplyTo,
		GroupRaw: d.Group,
	}
}

func (m *Annmsg) Cancel(u *UserManager) {
	if m.Status != 0 {
		go u.ContextModel.Bot.Send(tgbotapi.NewDeleteMessage(u.ID64(), m.SenderMessageID))
		go u.ContextModel.Bot.Send(tgbotapi.NewDeleteMessage(u.ID64(), m.BotPreviewMessageID))
	}
}

func (m *Annmsg) SaveCache(c *redis.RedisCache) error {
	key := c.Key("annmsg", strconv.Itoa(m.ID))
	b, err := json.Marshal(m)
	if err != nil {
		logrus.WithField("action", "SaveAnnmsgMarshal").Error(err)
		return err
	}
	c.Set(key, string(b), time.Duration(24*7)*time.Hour, context.Background())
	return nil
}

func (m *Annmsg) Msgs() []int {
	if len(m.Group) != 0 {
		return m.Group
	}
	var ids []int
	if e := json.Unmarshal([]byte(m.GroupRaw), &ids); e != nil {
		return []int{}
	} else {
		return ids
	}
}

func NewAnnmsg(t, d, c string, fromId, toId int, replyTo int) *Annmsg {
	return &Annmsg{
		Type:    t,
		Data:    d,
		Caption: c,
		FromId:  fromId,
		ToId:    toId,
		ID:      0,
		ReplyTo: replyTo,
	}
}

func GetMessage(msgId int, c *context2.ModelContext) *Annmsg {
	annmsg := &Annmsg{}
	if msgId == 0 {
		return annmsg
	}
	cacheKey := c.Redis.Key("annmsg", strconv.Itoa(msgId))
	d := c.Redis.Get(cacheKey, "", context.Background())
	if e := json.Unmarshal([]byte(d), annmsg); e == nil {
		return annmsg
	} else {
		dbannmsg := &DBAnnmsg{}
		e := c.DB.Where("id=?", msgId).Take(dbannmsg)
		if e.Error != nil {
			return nil
		}
		annmsg = dbannmsg.ToMessage()
		annmsg.SaveCache(c.Redis)
		return annmsg
	}

}

func (annmsg *Annmsg) Save(db *gorm.DB, c *redis.RedisCache) (*DBAnnmsg, error) {
	groupItems, _ := json.Marshal(annmsg.Group)
	dbAnnmsg := &DBAnnmsg{
		Type:    annmsg.Type,
		Data:    annmsg.Data,
		Caption: annmsg.Caption,
		FromId:  annmsg.FromId,
		ToId:    annmsg.ToId,
		ReplyTo: annmsg.ReplyTo,
		Group:   string(groupItems),
	}
	t := db.Save(dbAnnmsg)
	if t.Error != nil {
		logrus.WithField("action", "SaveAnnmsg").Error(t.Error)
		return nil, t.Error
	} else {
		annmsg.ID = int(dbAnnmsg.ID)
		if e := annmsg.SaveCache(c); e != nil {
			return nil, e
		}
		return dbAnnmsg, nil
	}
}

func ConvertUpdateToAnnmsg(u *tgbotapi.Update) *Annmsg {
	if u.Message.Text != "" {
		return &Annmsg{
			Type:    "Text",
			Data:    u.Message.Text,
			Caption: "",
		}
	}
	if u.Message.Photo != nil {
		photos := u.Message.Photo
		return &Annmsg{
			Type:    "Photo",
			Data:    (*photos)[0].FileID,
			Caption: u.Message.Caption,
		}
	}
	if u.Message.Video != nil {
		return &Annmsg{
			Type:    "Video",
			Data:    u.Message.Video.FileID,
			Caption: u.Message.Caption,
		}
	}

	if u.Message.VideoNote != nil {
		return &Annmsg{
			Type:    "VideoNote",
			Data:    u.Message.VideoNote.FileID,
			Caption: u.Message.Caption,
		}
	}

	if u.Message.Voice != nil {
		return &Annmsg{
			Type:    "Voice",
			Data:    u.Message.Voice.FileID,
			Caption: u.Message.Caption,
		}
	}

	if u.Message.Audio != nil {
		return &Annmsg{
			Type:    "Audio",
			Data:    u.Message.Audio.FileID,
			Caption: u.Message.Caption,
		}
	}

	if u.Message.Document != nil {
		return &Annmsg{
			Type:    "Document",
			Data:    u.Message.Document.FileID,
			Caption: u.Message.Caption,
		}
	}

	if u.Message.Sticker != nil {
		return &Annmsg{
			Type:    "Sticker",
			Data:    u.Message.Sticker.FileID,
			Caption: u.Message.Caption,
		}
	}

	return nil
}

func SendMessage(msg *Annmsg, to int64, replyMarkup *tgbotapi.InlineKeyboardMarkup, replyTo int) tgbotapi.Chattable {
	switch msg.Type {
	case "Sticker":
		a := tgbotapi.NewStickerShare(to, msg.Data)
		if replyMarkup != nil {
			a.ReplyMarkup = replyMarkup
		}
		if replyTo > 0 {
			a.ReplyToMessageID = replyTo
		}
		return a
	case "Photo":
		ph := tgbotapi.NewPhotoShare(to, msg.Data)
		ph.Caption = msg.Caption
		if replyMarkup != nil {
			ph.ReplyMarkup = replyMarkup
		}
		if replyTo > 0 {
			ph.ReplyToMessageID = replyTo
		}
		return ph
	case "Video":
		ph := tgbotapi.NewVideoShare(to, msg.Data)
		ph.Caption = msg.Caption
		if replyMarkup != nil {
			ph.ReplyMarkup = replyMarkup
		}
		if replyTo > 0 {
			ph.ReplyToMessageID = replyTo
		}
		return ph
	case "VideoNote":
		ph := tgbotapi.NewVideoNoteShare(to, 0, msg.Data)
		if replyMarkup != nil {
			ph.ReplyMarkup = replyMarkup
		}
		if replyTo > 0 {
			ph.ReplyToMessageID = replyTo
		}
		return ph
	case "Document":
		ph := tgbotapi.NewDocumentShare(to, msg.Data)
		ph.Caption = msg.Caption
		if replyMarkup != nil {
			ph.ReplyMarkup = replyMarkup
		}
		if replyTo > 0 {
			ph.ReplyToMessageID = replyTo
		}
		return ph

	case "Audio":
		ph := tgbotapi.NewAudioShare(to, msg.Data)
		ph.Caption = msg.Caption
		if replyMarkup != nil {
			ph.ReplyMarkup = replyMarkup
		}
		if replyTo > 0 {
			ph.ReplyToMessageID = replyTo
		}
		return ph
	case "Voice":
		ph := tgbotapi.NewVoiceShare(to, msg.Data)
		ph.Caption = msg.Caption
		if replyMarkup != nil {
			ph.ReplyMarkup = replyMarkup
		}
		if replyTo > 0 {
			ph.ReplyToMessageID = replyTo
		}
		return ph
	case "Text":
		ph := tgbotapi.NewMessage(to, msg.Data)
		if replyMarkup != nil {
			ph.ReplyMarkup = replyMarkup
		}
		if replyTo > 0 {
			ph.ReplyToMessageID = replyTo
		}
		return ph
	}
	return nil
}

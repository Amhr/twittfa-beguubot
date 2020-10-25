package models

import (
	"context"
	"encoding/json"
	"fmt"
	context2 "github.com/amhr/begubot/internal/context"
	"github.com/amhr/begubot/internal/epimetheus"
	"github.com/amhr/begubot/internal/keyboards"
	"github.com/amhr/begubot/internal/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type DBUser struct {
	gorm.Model
	FirstName  string
	LastName   string
	Username   string `gorm:"index"`
	TelegramID int    `gorm:"index"`
}

func (u *DBUser) Update(usr *UserManager) {
	u.FirstName = usr.TgBotUser.FirstName
	u.LastName = usr.TgBotUser.LastName
	u.Username = strings.ToLower(usr.TgBotUser.UserName)
	usr.DB.Save(u)
}

type UserManager struct {
	TgBotUser    *tgbotapi.User
	Cache        *redis.RedisCache
	DB           *gorm.DB
	Metrics      *epimetheus.MetricsManager
	UserMessage  *UserMessage
	ContextModel *context2.ModelContext
}

func NewUser(tgu *tgbotapi.User, c *redis.RedisCache, db *gorm.DB, m *epimetheus.MetricsManager, cm *context2.ModelContext) *UserManager {
	return &UserManager{
		TgBotUser:    tgu,
		Cache:        c,
		DB:           db,
		Metrics:      m,
		ContextModel: cm,
	}
}

func (u *UserManager) Key(items ...string) string {
	items = append(items, strconv.Itoa(u.UserMessage.TelegramID))
	return u.Cache.Key(items...)
}

func (u *UserManager) GetLocation() string {
	return u.Cache.Get(u.Key("location"), "home", context.Background())
}
func (u *UserManager) SetLocation(loc string) {
	u.Cache.Set(u.Key("location"), loc, time.Duration(24*7)*time.Hour, context.Background())
}

func (u *UserManager) ClearCache() {
	u.SetLocation("home")
	u.SetStep("1")
	u.SetCache("annmsg_id", "")
	u.SetCache("annmsg_reply", "")
}

func (u *UserManager) GetWaitingMsgs() []int {
	d := u.GetCache("waitingmsgs")
	var msgs []int
	if e := json.Unmarshal([]byte(d), &msgs); e == nil {
		return msgs
	} else {
		return []int{}
	}
}

func (u *UserManager) SetWaitingMsgs(msgs []int) {
	b, e := json.Marshal(msgs)
	if e != nil {
		return
	}
	u.SetCache("waitingmsgs", string(b))
}

func (u *UserManager) UnsetFromWaitingMsgs(msgId int) {
	ids := u.GetWaitingMsgs()
	nIds := make([]int, 0)
	for _, id := range ids {
		if id != msgId {
			nIds = append(nIds, id)
		}
	}
	u.SetWaitingMsgs(nIds)
}

func (u *UserManager) AddWaitingMsg(id int) {
	ids := u.GetWaitingMsgs()
	ids = append(ids, id)
	u.SetWaitingMsgs(ids)
}

func (u *UserManager) GetStep() string {
	return u.Cache.Get(u.Key("step"), "home", context.Background())
}
func (u *UserManager) SetStep(loc string) {
	u.Cache.Set(u.Key("step"), loc, time.Duration(24*7)*time.Hour, context.Background())
}

func (u *UserManager) GetCache(key string) string {
	return u.Cache.Get(u.Key("cache", key), "home", context.Background())
}
func (u *UserManager) SetCache(key string, val string) {
	u.Cache.Set(u.Key("cache", key), val, time.Duration(1)*time.Hour, context.Background())
}

func (u *UserManager) ID64() int64 {
	return int64(u.UserMessage.TelegramID)
}

func (u *UserManager) MyLinkIdentifier() string {
	return fmt.Sprintf("%d_%s", u.UserMessage.DatabaseID, randstr.Hex(4))
}

func (u *UserManager) UserByTID(teleid int) *DBUser {
	us := &DBUser{}
	f := u.DB.Where("telegram_id=?", teleid).First(us)
	if f.Error != nil {
		return nil
	}
	return us
}

func (u *UserManager) UserByID(id int) *DBUser {
	us := &DBUser{}
	f := u.DB.Where("id=?", id).First(us)
	if f.Error != nil {
		return nil
	}
	return us
}

func (u *UserManager) Error(txt string, bot *tgbotapi.BotAPI) {
	c := tgbotapi.NewMessage(u.ID64(), txt)
	c.ReplyMarkup = keyboards.HomeKeyboard()
	bot.Send(c)
}

func (u *UserManager) SaveDB() uint {
	usr := u.UserByTID(u.TgBotUser.ID)
	dbuid := uint(0)
	if usr == nil {
		us := &DBUser{
			FirstName:  u.TgBotUser.FirstName,
			LastName:   u.TgBotUser.LastName,
			Username:   strings.ToLower(u.TgBotUser.UserName),
			TelegramID: u.TgBotUser.ID,
		}
		u.DB.Save(us)
		dbuid = us.ID
	} else {
		usr.Update(u)
		dbuid = usr.ID
	}
	return dbuid
}

func (u *UserManager) SaveUserCache(us *UserMessage) {
	u.Cache.Set(u.Cache.Key("user", "tg", strconv.Itoa(us.TelegramID)), string(us.Byte()), time.Duration(24)*time.Hour, context.Background())
	u.Cache.Set(u.Cache.Key("user", "db", strconv.Itoa(us.DatabaseID)), string(us.Byte()), time.Duration(24)*time.Hour, context.Background())
}

func (u *UserManager) UserFromCache(t string, id int) *UserMessage {
	usr := u.Cache.Get(u.Cache.Key("user", t, strconv.Itoa(id)), "", context.Background())
	if usr == "" {
		return nil
	}
	nusr := &UserMessage{}
	err := json.Unmarshal([]byte(usr), nusr)
	if err != nil {
		return nil
	}
	return nusr
}

func (u *UserManager) MakeUserCache(DbId int) *UserMessage {
	return &UserMessage{
		FirstName:  u.TgBotUser.FirstName,
		Username:   strings.ToLower(u.TgBotUser.UserName),
		TelegramID: u.TgBotUser.ID,
		DatabaseID: DbId,
	}
}

//todo: implement better user cache system

func (u *UserManager) SaveCache() {
	u.SaveUserCache(u.UserMessage)
}

// searchs in cache too
func (u *UserManager) GetUserBy(t string, id int) *UserMessage {
	usrMsg := u.UserFromCache(t, id)
	// load from db
	if usrMsg == nil {
		u.Metrics.CacheNotExists("LoadUser")
		if t == "tid" {
			dbUser := u.UserByTID(id)
			if dbUser == nil {
				return nil
			}
			usrMsg = UserMessageFromDBUser(dbUser)
		} else {
			dbUser := u.UserByID(id)
			if dbUser == nil {
				return nil
			}
			usrMsg = UserMessageFromDBUser(dbUser)
		}
	}
	u.SaveUserCache(usrMsg)
	return usrMsg
}

func (u *UserManager) Load() {
	usrMsg := u.UserFromCache("tg", u.TgBotUser.ID)
	// load from db
	if usrMsg == nil {
		u.Metrics.CacheNotExists("LoadUser")
		dbid := u.SaveDB()
		u.UserMessage = UserMessageFromTelegramUser(u.TgBotUser, int(dbid))
	} else {
		u.Metrics.CacheExists("LoadUser")
		u.UserMessage = usrMsg
	}
	u.SaveUserCache(u.UserMessage)
}

// blocking

func (u *UserManager) Block(dbId int) {
	k := fmt.Sprintf("block:%d:%d", u.UserMessage.DatabaseID, dbId)
	u.Cache.Set(k, "yes", time.Hour*time.Duration(24*30), context.Background())
}

func (u *UserManager) Unblock(dbId int) {
	k := fmt.Sprintf("block:%d:%d", u.UserMessage.DatabaseID, dbId)
	u.Cache.Set(k, "", time.Hour*time.Duration(24*30), context.Background())
}
func (u *UserManager) IsBlocked(dbId int) bool {
	k := fmt.Sprintf("block:%d:%d", u.UserMessage.DatabaseID, dbId)
	return u.Cache.Get(k, "", context.Background()) != ""
}
func (u *UserManager) ImBlocked(dbId int) bool {
	k := fmt.Sprintf("block:%d:%d", dbId, u.UserMessage.DatabaseID)
	return u.Cache.Get(k, "", context.Background()) != ""
}

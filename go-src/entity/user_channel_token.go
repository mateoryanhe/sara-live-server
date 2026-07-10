package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbUserChannelToken db.TbName = "user_channel_tokens"
)

const (
	UserChannelTokenUserId      db.TbCol = "user_id"
	UserChannelTokenToken       db.TbCol = "token"
	UserChannelTokenExpireAt    db.TbCol = "expire_at"
	UserChannelTokenChannelName db.TbCol = "channel_name"
)

// UserChannelToken 用户频道Token(独立表,与通话业务无关)
// 主键 userId,不重新生成
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserChannelToken struct {
	UserId      uint64    `gorm:"primaryKey;column:user_id;comment:用户ID" json:"userId"`
	Token       string    `gorm:"size:512;default:'';comment:Token字符串" json:"token"`
	ExpireAt    time.Time `gorm:"index;comment:过期时间" json:"expireAt"`
	ChannelName string    `gorm:"size:128;default:'';comment:频道名称" json:"channelName"`
}

func NewUserChannelToken(userId uint64, channelName, token string, expireAt time.Time) *UserChannelToken {
	ret := &UserChannelToken{}
	ret.UserId = userId
	ret.SetChannelName(channelName)
	ret.SetToken(token)
	ret.SetExpireAt(expireAt)
	return ret
}

func (m *UserChannelToken) SetToken(v string) {
	m.Token = v
	syndb.AddDataToQuickChan(TbUserChannelToken, UserChannelTokenToken, &syndb.ColData{IdVal: m.UserId, ColVal: v})
}

func (m *UserChannelToken) SetExpireAt(v time.Time) {
	m.ExpireAt = v
	syndb.AddDataToQuickChan(TbUserChannelToken, UserChannelTokenExpireAt, &syndb.ColData{IdVal: m.UserId, ColVal: v})
}

func (m *UserChannelToken) SetChannelName(v string) {
	m.ChannelName = v
	syndb.AddDataToQuickChan(TbUserChannelToken, UserChannelTokenChannelName, &syndb.ColData{IdVal: m.UserId, ColVal: v})
}

func initUserChannelToken() {
	syndb.RegQuick(TbUserChannelToken, UserChannelTokenToken)
	syndb.RegQuick(TbUserChannelToken, UserChannelTokenExpireAt)
	syndb.RegQuick(TbUserChannelToken, UserChannelTokenChannelName)
	migrate.AutoMigrate(&UserChannelToken{})
}

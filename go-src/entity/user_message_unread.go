package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbUserMessageUnread db.TbName = "user_message_unreads"
)

const (
	UserMessageUnreadSystemUnread  db.TbCol = "system_unread"
	UserMessageUnreadPrivateUnread db.TbCol = "private_unread"
)

// UserMessageUnread 用户消息未读数(与用户一一对应,主键ID即用户ID)
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserMessageUnread struct {
	migrate.OneModel
	SystemUnread  uint64 `gorm:"default:0;comment:系统消息未读数" json:"systemUnread"`
	PrivateUnread uint64 `gorm:"default:0;comment:私信未读数" json:"privateUnread"`
}

func NewUserMessageUnread(userId uint64) *UserMessageUnread {
	ret := &UserMessageUnread{}
	ret.ID = userId
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

func (m *UserMessageUnread) AddSystemUnread(v uint64) {
	m.SystemUnread = math.Add(m.SystemUnread, v)
	m.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserMessageUnread, UserMessageUnreadSystemUnread, &syndb.ColData{
		IdVal: m.ID, ColVal: m.SystemUnread,
	})
}

func (m *UserMessageUnread) SubSystemUnread(v uint64) {
	if v >= m.SystemUnread {
		m.SystemUnread = 0
	} else {
		m.SystemUnread -= v
	}
	m.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserMessageUnread, UserMessageUnreadSystemUnread, &syndb.ColData{
		IdVal: m.ID, ColVal: m.SystemUnread,
	})
}

func (m *UserMessageUnread) AddPrivateUnread(v uint64) {
	m.PrivateUnread = math.Add(m.SystemUnread, v)
	m.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserMessageUnread, UserMessageUnreadPrivateUnread, &syndb.ColData{
		IdVal: m.ID, ColVal: m.PrivateUnread,
	})
}

func (m *UserMessageUnread) SubPrivateUnread(v uint64) {
	m.PrivateUnread = math.Sub(m.SystemUnread, v)
	syndb.AddDataToQuickChan(TbUserMessageUnread, UserMessageUnreadPrivateUnread, &syndb.ColData{
		IdVal: m.ID, ColVal: m.PrivateUnread,
	})
}

func (m *UserMessageUnread) SetCreatedAt(v time.Time) {
	m.CreatedAt = v
	syndb.AddDataToQuickChan(TbUserMessageUnread, db.CreatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserMessageUnread) SetUpdatedAt(v time.Time) {
	m.UpdatedAt = v
	syndb.AddDataToQuickChan(TbUserMessageUnread, db.UpdatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func initUserMessageUnread() {
	syndb.RegQuickWithLarge(TbUserMessageUnread, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbUserMessageUnread, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbUserMessageUnread, UserMessageUnreadSystemUnread)
	syndb.RegQuickWithLarge(TbUserMessageUnread, UserMessageUnreadPrivateUnread)
	migrate.AutoMigrate(&UserMessageUnread{})
}

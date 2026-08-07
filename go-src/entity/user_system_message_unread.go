package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbUserSystemMessageUnread db.TbName = "user_system_message_unreads"
)

const (
	UserSystemMessageUnreadUserId      db.TbCol = "user_id"
	UserSystemMessageUnreadType        db.TbCol = "type"
	UserSystemMessageUnreadUnreadCount db.TbCol = "unread_count"
)

const (
	UserSystemMessageUnreadTypeActivity uint8 = 1 // 活动消息
	UserSystemMessageUnreadTypePersonal uint8 = 2 // 个人系统消息
)

// UserSystemMessageUnread 系统消息未读明细(按用户+类型维度)
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type UserSystemMessageUnread struct {
	ID          string    `gorm:"primaryKey;size:64;comment:复合ID(userId_type)" json:"id"`
	UserId      uint64    `gorm:"default:0;comment:用户id" json:"userId"`
	Type        uint8     `gorm:"default:0;comment:类型(1活动消息,2个人系统消息)" json:"type"`
	UnreadCount uint64    `gorm:"default:0;comment:未读数量" json:"unreadCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func BuildUserSystemMessageUnreadId(userId uint64, msgType uint8) string {
	return fmt.Sprintf("%d_%d", userId, msgType)
}

func NewUserSystemMessageUnread(userId uint64, msgType uint8) *UserSystemMessageUnread {
	row := &UserSystemMessageUnread{}
	row.ID = BuildUserSystemMessageUnreadId(userId, msgType)
	row.SetUserId(userId)
	row.SetType(msgType)
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	return row
}

func (m *UserSystemMessageUnread) SetUserId(v uint64) {
	m.UserId = v
	syndb.AddData(TbUserSystemMessageUnread, UserSystemMessageUnreadUserId, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserSystemMessageUnread) SetType(v uint8) {
	m.Type = v
	syndb.AddData(TbUserSystemMessageUnread, UserSystemMessageUnreadType, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserSystemMessageUnread) AddUnread(v uint64) {
	m.UnreadCount = math.Add(m.UnreadCount, v)
	m.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserSystemMessageUnread, UserSystemMessageUnreadUnreadCount, &syndb.ColData{
		IdVal: m.ID, ColVal: m.UnreadCount,
	})
}

func (m *UserSystemMessageUnread) SubUnread(v uint64) {
	if v >= m.UnreadCount {
		m.UnreadCount = 0
	} else {
		m.UnreadCount -= v
	}
	m.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserSystemMessageUnread, UserSystemMessageUnreadUnreadCount, &syndb.ColData{
		IdVal: m.ID, ColVal: m.UnreadCount,
	})
}

func (m *UserSystemMessageUnread) SetCreatedAt(v time.Time) {
	m.CreatedAt = v
	syndb.AddData(TbUserSystemMessageUnread, db.CreatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func (m *UserSystemMessageUnread) SetUpdatedAt(v time.Time) {
	m.UpdatedAt = v
	syndb.AddData(TbUserSystemMessageUnread, db.UpdatedAtName, &syndb.ColData{
		IdVal: m.ID, ColVal: v,
	})
}

func initUserSystemMessageUnread() {
	syndb.RegQuick(TbUserSystemMessageUnread, db.CreatedAtName)
	syndb.RegQuick(TbUserSystemMessageUnread, db.UpdatedAtName)
	syndb.RegQuick(TbUserSystemMessageUnread, UserSystemMessageUnreadUserId)
	syndb.RegQuick(TbUserSystemMessageUnread, UserSystemMessageUnreadType)
	syndb.RegQuick(TbUserSystemMessageUnread, UserSystemMessageUnreadUnreadCount)
	migrate.AutoMigrate(&UserSystemMessageUnread{})
}

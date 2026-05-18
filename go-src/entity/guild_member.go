package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveGuildMember db.TbName = "live_guild_members"
)

const (
	GuildMemberGuildId db.TbCol = "guild_id"
	GuildMemberUserId  db.TbCol = "user_id"
	GuildMemberStatus  db.TbCol = "status"
)

// 工会成员状态
const (
	GuildMemberStatusPending  uint8 = 0 // 申请中
	GuildMemberStatusApproved uint8 = 1 // 已加入
	GuildMemberStatusRejected uint8 = 2 // 已拒绝
	GuildMemberStatusKicked   uint8 = 3 // 已剔除
)

// LiveGuildMember 直播工会成员
type LiveGuildMember struct {
	migrate.OneModel
	GuildId uint64 `gorm:"index;default:0;comment:工会ID" json:"guildId"`
	UserId  uint64 `gorm:"index;default:0;comment:用户ID" json:"userId"`
	Status  uint8  `gorm:"default:0;comment:状态(0申请中,1已加入,2已拒绝,3已剔除)" json:"status"`
}

// NewLiveGuildMember 创建新的工会成员对象(只构造内存对象,通过syndb异步入库)
func NewLiveGuildMember(memberId, guildId, userId uint64) *LiveGuildMember {
	m := &LiveGuildMember{}
	m.ID = memberId
	now := time.Now()
	m.SetCreatedAt(now)
	m.SetUpdatedAt(now)
	m.SetGuildId(guildId)
	m.SetUserId(userId)
	m.SetStatus(GuildMemberStatusPending)
	return m
}

func (m *LiveGuildMember) SetGuildId(v uint64) {
	m.GuildId = v
	syndb.AddDataToQuickChan(TbLiveGuildMember, GuildMemberGuildId, &syndb.ColData{
		IdVal:  m.ID,
		ColVal: v,
	})
}

func (m *LiveGuildMember) SetUserId(v uint64) {
	m.UserId = v
	syndb.AddDataToQuickChan(TbLiveGuildMember, GuildMemberUserId, &syndb.ColData{
		IdVal:  m.ID,
		ColVal: v,
	})
}

func (m *LiveGuildMember) SetStatus(v uint8) {
	m.Status = v
	m.SetUpdatedAtOnly(time.Now())
	syndb.AddDataToQuickChan(TbLiveGuildMember, GuildMemberStatus, &syndb.ColData{
		IdVal:  m.ID,
		ColVal: v,
	})
}

func (m *LiveGuildMember) SetCreatedAt(v time.Time) {
	m.CreatedAt = v
	syndb.AddDataToQuickChan(TbLiveGuildMember, db.CreatedAtName, &syndb.ColData{
		IdVal:  m.ID,
		ColVal: v,
	})
}

func (m *LiveGuildMember) SetUpdatedAt(v time.Time) {
	m.UpdatedAt = v
	syndb.AddDataToQuickChan(TbLiveGuildMember, db.UpdatedAtName, &syndb.ColData{
		IdVal:  m.ID,
		ColVal: v,
	})
}

// SetUpdatedAtOnly 仅刷新内存中的UpdatedAt并入队同步,避免在SetStatus中重复调用产生多余日志
func (m *LiveGuildMember) SetUpdatedAtOnly(v time.Time) {
	m.UpdatedAt = v
	syndb.AddDataToQuickChan(TbLiveGuildMember, db.UpdatedAtName, &syndb.ColData{
		IdVal:  m.ID,
		ColVal: v,
	})
}

func initGuildMember() {
	syndb.RegQuickWithLarge(TbLiveGuildMember, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbLiveGuildMember, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbLiveGuildMember, GuildMemberGuildId)
	syndb.RegQuickWithLarge(TbLiveGuildMember, GuildMemberUserId)
	syndb.RegQuickWithLarge(TbLiveGuildMember, GuildMemberStatus)
	migrate.AutoMigrate(&LiveGuildMember{})
}

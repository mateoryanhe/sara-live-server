package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveGuild db.TbName = "live_guilds"
)

const (
	LiveGuildName        db.TbCol = "name"
	LiveGuildLeaderId    db.TbCol = "leader_id"
	LiveGuildLeaderName  db.TbCol = "leader_name"
	LiveGuildBankCard    db.TbCol = "bank_card"
	LiveGuildContact     db.TbCol = "contact"
	LiveGuildDescription db.TbCol = "description"
	LiveGuildStatus      db.TbCol = "status"
)

// LiveGuild 直播工会(启动预加载至内存,永不过期;字段写入通过 syndb 异步入库,CMS 列表直接读缓存)
type LiveGuild struct {
	migrate.OneModel
	Name        string `gorm:"size:64;comment:工会名称" json:"name"`
	LeaderId    uint64 `gorm:"default:0;comment:会长/负责人ID" json:"leaderId"`
	LeaderName  string `gorm:"size:64;default:'';comment:会长名称" json:"leaderName"`
	BankCard    string `gorm:"size:64;default:'';comment:银行卡" json:"bankCard"`
	Contact     string `gorm:"size:64;comment:联系方式" json:"contact"`
	Description string `gorm:"size:255;comment:工会简介" json:"description"`
	Status      uint8  `gorm:"default:1;comment:状态(0-禁用,1-启用)" json:"status"`
}

// NewLiveGuild 构造内存对象,字段写入通过 syndb 异步入库
func NewLiveGuild(id uint64, name string, leaderId uint64, leaderName, contact, description string, status uint8) *LiveGuild {
	g := &LiveGuild{}
	g.ID = id
	now := time.Now()
	g.SetCreatedAt(now)
	g.SetUpdatedAt(now)
	g.SetName(name)
	g.SetLeaderId(leaderId)
	g.SetLeaderName(leaderName)
	g.SetBankCard("")
	g.SetContact(contact)
	g.SetDescription(description)
	if status == 0 {
		status = 1
	}
	g.SetStatus(status)
	return g
}

func (g *LiveGuild) SetName(v string) {
	g.Name = v
	g.touchUpdatedAt()
	syndb.AddData(TbLiveGuild, LiveGuildName, &syndb.ColData{
		IdVal: g.ID, ColVal: v,
	})
}

func (g *LiveGuild) SetLeaderId(v uint64) {
	g.LeaderId = v
	g.touchUpdatedAt()
	syndb.AddData(TbLiveGuild, LiveGuildLeaderId, &syndb.ColData{
		IdVal: g.ID, ColVal: v,
	})
}

func (g *LiveGuild) SetLeaderName(v string) {
	g.LeaderName = v
	g.touchUpdatedAt()
	syndb.AddData(TbLiveGuild, LiveGuildLeaderName, &syndb.ColData{
		IdVal: g.ID, ColVal: v,
	})
}

func (g *LiveGuild) SetBankCard(v string) {
	g.BankCard = v
	g.touchUpdatedAt()
	syndb.AddData(TbLiveGuild, LiveGuildBankCard, &syndb.ColData{
		IdVal: g.ID, ColVal: v,
	})
}

func (g *LiveGuild) SetContact(v string) {
	g.Contact = v
	g.touchUpdatedAt()
	syndb.AddData(TbLiveGuild, LiveGuildContact, &syndb.ColData{
		IdVal: g.ID, ColVal: v,
	})
}

func (g *LiveGuild) SetDescription(v string) {
	g.Description = v
	g.touchUpdatedAt()
	syndb.AddData(TbLiveGuild, LiveGuildDescription, &syndb.ColData{
		IdVal: g.ID, ColVal: v,
	})
}

func (g *LiveGuild) SetStatus(v uint8) {
	g.Status = v
	g.touchUpdatedAt()
	syndb.AddData(TbLiveGuild, LiveGuildStatus, &syndb.ColData{
		IdVal: g.ID, ColVal: v,
	})
}

func (g *LiveGuild) SetCreatedAt(v time.Time) {
	g.CreatedAt = v
	syndb.AddData(TbLiveGuild, db.CreatedAtName, &syndb.ColData{
		IdVal: g.ID, ColVal: v,
	})
}

func (g *LiveGuild) SetUpdatedAt(v time.Time) {
	g.UpdatedAt = v
	syndb.AddData(TbLiveGuild, db.UpdatedAtName, &syndb.ColData{
		IdVal: g.ID, ColVal: v,
	})
}

func (g *LiveGuild) touchUpdatedAt() {
	g.UpdatedAt = time.Now()
	syndb.AddData(TbLiveGuild, db.UpdatedAtName, &syndb.ColData{
		IdVal: g.ID, ColVal: g.UpdatedAt,
	})
}

func InitLiveGuild() {
	initLiveGuild()
}

func initLiveGuild() {
	syndb.RegQuick(TbLiveGuild, db.CreatedAtName)
	syndb.RegQuick(TbLiveGuild, db.UpdatedAtName)
	syndb.RegQuick(TbLiveGuild, LiveGuildName)
	syndb.RegQuick(TbLiveGuild, LiveGuildLeaderId)
	syndb.RegQuick(TbLiveGuild, LiveGuildLeaderName)
	syndb.RegQuick(TbLiveGuild, LiveGuildBankCard)
	syndb.RegQuick(TbLiveGuild, LiveGuildContact)
	syndb.RegQuick(TbLiveGuild, LiveGuildDescription)
	syndb.RegQuick(TbLiveGuild, LiveGuildStatus)
	migrate.AutoMigrate(&LiveGuild{})
}

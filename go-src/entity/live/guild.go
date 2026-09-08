package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveGuild db.TbName = "live_guilds"
)

const (
	LiveGuildName        db.TbCol = "name"
	LiveGuildLeaderId    db.TbCol = "leader_id"
	LiveGuildLeaderName  db.TbCol = "leader_name"
	LiveGuildDescription db.TbCol = "description"
	LiveGuildStatus      db.TbCol = "status"
	LiveGuildGuildType    db.TbCol = "guild_type"
	LiveGuildSharePercent db.TbCol = "share_percent"
	LiveGuildCreatorId    db.TbCol = "creator_id"
	LiveGuildCreatorName  db.TbCol = "creator_name"
)

// 工会上下架状态
const (
	LiveGuildStatusOffShelf uint8 = 0 // 下架(列表不可见)
	LiveGuildStatusOnShelf  uint8 = 1 // 上架
)

// 工会类型
const (
	LiveGuildTypeNormal       uint8 = 0 // 普通工会
	LiveGuildTypeCoinMerchant uint8 = 1 // 币商工会
)

// LiveGuild 直播工会(读写直连数据库,status=0 表示下架)
type LiveGuild struct {
	migrate.OneModel
	Name         string  `gorm:"size:64;comment:工会名称" json:"name"`
	LeaderId     uint64  `gorm:"default:0;comment:会长/负责人ID" json:"leaderId"`
	LeaderName   string  `gorm:"size:64;default:'';comment:会长名称" json:"leaderName"`
	Description  string  `gorm:"size:255;comment:工会简介" json:"description"`
	Status       uint8   `gorm:"default:1;comment:状态(0-下架,1-上架)" json:"status"`
	GuildType    uint8   `gorm:"default:0;comment:工会类型(0普通,1币商)" json:"guildType"`
	SharePercent float64 `gorm:"type:decimal(6,2);default:10;comment:分佣比例(%);币商工会可自定义" json:"sharePercent"`
	CreatorId    uint64  `gorm:"default:0;comment:创建者CMS用户ID" json:"creatorId"`
	CreatorName  string  `gorm:"size:64;default:'';comment:创建者CMS用户名" json:"creatorName"`
}

// NewLiveGuild 构造工会对象(不写库);分佣比例由调用方按全局配置/默认值填入
func NewLiveGuild(id uint64, name string, leaderId uint64, leaderName, description string) *LiveGuild {
	now := time.Now()
	return &LiveGuild{
		OneModel: migrate.OneModel{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:         name,
		LeaderId:     leaderId,
		LeaderName:   leaderName,
		Description:  description,
		Status:       LiveGuildStatusOnShelf,
		GuildType:    LiveGuildTypeNormal,
		SharePercent: 10,
	}
}

func InitLiveGuild() {
	initLiveGuild()
}

func initLiveGuild() {
	migrate.AutoMigrate(&LiveGuild{})
}

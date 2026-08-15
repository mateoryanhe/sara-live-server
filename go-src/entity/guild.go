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
)

// 工会状态(软删除)
const (
	LiveGuildStatusDeleted uint8 = 0 // 已删除(列表不可见)
	LiveGuildStatusNormal  uint8 = 1 // 正常
)

// LiveGuild 直播工会(读写直连数据库,status=0 表示软删除)
type LiveGuild struct {
	migrate.OneModel
	Name        string `gorm:"size:64;comment:工会名称" json:"name"`
	LeaderId    uint64 `gorm:"default:0;comment:会长/负责人ID" json:"leaderId"`
	LeaderName  string `gorm:"size:64;default:'';comment:会长名称" json:"leaderName"`
	Description string `gorm:"size:255;comment:工会简介" json:"description"`
	Status      uint8  `gorm:"default:1;comment:状态(0-已删除,1-正常)" json:"status"`
}

// NewLiveGuild 构造工会对象(不写库)
func NewLiveGuild(id uint64, name string, leaderId uint64, leaderName, description string) *LiveGuild {
	now := time.Now()
	return &LiveGuild{
		OneModel: migrate.OneModel{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        name,
		LeaderId:    leaderId,
		LeaderName:  leaderName,
		Description: description,
		Status:      LiveGuildStatusNormal,
	}
}

func InitLiveGuild() {
	initLiveGuild()
}

func initLiveGuild() {
	migrate.AutoMigrate(&LiveGuild{})
}

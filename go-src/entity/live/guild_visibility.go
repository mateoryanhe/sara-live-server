package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
)

const (
	TbLiveGuildVisibility db.TbName = "live_guild_visibilities"
)

const (
	LiveGuildVisibilityGuildId   db.TbCol = "guild_id"
	LiveGuildVisibilityCmsUserId db.TbCol = "cms_user_id"
)

// LiveGuildVisibility CMS 用户对工会的可见性(创建时自动写入创建人)
type LiveGuildVisibility struct {
	migrate.OneModel
	GuildId   uint64 `gorm:"uniqueIndex:idx_guild_visibility_user;index;default:0;comment:工会ID" json:"guildId"`
	CmsUserId uint64 `gorm:"uniqueIndex:idx_guild_visibility_user;index;default:0;comment:可见的CMS用户ID" json:"cmsUserId"`
}

// NewLiveGuildVisibility 构造可见性记录(不写库)
func NewLiveGuildVisibility(guildId, cmsUserId uint64) *LiveGuildVisibility {
	now := time.Now()
	return &LiveGuildVisibility{
		OneModel: migrate.OneModel{
			ID:        snowflake.GetId(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		GuildId:   guildId,
		CmsUserId: cmsUserId,
	}
}

func initLiveGuildVisibility() {
	migrate.AutoMigrate(&LiveGuildVisibility{})
}

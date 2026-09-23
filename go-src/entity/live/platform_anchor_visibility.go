package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
)

const (
	TbLivePlatformAnchorVisibility db.TbName = "live_platform_anchor_visibilities"
)

const (
	LivePlatformAnchorVisibilityAnchorId  db.TbCol = "anchor_id"
	LivePlatformAnchorVisibilityCmsUserId db.TbCol = "cms_user_id"
)

// LivePlatformAnchorVisibility CMS 用户对平台主播的可见性。
type LivePlatformAnchorVisibility struct {
	migrate.OneModel
	AnchorId  uint64 `gorm:"uniqueIndex:idx_platform_anchor_visibility_user;index;default:0;comment:平台主播ID" json:"anchorId"`
	CmsUserId uint64 `gorm:"uniqueIndex:idx_platform_anchor_visibility_user;index;default:0;comment:可见的CMS用户ID" json:"cmsUserId"`
}

// NewLivePlatformAnchorVisibility 构造平台主播可见性记录。
func NewLivePlatformAnchorVisibility(anchorId, cmsUserId uint64) *LivePlatformAnchorVisibility {
	now := time.Now()
	return &LivePlatformAnchorVisibility{
		OneModel: migrate.OneModel{
			ID:        snowflake.GetId(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		AnchorId:  anchorId,
		CmsUserId: cmsUserId,
	}
}

func initLivePlatformAnchorVisibility() {
	migrate.AutoMigrate(&LivePlatformAnchorVisibility{})
}

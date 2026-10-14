package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbPromotionChannelUrl db.TbName = "promotion_channel_urls"
)

const (
	PromotionChannelUrlChannelName db.TbCol = "channel_name"
	PromotionChannelUrlChannelCode db.TbCol = "channel_code"
	PromotionChannelUrlURL         db.TbCol = "url"
)

// PromotionChannelUrl CMS 推广渠道 URL 配置
type PromotionChannelUrl struct {
	migrate.OneModel
	ChannelName string `gorm:"size:128;default:'';comment:渠道名称" json:"channelName"`
	ChannelCode int    `gorm:"uniqueIndex;comment:渠道编码" json:"channelCode"`
	URL         string `gorm:"size:512;default:'';comment:推广链接URL" json:"url"`
}

func initPromotionChannelUrl() {
	migrate.AutoMigrate(&PromotionChannelUrl{})
}

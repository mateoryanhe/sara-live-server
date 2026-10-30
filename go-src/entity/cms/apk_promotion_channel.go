package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbApkPromotionChannel db.TbName = "apk_promotion_channels"
)

const (
	ApkPromotionChannelChannelName     db.TbCol = "channel_name"
	ApkPromotionChannelChannelCode     db.TbCol = "channel_code"
	ApkPromotionChannelApkFileName   db.TbCol = "apk_file_name"
	ApkPromotionChannelApkDownloadURL db.TbCol = "apk_download_url"
)

// ApkPromotionChannel CMS APK 推广渠道(字符串渠道码 + 专用 APK)
type ApkPromotionChannel struct {
	migrate.OneModel
	ChannelName     string `gorm:"size:128;default:'';comment:渠道名称" json:"channelName"`
	ChannelCode     string `gorm:"size:64;uniqueIndex;comment:渠道编码(字符串,唯一)" json:"channelCode"`
	ApkFileName     string `gorm:"size:256;default:'';comment:APK云桶对象名(含版本时间戳)" json:"apkFileName"`
	ApkDownloadURL  string `gorm:"size:512;default:'';comment:APK公开下载URL" json:"apkDownloadUrl"`
}

func initApkPromotionChannel() {
	migrate.AutoMigrate(&ApkPromotionChannel{})
}

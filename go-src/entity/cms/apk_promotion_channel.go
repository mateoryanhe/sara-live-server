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
	ApkPromotionChannelH5AppChannelCode db.TbCol = "h5_app_channel_code"
	ApkPromotionChannelApkFileName   db.TbCol = "apk_file_name"
	ApkPromotionChannelApkDownloadURL db.TbCol = "apk_download_url"
)

// ApkPromotionChannel CMS APK 推广渠道(字符串渠道码 + 专用 APK)
type ApkPromotionChannel struct {
	migrate.OneModel
	ChannelName      string `gorm:"size:128;default:'';comment:渠道名称" json:"channelName"`
	H5AppChannelCode string `gorm:"size:64;uniqueIndex:idx_apk_h5_pair,priority:1;default:'';comment:H5-App渠道码(须存在于推广URL列表)" json:"h5AppChannelCode"`
	ChannelCode      string `gorm:"size:64;uniqueIndex:idx_apk_h5_pair,priority:2;comment:APK渠道编码" json:"channelCode"`
	ApkFileName      string `gorm:"size:256;default:'';comment:APK云桶对象名(含版本时间戳)" json:"apkFileName"`
	ApkDownloadURL   string `gorm:"size:512;default:'';comment:APK公开下载URL" json:"apkDownloadUrl"`
}

func initApkPromotionChannel() {
	migrate.AutoMigrate(&ApkPromotionChannel{})
}

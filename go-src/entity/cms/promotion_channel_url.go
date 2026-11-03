package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbPromotionChannelUrl db.TbName = "promotion_channel_urls"
)

const (
	PromotionChannelUrlChannelName     db.TbCol = "channel_name"
	PromotionChannelUrlChannelCode     db.TbCol = "channel_code"
	PromotionChannelUrlURL             db.TbCol = "url"
	PromotionChannelUrlH5AppChannelCode db.TbCol = "h5_app_channel_code"
	PromotionChannelUrlApkFileName     db.TbCol = "apk_file_name"
	PromotionChannelUrlApkDownloadURL  db.TbCol = "apk_download_url"
)

// PromotionChannelUrl CMS 推广渠道 URL 配置(H5 分享 + lookupApkDownload 渠道包)
type PromotionChannelUrl struct {
	migrate.OneModel
	ChannelName      string `gorm:"size:128;default:'';comment:渠道名称" json:"channelName"`
	ChannelCode      int    `gorm:"comment:渠道编码(H5链接ch整数,CMS校验唯一)" json:"channelCode"`
	H5AppChannelCode string `gorm:"size:64;default:'';comment:H5-App渠道码(与APK渠道码配对,非空时CMS校验唯一)" json:"h5AppChannelCode"`
	URL               string `gorm:"size:512;default:'';comment:推广链接URL" json:"url"`
	ApkFileName       string `gorm:"size:512;default:'';comment:云桶对象名" json:"apkFileName"`
	ApkDownloadURL    string `gorm:"size:1024;default:'';comment:渠道APK公开URL" json:"apkDownloadUrl"`
}

func initPromotionChannelUrl() {
	migrate.AutoMigrate(&PromotionChannelUrl{})
}

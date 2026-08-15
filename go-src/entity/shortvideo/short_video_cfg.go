package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbShortVideoCfg db.TbName = "short_video_cfgs"
)

const (
	ShortVideoCfgEntryDisabled uint8 = 0 // 入口关闭
	ShortVideoCfgEntryEnabled  uint8 = 1 // 入口开启
)

// ShortVideoCfg 短视频全局配置(CMS管理,通常仅一条记录)
type ShortVideoCfg struct {
	migrate.OneModel
	MaxFileSize                uint64 `gorm:"default:0;comment:最大文件大小(字节)" json:"maxFileSize"`
	MaxCoverFileSize           uint32 `gorm:"default:5;comment:封面图片最大大小(M)" json:"maxCoverFileSize"`
	MaxDuration                uint32 `gorm:"default:0;comment:最大时长(秒)" json:"maxDuration"`
	FreeWatchSeconds           uint32 `gorm:"default:7;comment:免费观看时长(秒)" json:"freeWatchSeconds"`
	EntryEnabled               uint8  `gorm:"default:0;comment:入口开关(0关闭,1开启)" json:"entryEnabled"`
	AnchorDailyUploadLimit     uint32 `gorm:"default:100;comment:主播每日最大上传数量" json:"anchorDailyUploadLimit"`
	NormalUserDailyUploadLimit uint32 `gorm:"default:1;comment:普通用户(类型0)每日最大上传数量" json:"normalUserDailyUploadLimit"`
}

func initShortVideoCfg() {
	migrate.AutoMigrate(&ShortVideoCfg{})
}

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
	MaxFileSize      uint64 `gorm:"default:0;comment:最大文件大小(字节)" json:"maxFileSize"`
	MaxDuration      uint32 `gorm:"default:0;comment:最大时长(秒)" json:"maxDuration"`
	FreeWatchSeconds uint32 `gorm:"default:7;comment:免费观看时长(秒)" json:"freeWatchSeconds"`
	EntryEnabled     uint8  `gorm:"default:0;comment:入口开关(0关闭,1开启)" json:"entryEnabled"`
}

func initShortVideoCfg() {
	migrate.AutoMigrate(&ShortVideoCfg{})
}

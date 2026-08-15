package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbPreloadCfg db.TbName = "preload_cfgs"
)

// PreloadCfg 启动预热配置(CMS 管理,通常仅一条)
type PreloadCfg struct {
	migrate.OneModel
	RecentLoginLimit int `gorm:"default:100;comment:最近登录用户预热数量" json:"recentLoginLimit"`
}

func initPreloadCfg() {
	migrate.AutoMigrate(&PreloadCfg{})
}

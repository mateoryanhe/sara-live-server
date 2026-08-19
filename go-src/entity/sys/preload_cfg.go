package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbPreloadCfg db.TbName = "preload_cfgs"
)

// PreloadCfg 启动预热与服务器运行配置(CMS 管理,通常仅一条)
type PreloadCfg struct {
	migrate.OneModel
	RecentLoginLimit int    `gorm:"default:100;comment:最近登录用户预热数量" json:"recentLoginLimit"`
	HotRestartAuth   string `gorm:"size:128;comment:热重启接口密钥" json:"hotRestartAuth"`
	MemoryLimitM     int    `gorm:"default:300;comment:Go堆内存软上限MB" json:"memoryLimitM"`
	IpGeoDbPath      string `gorm:"size:512;comment:GeoLite2-Country.mmdb路径" json:"ipGeoDbPath"`
}

func initPreloadCfg() {
	migrate.AutoMigrate(&PreloadCfg{})
}

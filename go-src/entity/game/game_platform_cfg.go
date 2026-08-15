package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbGamePlatformCfg db.TbName = "game_platform_cfgs"
)

// GamePlatformDefaultVendorUrl 厂家 API 默认地址
const GamePlatformDefaultVendorUrl = "https://gapi.win12.best"

// GamePlatformCfg 第三方游戏平台接入配置(CMS 管理,通常仅一条)
type GamePlatformCfg struct {
	migrate.OneModel
	VendorUrl string `gorm:"size:255;default:'';comment:厂家API根地址" json:"vendorUrl"`
	Token     string `gorm:"size:512;default:'';comment:接入Token(x-token)" json:"token"`
	SecretKey string `gorm:"size:255;default:'';comment:SecretKey" json:"secretKey"`
	IconUrl   string `gorm:"size:512;default:'';comment:游戏封面CDN根地址" json:"iconUrl"`
}

func initGamePlatformCfg() {
	migrate.AutoMigrate(&GamePlatformCfg{})
}

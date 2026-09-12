package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbCountryFlagCfg db.TbName = "country_flag_cfgs"
)

// CountryFlagCfg 国旗资源部署配置(通常一条;Version 为当前生效目录名)
type CountryFlagCfg struct {
	migrate.OneModel
	Version string `gorm:"size:32;default:'';comment:当前国旗资源版本目录名" json:"version"`
}

func initCountryFlagCfg() {
	migrate.AutoMigrate(&CountryFlagCfg{})
}

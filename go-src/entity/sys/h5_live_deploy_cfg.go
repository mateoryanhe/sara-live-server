package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbH5LiveDeployCfg db.TbName = "h5_live_deploy_cfgs"
)

// H5LiveDeployCfg H5 直播部署配置(CMS 管理,通常仅一条)
type H5LiveDeployCfg struct {
	migrate.OneModel
	DeploySecret string `gorm:"size:128;default:'';comment:H5部署密钥" json:"deploySecret"`
}

func initH5LiveDeployCfg() {
	migrate.AutoMigrate(&H5LiveDeployCfg{})
}

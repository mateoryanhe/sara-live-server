package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbShareActivityDeployCfg db.TbName = "share_activity_deploy_cfgs"
)

// ShareActivityDeployCfg 分享活动页部署配置(CMS 管理,通常仅一条;域名映射见 domain_site_mappings)
type ShareActivityDeployCfg struct {
	migrate.OneModel
}

func initShareActivityDeployCfg() {
	migrate.AutoMigrate(&ShareActivityDeployCfg{})
}

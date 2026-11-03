package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbChannelDownloadLandingDeployCfg db.TbName = "channel_download_landing_deploy_cfgs"
)

// ChannelDownloadLandingDeployCfg 渠道下载落地页部署配置(CMS 管理,通常仅一条;域名映射见 domain_site_mappings)
type ChannelDownloadLandingDeployCfg struct {
	migrate.OneModel
}

func initChannelDownloadLandingDeployCfg() {
	migrate.AutoMigrate(&ChannelDownloadLandingDeployCfg{})
}

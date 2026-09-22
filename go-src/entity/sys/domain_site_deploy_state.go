package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const TbDomainSiteDeployState db.TbName = "domain_site_deploy_states"

// DomainSiteDeployState 保存每个静态站点最后一次成功完成 ZIP 部署的时间。
type DomainSiteDeployState struct {
	migrate.OneModel
	SiteKey      string    `gorm:"size:64;uniqueIndex;default:'';comment:站点业务键" json:"siteKey"`
	LastUploadAt time.Time `gorm:"comment:最后一次成功上传并解压时间" json:"lastUploadAt"`
}

func initDomainSiteDeployState() {
	migrate.AutoMigrate(&DomainSiteDeployState{})
}

package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbDomainSiteMapping db.TbName = "domain_site_mappings"
)

// DomainSiteMapping 保存可由 CMS 热更新的单域名、URL 前缀与物理目录映射。
// SiteKey 是业务稳定键；Prefix 用于覆盖 config.yaml 中同前缀的静态站点。
type DomainSiteMapping struct {
	migrate.OneModel
	SiteKey string `gorm:"size:64;uniqueIndex;default:'';comment:站点业务键" json:"siteKey"`
	Prefix  string `gorm:"size:128;uniqueIndex;default:'';comment:静态URL前缀" json:"prefix"`
	Domain  string `gorm:"size:512;default:'';comment:注册域名,只保存一个Host" json:"domain"`
	Root    string `gorm:"size:1024;default:'';comment:静态文件物理目录" json:"root"`
}

func initDomainSiteMapping() {
	migrate.AutoMigrate(&DomainSiteMapping{})
}

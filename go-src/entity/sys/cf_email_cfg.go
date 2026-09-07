package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbCfEmailCfg db.TbName = "cf_email_cfgs"
)

// CfEmailCfg Cloudflare Email Sending 发信配置(CMS 管理,通常仅一条)
// 域名如 mail.saralive.net 需在 CF Email Service 中 Onboard,并由 CF DNS 托管解析
type CfEmailCfg struct {
	migrate.OneModel
	Enabled   bool   `gorm:"default:0;comment:是否启用" json:"enabled"`
	AccountId string `gorm:"size:64;default:'';comment:Cloudflare Account ID" json:"accountId"`
	ApiToken  string `gorm:"size:256;default:'';comment:API Token(需Email Sending Edit权限)" json:"apiToken"`
	FromEmail string `gorm:"size:256;default:'';comment:发件地址(如noreply@mail.saralive.net)" json:"fromEmail"`
}

func (CfEmailCfg) TableName() string {
	return string(TbCfEmailCfg)
}

func initCfEmailCfg() {
	migrate.AutoMigrate(&CfEmailCfg{})
}

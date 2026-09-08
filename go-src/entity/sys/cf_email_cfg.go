package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbCfEmailCfg db.TbName = "cf_email_cfgs"
)

// CfEmailCfg 邮件发信配置(CMS 管理,通常仅一条);现走 AWS SES
// 表名沿用 cf_email_cfgs;发件地址须在 SES 完成身份/域验证
type CfEmailCfg struct {
	migrate.OneModel
	Enabled         bool   `gorm:"default:0;comment:是否启用" json:"enabled"`
	Region          string `gorm:"size:64;default:'';comment:AWS Region(如us-west-1)" json:"region"`
	AccessKeyId     string `gorm:"size:128;default:'';comment:AWS Access Key ID" json:"accessKeyId"`
	SecretAccessKey string `gorm:"size:256;default:'';comment:AWS Secret Access Key" json:"secretAccessKey"`
	FromEmail       string `gorm:"size:256;default:'';comment:发件地址(须已在SES验证)" json:"fromEmail"`
}

func (CfEmailCfg) TableName() string {
	return string(TbCfEmailCfg)
}

func initCfEmailCfg() {
	migrate.AutoMigrate(&CfEmailCfg{})
}

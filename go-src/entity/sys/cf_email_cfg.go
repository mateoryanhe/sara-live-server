package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbCfEmailCfg db.TbName = "cf_email_cfgs"
)

// CfEmailCfg 邮件发信配置(CMS 管理,通常仅一条);走 Amazon SES SMTP(STARTTLS,通常 587).
// 表名沿用 cf_email_cfgs;发件地址须在 SES 完成身份/域验证.
type CfEmailCfg struct {
	migrate.OneModel
	Enabled      bool   `gorm:"default:0;comment:是否启用" json:"enabled"`
	SmtpHost     string `gorm:"size:256;default:'';comment:SMTP服务器地址" json:"smtpHost"`
	SmtpPort     int    `gorm:"default:587;comment:SMTP端口" json:"smtpPort"`
	SmtpUsername string `gorm:"size:128;default:'';comment:SMTP用户名" json:"smtpUsername"`
	SmtpPassword string `gorm:"size:256;default:'';comment:SMTP密码" json:"smtpPassword"`
	FromEmail    string `gorm:"size:256;default:'';comment:发件地址(须已在SES验证)" json:"fromEmail"`
	// 以下字段已废弃(原 SES API),保留列兼容旧数据
	Region          string `gorm:"size:64;default:'';comment:已废弃" json:"region,omitempty"`
	AccessKeyId     string `gorm:"size:128;default:'';comment:已废弃" json:"accessKeyId,omitempty"`
	SecretAccessKey string `gorm:"size:256;default:'';comment:已废弃" json:"secretAccessKey,omitempty"`
}

func (CfEmailCfg) TableName() string {
	return string(TbCfEmailCfg)
}

func initCfEmailCfg() {
	migrate.AutoMigrate(&CfEmailCfg{})
}

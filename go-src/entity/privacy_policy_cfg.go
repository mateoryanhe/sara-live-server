package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbPrivacyPolicyCfg db.TbName = "privacy_policy_cfgs"
)

// PrivacyPolicyCfg 隐私政策配置(CMS 管理,通常仅一条)
type PrivacyPolicyCfg struct {
	migrate.OneModel
	PrivacyPolicyUrl  string `gorm:"size:512;default:'';comment:隐私政策页面URL" json:"privacyPolicyUrl"`
	TermsOfServiceUrl string `gorm:"size:512;default:'';comment:用户服务协议页面URL" json:"termsOfServiceUrl"`
}

func initPrivacyPolicyCfg() {
	migrate.AutoMigrate(&PrivacyPolicyCfg{})
}

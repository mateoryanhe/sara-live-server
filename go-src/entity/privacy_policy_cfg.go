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
	CreatorTermsUrl   string `gorm:"size:512;default:'';comment:短视频创作者上传合规条款URL" json:"creatorTermsUrl"`
	RoomOwnerTermsUrl string `gorm:"size:512;default:'';comment:房间房主责任条款URL" json:"roomOwnerTermsUrl"`
	VipDescUrl        string `gorm:"size:512;default:'';comment:VIP描述文档URL" json:"vipDescUrl"`
	AboutSiteUrl      string `gorm:"size:512;default:'';comment:About页面URL" json:"aboutSiteUrl"`
	SafetyCenterUrl   string `gorm:"size:512;default:'';comment:安全中心页面URL" json:"safetyCenterUrl"`
}

func initPrivacyPolicyCfg() {
	migrate.AutoMigrate(&PrivacyPolicyCfg{})
}

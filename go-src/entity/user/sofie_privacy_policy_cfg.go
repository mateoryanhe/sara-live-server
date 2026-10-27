package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbSofiePrivacyPolicyCfg db.TbName = "sofie_privacy_policy_cfgs"
)

// SofiePrivacyPolicyCfg Sofie 隐私/条款配置(CMS 管理,通常仅一条;与 Sara privacy_policy_cfgs 独立)
// PackageName 非空时,仅当请求头包名匹配才下发本配置;否则回落 Sara 隐私政策.
type SofiePrivacyPolicyCfg struct {
	migrate.OneModel
	PackageName       string `gorm:"size:128;default:'';index;comment:适用App包名(空则不下发本配置)" json:"packageName"`
	ApiBase           string `gorm:"size:512;default:'';comment:页面URL前缀(如 https://sofielive.net)" json:"apiBase"`
	PrivacyPolicyUrl  string `gorm:"size:512;default:'';comment:隐私政策页面URL" json:"privacyPolicyUrl"`
	TermsOfServiceUrl string `gorm:"size:512;default:'';comment:用户服务协议页面URL" json:"termsOfServiceUrl"`
	CreatorTermsUrl   string `gorm:"size:512;default:'';comment:短视频创作者上传合规条款URL" json:"creatorTermsUrl"`
	RoomOwnerTermsUrl string `gorm:"size:512;default:'';comment:房间房主责任条款URL" json:"roomOwnerTermsUrl"`
	VipDescUrl        string `gorm:"size:512;default:'';comment:VIP描述文档URL" json:"vipDescUrl"`
	AboutSiteUrl      string `gorm:"size:512;default:'';comment:About页面URL" json:"aboutSiteUrl"`
	SafetyCenterUrl   string `gorm:"size:512;default:'';comment:安全中心页面URL" json:"safetyCenterUrl"`
}

func initSofiePrivacyPolicyCfg() {
	migrate.AutoMigrate(&SofiePrivacyPolicyCfg{})
}

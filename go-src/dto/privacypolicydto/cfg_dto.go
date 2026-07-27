package privacypolicydto

import "github.com/gogf/gf/v2/frame/g"

type GetPrivacyPolicyCfgReq struct {
	g.Meta `path:"/getPrivacyPolicyCfg" method:"post" summary:"查询隐私政策配置" tags:"隐私政策配置"`
}

type PrivacyPolicyCfgItem struct {
	ID                string `json:"id"`
	PrivacyPolicyUrl  string `json:"privacyPolicyUrl"`
	TermsOfServiceUrl string `json:"termsOfServiceUrl"`
	CreatorTermsUrl   string `json:"creatorTermsUrl"`
	RoomOwnerTermsUrl string `json:"roomOwnerTermsUrl"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

type GetPrivacyPolicyCfgRes struct {
	Cfg *PrivacyPolicyCfgItem `json:"cfg"`
}

type SavePrivacyPolicyCfgReq struct {
	g.Meta            `path:"/savePrivacyPolicyCfg" method:"post" summary:"保存隐私政策配置" tags:"隐私政策配置"`
	ID                uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	PrivacyPolicyUrl  string `json:"privacyPolicyUrl" v:"max-length:512#隐私政策URL长度不能超过512" dc:"隐私政策页面URL"`
	TermsOfServiceUrl string `json:"termsOfServiceUrl" v:"max-length:512#用户服务协议URL长度不能超过512" dc:"用户服务协议页面URL"`
	CreatorTermsUrl   string `json:"creatorTermsUrl" v:"max-length:512#创作者条款URL长度不能超过512" dc:"短视频创作者上传合规条款URL"`
	RoomOwnerTermsUrl string `json:"roomOwnerTermsUrl" v:"max-length:512#房主责任条款URL长度不能超过512" dc:"房间房主责任条款URL"`
}

type SavePrivacyPolicyCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

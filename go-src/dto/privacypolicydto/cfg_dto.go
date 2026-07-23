package privacypolicydto

import "github.com/gogf/gf/v2/frame/g"

type GetPrivacyPolicyCfgReq struct {
	g.Meta `path:"/getPrivacyPolicyCfg" method:"post" summary:"查询隐私政策配置" tags:"隐私政策配置"`
}

type PrivacyPolicyCfgItem struct {
	ID               string `json:"id"`
	PrivacyPolicyUrl string `json:"privacyPolicyUrl"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type GetPrivacyPolicyCfgRes struct {
	Cfg *PrivacyPolicyCfgItem `json:"cfg"`
}

type SavePrivacyPolicyCfgReq struct {
	g.Meta           `path:"/savePrivacyPolicyCfg" method:"post" summary:"保存隐私政策配置" tags:"隐私政策配置"`
	ID               uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	PrivacyPolicyUrl string `json:"privacyPolicyUrl" v:"max-length:512#隐私政策URL长度不能超过512" dc:"隐私政策页面URL"`
}

type SavePrivacyPolicyCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

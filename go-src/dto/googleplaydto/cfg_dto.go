package googleplaydto

import "github.com/gogf/gf/v2/frame/g"

type GetGooglePlayCfgReq struct {
	g.Meta `path:"/getGooglePlayCfg" method:"post" summary:"查询Google Play配置" tags:"Google Play配置"`
}

type GooglePlayCfgItem struct {
	ID                 string `json:"id"`
	Enabled            bool   `json:"enabled"`
	PackageName        string `json:"packageName"`
	ServiceAccountJson string `json:"serviceAccountJson"`
	RtdnAudience       string `json:"rtdnAudience"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type GetGooglePlayCfgRes struct {
	Cfg *GooglePlayCfgItem `json:"cfg"`
}

type SaveGooglePlayCfgReq struct {
	g.Meta             `path:"/saveGooglePlayCfg" method:"post" summary:"保存Google Play配置" tags:"Google Play配置"`
	ID                 uint64 `json:"id" dc:"配置ID,新建传0"`
	Enabled            bool   `json:"enabled" dc:"是否启用RTDN充值到账"`
	PackageName        string `json:"packageName" v:"max-length:128#包名最长128字符" dc:"Android包名"`
	ServiceAccountJson string `json:"serviceAccountJson" dc:"Google服务账号JSON全文"`
	RtdnAudience       string `json:"rtdnAudience" v:"max-length:512#RTDN Audience最长512字符" dc:"Pub/Sub Push JWT aud"`
}

type SaveGooglePlayCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

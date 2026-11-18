package appsflyercfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetAppsFlyerCfgReq struct {
	g.Meta `path:"/getAppsFlyerCfg" method:"post" summary:"查询AppsFlyer配置" tags:"AppsFlyer配置"`
}

type AppsFlyerCfgItem struct {
	ID            string `json:"id"`
	AndroidDevKey string `json:"androidDevKey"`
	IosDevKey     string `json:"iosDevKey"`
	IosAppId      string `json:"iosAppId"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type GetAppsFlyerCfgRes struct {
	Cfg *AppsFlyerCfgItem `json:"cfg"`
}

type SaveAppsFlyerCfgReq struct {
	g.Meta        `path:"/saveAppsFlyerCfg" method:"post" summary:"保存AppsFlyer配置" tags:"AppsFlyer配置"`
	ID            uint64 `json:"id" dc:"配置ID,新建传0"`
	AndroidDevKey string `json:"androidDevKey" v:"max-length:128#Android Dev Key最长128字符" dc:"Android Dev Key"`
	IosDevKey     string `json:"iosDevKey" v:"max-length:128#iOS Dev Key最长128字符" dc:"iOS Dev Key"`
	IosAppId      string `json:"iosAppId" v:"max-length:128#iOS App ID最长128字符" dc:"iOS App ID"`
}

type SaveAppsFlyerCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type GetAppsFlyerCfgForAppReq struct {
	g.Meta `path:"/getCfgForApp" method:"post" summary:"App查询AppsFlyer配置(读缓存,按X-Platform)" tags:"AppsFlyer配置"`
}

type AppsFlyerCfgForApp struct {
	AppsFlyerDevKey string `json:"appsFlyerDevKey" dc:"AppsFlyer Dev Key"`
	AppsFlyerAppId  string `json:"appsFlyerAppId" dc:"iOS App ID,Android为空"`
}

type GetAppsFlyerCfgForAppRes struct {
	Cfg *AppsFlyerCfgForApp `json:"cfg" dc:"未配置或平台未知时为null"`
}
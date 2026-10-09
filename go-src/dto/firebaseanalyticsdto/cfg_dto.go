package firebaseanalyticsdto

import "github.com/gogf/gf/v2/frame/g"

type GetFirebaseAnalyticsCfgReq struct {
	g.Meta `path:"/getFirebaseAnalyticsCfg" method:"post" summary:"查询Firebase Analytics埋点配置" tags:"Firebase Analytics配置"`
}

type FirebaseAnalyticsCfgItem struct {
	ID                 string `json:"id"`
	Enabled            uint8  `json:"enabled" dc:"0关闭1启用"`
	ProjectId          string `json:"projectId"`
	ClientConfigJson   string `json:"clientConfigJson"`
	ServiceAccountJson   string `json:"serviceAccountJson"`
	MeasurementApiSecret string `json:"measurementApiSecret"`
	CreatedAt            string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type GetFirebaseAnalyticsCfgRes struct {
	Cfg *FirebaseAnalyticsCfgItem `json:"cfg"`
}

type SaveFirebaseAnalyticsCfgReq struct {
	g.Meta           `path:"/saveFirebaseAnalyticsCfg" method:"post" summary:"保存Firebase Analytics埋点配置" tags:"Firebase Analytics配置"`
	ID               uint64 `json:"id" dc:"配置ID,新建传0"`
	Enabled          uint8  `json:"enabled" v:"in:0,1#启用状态仅支持0或1" dc:"0关闭1启用"`
	ProjectId          string `json:"projectId" v:"max-length:128#Project ID最长128字符" dc:"Firebase Project ID"`
	ClientConfigJson   string `json:"clientConfigJson" dc:"Firebase客户端配置JSON(Analytics SDK初始化)"`
	ServiceAccountJson   string `json:"serviceAccountJson" dc:"Firebase服务账号JSON(Admin SDK/服务端上报)"`
	MeasurementApiSecret string `json:"measurementApiSecret" v:"max-length:128#API Secret最长128字符" dc:"GA4 Measurement Protocol API Secret"`
}

type SaveFirebaseAnalyticsCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

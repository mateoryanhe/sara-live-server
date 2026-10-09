package firebaseanalyticsdto

import "github.com/gogf/gf/v2/frame/g"

type GetFirebaseAnalyticsClientCfgForAppReq struct {
	g.Meta `path:"/getClientCfgForApp" method:"post" summary:"App免登录查询Firebase Analytics客户端配置" tags:"Firebase Analytics"`
}

// FirebaseAnalyticsClientCfgForApp 可公开给 App/H5 的 Analytics SDK 初始化字段.
type FirebaseAnalyticsClientCfgForApp struct {
	ApiKey            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain,omitempty"`
	ProjectId         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket,omitempty"`
	MessagingSenderId string `json:"messagingSenderId,omitempty"`
	AppId             string `json:"appId"`
	MeasurementId     string `json:"measurementId,omitempty"`
	DatabaseUrl       string `json:"databaseURL,omitempty"`
}

type GetFirebaseAnalyticsClientCfgForAppRes struct {
	Enabled uint8                               `json:"enabled" dc:"0关闭1启用"`
	Cfg     *FirebaseAnalyticsClientCfgForApp `json:"cfg" dc:"未启用或未配置时为null"`
}

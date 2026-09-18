package firebasedto

import "github.com/gogf/gf/v2/frame/g"

type GetFirebaseCfgReq struct {
	g.Meta `path:"/getFirebaseCfg" method:"post" summary:"查询Firebase登录配置" tags:"Firebase登录配置"`
}

type FirebaseCfgItem struct {
	ID                 string `json:"id"`
	ProjectId          string `json:"projectId"`
	ClientConfigJson   string `json:"clientConfigJson"`
	ServiceAccountJson string `json:"serviceAccountJson"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type GetFirebaseCfgRes struct {
	Cfg *FirebaseCfgItem `json:"cfg"`
}

type SaveFirebaseCfgReq struct {
	g.Meta             `path:"/saveFirebaseCfg" method:"post" summary:"保存Firebase登录配置" tags:"Firebase登录配置"`
	ID                 uint64 `json:"id" dc:"配置ID,新建传0"`
	ProjectId          string `json:"projectId" v:"required|max-length:128#Project ID不能为空|Project ID最长128字符" dc:"Firebase Project ID"`
	ClientConfigJson   string `json:"clientConfigJson" v:"required#客户端配置JSON不能为空" dc:"Firebase客户端配置JSON全文"`
	ServiceAccountJson string `json:"serviceAccountJson" v:"required#服务账号JSON不能为空" dc:"Firebase服务账号JSON全文"`
}

type SaveFirebaseCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type GetFirebaseClientCfgForAppReq struct {
	g.Meta `path:"/getClientCfgForApp" method:"post" summary:"App免登录查询Firebase客户端配置" tags:"Firebase登录配置"`
}

// FirebaseClientCfgForApp 仅包含可公开给客户端的 Firebase 初始化字段。
type FirebaseClientCfgForApp struct {
	ApiKey            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain,omitempty"`
	ProjectId         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket,omitempty"`
	MessagingSenderId string `json:"messagingSenderId,omitempty"`
	AppId             string `json:"appId"`
	MeasurementId     string `json:"measurementId,omitempty"`
	DatabaseUrl       string `json:"databaseURL,omitempty"`
}

type GetFirebaseClientCfgForAppRes struct {
	Cfg *FirebaseClientCfgForApp `json:"cfg" dc:"未配置时为null"`
}

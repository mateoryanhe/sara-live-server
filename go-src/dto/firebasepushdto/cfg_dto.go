package firebasepushdto

import "github.com/gogf/gf/v2/frame/g"

type GetFirebasePushCfgReq struct {
	g.Meta `path:"/getFirebasePushCfg" method:"post" summary:"查询Firebase推送配置" tags:"Firebase推送"`
}

type FirebasePushCfgItem struct {
	ID               string `json:"id"`
	Enabled          uint8  `json:"enabled" dc:"0关1开"`
	OfflineCallTitle string `json:"offlineCallTitle" dc:"1v1主播离线推送标题"`
	OfflineCallBody  string `json:"offlineCallBody" dc:"1v1主播离线推送文案,支持{nickname}"`
	OfflineCallIcon  string `json:"offlineCallIcon" dc:"1v1主播离线推送图标URL,空=App默认"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type GetFirebasePushCfgRes struct {
	Cfg *FirebasePushCfgItem `json:"cfg"`
}

type SaveFirebasePushCfgReq struct {
	g.Meta           `path:"/saveFirebasePushCfg" method:"post" summary:"保存Firebase推送配置" tags:"Firebase推送"`
	ID               uint64 `json:"id" dc:"配置ID,新建传0"`
	Enabled          uint8  `json:"enabled" dc:"0关1开"`
	OfflineCallTitle string `json:"offlineCallTitle" v:"max-length:128#标题最长128字符" dc:"1v1主播离线推送标题"`
	OfflineCallBody  string `json:"offlineCallBody" v:"max-length:512#文案最长512字符" dc:"1v1主播离线推送文案"`
	OfflineCallIcon  string `json:"offlineCallIcon" v:"max-length:512#图标URL最长512字符" dc:"1v1主播离线推送图标URL"`
}

type SaveFirebasePushCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

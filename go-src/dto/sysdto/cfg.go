package sysdto

import "github.com/gogf/gf/v2/frame/g"

type SysCfgReq struct {
	g.Meta `path:"/cfg" method:"get" summary:"获取系统配置" tags:"系统"`
}

type SysCfgResp struct {
	SysTime                     int64   `json:"sysTime"`
	PaidDanmakuPrice            float64 `json:"paidDanmakuPrice" dc:"直播间付费弹幕价格(钻石)"`
	PrivateRoomFreeWatchSeconds uint32  `json:"privateRoomFreeWatchSeconds" dc:"私密直播间免费观看时长(秒)"`
	PrivacyPolicyUrl            string  `json:"privacyPolicyUrl" dc:"隐私政策页面URL"`
	TermsOfServiceUrl           string  `json:"termsOfServiceUrl" dc:"用户服务协议页面URL"`
	CreatorTermsUrl             string  `json:"creatorTermsUrl" dc:"短视频创作者上传合规条款URL"`
	RoomOwnerTermsUrl           string  `json:"roomOwnerTermsUrl" dc:"房间房主责任条款URL"`
	VipDescUrl                  string  `json:"vipDescUrl" dc:"VIP描述文档URL"`
	AppImageMaxSize             uint64  `json:"appImageMaxSize" dc:"App端图片上传大小上限(字节)"`
	GoldToDiamondRate           int     `json:"goldToDiamondRate" dc:"金币兑换钻石比例(1金币=N钻石)"`
	ExchangeFeePercent          float64 `json:"exchangeFeePercent" dc:"App手动兑换手续费(%)，大于0时收取"`
	AboutSiteUrl                string  `json:"aboutSiteUrl" dc:"About页面URL"`
	SafetyCenterUrl             string  `json:"safetyCenterUrl" dc:"安全中心页面URL"`
}

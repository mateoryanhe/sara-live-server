package activitydto

import "github.com/gogf/gf/v2/frame/g"

type GetFirstRechargeActivityCfgReq struct {
	g.Meta `path:"/getFirstRechargeActivityCfg" method:"post" summary:"查询首充活动配置" tags:"首充活动配置"`
}

type FirstRechargePrivilegeItem struct {
	Icon     string `json:"icon"`
	IconName string `json:"iconName"`
	DescEn   string `json:"descEn"`
	DescEs   string `json:"descEs"`
	DescPt   string `json:"descPt"`
	DescHi   string `json:"descHi"`
	DescId   string `json:"descId"`
}

type FirstRechargeActivityCfgItem struct {
	ID                string                        `json:"id"`
	Enabled           bool                          `json:"enabled"`
	Icon              string                        `json:"icon"`
	IconName          string                        `json:"iconName"`
	TitleEn           string                        `json:"titleEn"`
	TitleEs           string                        `json:"titleEs"`
	TitlePt           string                        `json:"titlePt"`
	TitleHi           string                        `json:"titleHi"`
	TitleId           string                        `json:"titleId"`
	RechargeBtnTextEn string                        `json:"rechargeBtnTextEn"`
	RechargeBtnTextEs string                        `json:"rechargeBtnTextEs"`
	RechargeBtnTextPt string                        `json:"rechargeBtnTextPt"`
	RechargeBtnTextHi string                        `json:"rechargeBtnTextHi"`
	RechargeBtnTextId string                        `json:"rechargeBtnTextId"`
	Privileges        []*FirstRechargePrivilegeItem `json:"privileges"`
	CreatedAt         string                        `json:"createdAt"`
	UpdatedAt         string                        `json:"updatedAt"`
}

type GetFirstRechargeActivityCfgRes struct {
	Cfg *FirstRechargeActivityCfgItem `json:"cfg"`
}

type SaveFirstRechargeActivityCfgReq struct {
	g.Meta            `path:"/saveFirstRechargeActivityCfg" method:"post" summary:"保存首充活动配置" tags:"首充活动配置"`
	ID                uint64                        `json:"id" dc:"配置ID,首次保存可为0"`
	Enabled           bool                          `json:"enabled" dc:"活动开关"`
	Icon              string                        `json:"icon" dc:"小图标资源名"`
	TitleEn           string                        `json:"titleEn"`
	TitleEs           string                        `json:"titleEs"`
	TitlePt           string                        `json:"titlePt"`
	TitleHi           string                        `json:"titleHi"`
	TitleId           string                        `json:"titleId"`
	RechargeBtnTextEn string                        `json:"rechargeBtnTextEn"`
	RechargeBtnTextEs string                        `json:"rechargeBtnTextEs"`
	RechargeBtnTextPt string                        `json:"rechargeBtnTextPt"`
	RechargeBtnTextHi string                        `json:"rechargeBtnTextHi"`
	RechargeBtnTextId string                        `json:"rechargeBtnTextId"`
	Privileges        []*FirstRechargePrivilegeItem `json:"privileges"`
}

type SaveFirstRechargeActivityCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type AppFirstRechargePrivilegeItem struct {
	Icon   string `json:"icon" dc:"特权图标CDN完整URL"`
	DescEn string `json:"descEn" dc:"特权描述(英文)"`
	DescEs string `json:"descEs" dc:"特权描述(西班牙语)"`
	DescPt string `json:"descPt" dc:"特权描述(葡萄牙语)"`
	DescHi string `json:"descHi" dc:"特权描述(印地语)"`
	DescId string `json:"descId" dc:"特权描述(印尼语)"`
}

type AppFirstRechargeActivityCfgReq struct {
	g.Meta `path:"/firstRechargeActivityCfgForApp" method:"post" summary:"App查询首充活动配置(读内存缓存)" tags:"首充活动配置"`
}

// AppFirstRechargeActivityCfgRes App端首充活动配置
// enabled=false 时建议隐藏活动入口; 文案按客户端语言从 En/Es/Pt/Hi/Id 字段取值
// FirstRechargeSuccessPushItem 首充成功推送载荷
type FirstRechargeSuccessPushItem struct {
	FirstRecharge bool   `json:"firstRecharge" dc:"false=已首充,客户端隐藏首充入口"`
	Gold          uint64 `json:"gold" dc:"本次到账金币数"`
	OrderId       string `json:"orderId" dc:"充值订单ID"`
}

type AppFirstRechargeActivityCfgRes struct {
	Enabled           bool                             `json:"enabled" dc:"活动开关,true=展示首充活动,false=隐藏"`
	Icon              string                           `json:"icon" dc:"活动小图标CDN完整URL"`
	TitleEn           string                           `json:"titleEn" dc:"活动标题(英文)"`
	TitleEs           string                           `json:"titleEs" dc:"活动标题(西班牙语)"`
	TitlePt           string                           `json:"titlePt" dc:"活动标题(葡萄牙语)"`
	TitleHi           string                           `json:"titleHi" dc:"活动标题(印地语)"`
	TitleId           string                           `json:"titleId" dc:"活动标题(印尼语)"`
	RechargeBtnTextEn string                           `json:"rechargeBtnTextEn" dc:"充值按钮文案(英文)"`
	RechargeBtnTextEs string                           `json:"rechargeBtnTextEs" dc:"充值按钮文案(西班牙语)"`
	RechargeBtnTextPt string                           `json:"rechargeBtnTextPt" dc:"充值按钮文案(葡萄牙语)"`
	RechargeBtnTextHi string                           `json:"rechargeBtnTextHi" dc:"充值按钮文案(印地语)"`
	RechargeBtnTextId string                           `json:"rechargeBtnTextId" dc:"充值按钮文案(印尼语)"`
	Privileges        []*AppFirstRechargePrivilegeItem `json:"privileges" dc:"特权列表(按CMS配置顺序),icon为CDN URL,描述按语言取值"`
}

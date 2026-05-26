package rechargeorderdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// ===== 公共订单项 =====

// RechargeOrderItem 充值订单单条(CMS / App 共用)
type RechargeOrderItem struct {
	ID           string  `json:"id"`
	UserId       uint64  `json:"userId"`
	CfgId        uint64  `json:"cfgId"`
	Price        float64 `json:"price"`    // 实付金额(USD)
	Currency     string  `json:"currency"` // 币种
	Gold         float64 `json:"gold"`     // 充值发放金币数
	Status       uint8   `json:"status"`   // 0待支付 1已完成 2已取消
	Source       uint8   `json:"source"`   // 1App 2后台手动
	PayChannel   uint8   `json:"payChannel"`
	ThirdOrderId string  `json:"thirdOrderId"`
	Remark       string  `json:"remark"`
	OperatorId   uint64  `json:"operatorId"`
	CreatedAt    int64   `json:"createdAt"` // 秒
	PaidAt       int64   `json:"paidAt"`    // 秒,0=未支付
}

// ===== CMS =====

// CMSRechargeOrderListReq CMS分页查询充值订单
type CMSRechargeOrderListReq struct {
	g.Meta `path:"/rechargeOrderList" method:"post" summary:"分页查询充值订单" tags:"充值订单"`
	httpserver.CMSQueryReq
	UserId       string `json:"userId"       dc:"按用户ID过滤(空=全部)"`
	OrderId      string `json:"orderId"      dc:"按订单ID精确查询(空=不过滤)"`
	StatusFilter int    `json:"statusFilter" dc:"状态过滤(0=全部,1=待支付,2=已完成,3=已取消)"`
	Source       int    `json:"source"       dc:"来源过滤(0=全部,1=App,2=后台手动)"`
	StartTime    int64  `json:"startTime"    dc:"创建时间起(秒, 0=不过滤)"`
	EndTime      int64  `json:"endTime"      dc:"创建时间止(秒, 0=不过滤)"`
}

// CMSRechargeOrderListItem CMS充值订单列表项
type CMSRechargeOrderListItem struct {
	ID           string  `json:"id"`
	UserId       string  `json:"userId"`
	Nickname     string  `json:"nickname"`
	CfgId        string  `json:"cfgId"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	Gold         float64 `json:"gold"`
	Status       uint8   `json:"status"`
	Source       uint8   `json:"source"`
	PayChannel   uint8   `json:"payChannel"`
	ThirdOrderId string  `json:"thirdOrderId"`
	Remark       string  `json:"remark"`
	OperatorId   string  `json:"operatorId"`
	CreatedAt    int64   `json:"createdAt"`
	PaidAt       int64   `json:"paidAt"`
}

// CMSManualRechargeReq 后台手动给玩家充值到账
// 同时支持两种方式(二选一):
//  1. 引用充值配置(CfgId>0): 自动取 Price/Gold/Currency
//  2. 手动输入(CfgId=0): 必须提供 Price 与 Gold,Currency 可选(默认 CNY)
type CMSManualRechargeReq struct {
	g.Meta  `path:"/manualRecharge" method:"post" summary:"后台手动给玩家充值到账" tags:"充值订单"`
	OrderId string `json:"orderId" v:"required#订单ID不能为空" dc:"订单ID"`
}

type CMSManualRechargeRes struct {
	OrderId string  `json:"orderId"`
	Gold    float64 `json:"gold"`  // 实际发放金币数
	After   float64 `json:"after"` // 玩家金币变更后余额
	Success bool    `json:"success"`
}

// CMSCreateRechargeOrderReq 后台人工创建充值订单(待支付)
type CMSCreateRechargeOrderReq struct {
	g.Meta `path:"/manualCreateOrder" method:"post" summary:"后台人工创建充值订单" tags:"充值订单"`
	UserId string  `json:"userId" v:"required#玩家ID不能为空" dc:"玩家用户ID"`
	Amount float64 `json:"amount" v:"required#订单金额不能为空" dc:"订单金额(USD)"`
}

type CMSCreateRechargeOrderRes struct {
	OrderId  string  `json:"orderId"`
	Price    float64 `json:"price"`
	Gold     float64 `json:"gold"`
	Currency string  `json:"currency"`
	Status   uint8   `json:"status"`
	Success  bool    `json:"success"`
}

// ===== App =====

// AppCreateRechargeOrderReq App 端创建充值订单(根据 cfgId 取价格与金币)
// 不允许 App 自由指定金额,防止刷单
type AppCreateRechargeOrderReq struct {
	g.Meta     `path:"/createRechargeOrder" method:"post" summary:"App创建充值订单" tags:"充值订单"`
	CfgId      uint64  `json:"cfgId"       dc:"充值档位ID(来自已上架的 rechargeCfg)"`
	PayChannel uint8   `json:"payChannel"  dc:"支付渠道(wechat/alipay/applepay 等)"`
	Amount     float64 `json:"amount" dc:"自定义金额"`
}

type AppCreateRechargeOrderRes struct {
	OrderId  string  `json:"orderId"`  // 订单ID
	Price    float64 `json:"price"`    // 应付金额(USD)
	Currency string  `json:"currency"` // 币种
	Status   uint8   `json:"status"`   // 订单状态(创建后=0 待支付)
}

// AppMyRechargeOrderListReq App 端查询自己的充值记录(分页)
type AppMyRechargeOrderListReq struct {
	g.Meta       `path:"/myRechargeOrderList" method:"post" summary:"App查询本人充值订单列表" tags:"充值订单"`
	PageIndex    int `json:"pageIndex"    dc:"页码(从1开始)"`
	PageSize     int `json:"pageSize"     dc:"每页数量"`
	StatusFilter int `json:"statusFilter" dc:"状态过滤(0=全部,1=待支付,2=已完成,3=已取消)"`
}

type AppMyRechargeOrderListRes struct {
	Total int                  `json:"total"`
	List  []*RechargeOrderItem `json:"list"`
}

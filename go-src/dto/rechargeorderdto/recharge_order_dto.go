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
	Price        uint64  `json:"price"`    // 实付金额(分)
	Currency     string  `json:"currency"` // 币种
	Gold         float64 `json:"gold"`     // 充值发放金币数
	Status       uint8   `json:"status"`   // 0待支付 1已完成 2已取消
	Source       uint8   `json:"source"`   // 1App 2后台手动
	PayChannel   string  `json:"payChannel"`
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
	UserId       uint64 `json:"userId"       dc:"按用户ID过滤(0=全部)"`
	OrderId      uint64 `json:"orderId"      dc:"按订单ID精确查询(0=不过滤)"`
	StatusFilter int    `json:"statusFilter" dc:"状态过滤(0=全部,1=待支付,2=已完成,3=已取消)"`
	Source       int    `json:"source"       dc:"来源过滤(0=全部,1=App,2=后台手动)"`
	StartTime    int64  `json:"startTime"    dc:"创建时间起(秒, 0=不过滤)"`
	EndTime      int64  `json:"endTime"      dc:"创建时间止(秒, 0=不过滤)"`
}

// CMSManualRechargeReq 后台手动给玩家充值到账
// 同时支持两种方式(二选一):
//  1. 引用充值配置(CfgId>0): 自动取 Price/Gold/Currency
//  2. 手动输入(CfgId=0): 必须提供 Price 与 Gold,Currency 可选(默认 CNY)
type CMSManualRechargeReq struct {
	g.Meta   `path:"/manualRecharge" method:"post" summary:"后台手动给玩家充值到账" tags:"充值订单"`
	UserId   uint64  `json:"userId"   v:"required|min:1#玩家ID不能为空|玩家ID非法" dc:"目标玩家ID"`
	CfgId    uint64  `json:"cfgId"    dc:"充值档位ID,>0时从配置取价格与金币;=0时使用 price/gold"`
	Price    uint64  `json:"price"    dc:"实付金额(单位:分),CfgId=0时必填"`
	Currency string  `json:"currency" v:"max-length:8#货币代码最长8字符" dc:"币种(默认CNY)"`
	Gold     float64 `json:"gold"     dc:"发放金币数,CfgId=0时必填"`
	Remark   string  `json:"remark"   v:"max-length:255#备注最长255字符" dc:"备注"`
}

type CMSManualRechargeRes struct {
	OrderId string  `json:"orderId"`
	Gold    float64 `json:"gold"`  // 实际发放金币数
	After   float64 `json:"after"` // 玩家金币变更后余额
	Success bool    `json:"success"`
}

// ===== App =====

// AppCreateRechargeOrderReq App 端创建充值订单(根据 cfgId 取价格与金币)
// 不允许 App 自由指定金额,防止刷单
type AppCreateRechargeOrderReq struct {
	g.Meta     `path:"/createRechargeOrder" method:"post" summary:"App创建充值订单" tags:"充值订单"`
	CfgId      uint64 `json:"cfgId"      v:"required|min:1#档位ID不能为空|档位ID非法" dc:"充值档位ID(来自已上架的 rechargeCfg)"`
	PayChannel string `json:"payChannel" v:"max-length:32#支付渠道最长32字符" dc:"支付渠道(wechat/alipay/applepay 等)"`
}

type AppCreateRechargeOrderRes struct {
	OrderId  string `json:"orderId"`  // 订单ID
	Price    uint64 `json:"price"`    // 应付金额(分)
	Currency string `json:"currency"` // 币种
	Status   uint8  `json:"status"`   // 订单状态(创建后=0 待支付)
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

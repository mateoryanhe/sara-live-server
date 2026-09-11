package gameevent

import (
	"xr-game-server/constants/currency"
	"xr-game-server/core/event"
	"xr-game-server/entity/recharge"
)

const (
	// UsdIncomeArrivedEvent 美金入账(订单待发币阶段发布;含白名单;用 Kind 区分类型)
	// 含:人工确认 / Google / iOS / 渠道 / 币商 / 充值白名单
	UsdIncomeArrivedEvent event.Type = "UsdIncomeArrivedEvent"

	// RechargeGoldArrivedEvent 充值金币到账(wallet 处理加赠并发币后发布)
	// IsAccountFirst / IsTierFirst 表示本笔是否账号首充、档位首充
	RechargeGoldArrivedEvent event.Type = "RechargeGoldArrivedEvent"
)

// UsdIncomeKind 美金入账类型
type UsdIncomeKind uint8

const (
	UsdIncomeKindUnknown      UsdIncomeKind = 0
	UsdIncomeKindManual       UsdIncomeKind = 1 // 人工确认订单
	UsdIncomeKindGoogle       UsdIncomeKind = 2
	UsdIncomeKindIOS          UsdIncomeKind = 3
	UsdIncomeKindChannel      UsdIncomeKind = 4 // 渠道支付
	UsdIncomeKindCoinMerchant UsdIncomeKind = 5 // 币商
	UsdIncomeKindWhitelist    UsdIncomeKind = 6 // 充值白名单(虚拟美金)
)

// UsdIncomeArrivedEventData 美金入账事件载荷
// wallet 订阅后完成加赠发币,回填 Granted/GoldBalance 等字段供 completeOrder 使用。
type UsdIncomeArrivedEventData struct {
	Order     *entity.RechargeOrder
	UsdAmount float64
	Kind      UsdIncomeKind
	// Reason 发币流水原因;为 ReasonRecharge 时由 wallet 按订单细分
	Reason currency.Reason

	// --- wallet 处理后回填 ---
	Granted        bool
	BaseGold       float64
	CreditedGold   float64
	GoldBalance    float64
	IsTierFirst    bool
	IsAccountFirst bool
}

func NewUsdIncomeArrivedEventData(order *entity.RechargeOrder, kind UsdIncomeKind, reason currency.Reason) *UsdIncomeArrivedEventData {
	if order == nil {
		return &UsdIncomeArrivedEventData{Kind: kind, Reason: reason}
	}
	return &UsdIncomeArrivedEventData{
		Order:     order,
		UsdAmount: order.Price,
		Kind:      kind,
		Reason:    reason,
	}
}

// IsVirtual 是否虚拟美金入账(充值白名单)
func (d *UsdIncomeArrivedEventData) IsVirtual() bool {
	return d != nil && d.Kind == UsdIncomeKindWhitelist
}

// RechargeGoldArrivedEventData 充值金币到账事件载荷
type RechargeGoldArrivedEventData struct {
	Order          *entity.RechargeOrder
	Kind           UsdIncomeKind
	BaseGold       float64 // 加赠前
	CreditedGold   float64 // 实际到账(含加赠)
	BonusGold      float64 // 加赠部分
	GoldBalance    float64 // 发币后余额
	IsTierFirst    bool    // 档位首充(含加赠)
	IsAccountFirst bool    // 账号首充
}

func NewRechargeGoldArrivedEventData(
	order *entity.RechargeOrder,
	kind UsdIncomeKind,
	baseGold, creditedGold, goldBalance float64,
	isTierFirst, isAccountFirst bool,
) *RechargeGoldArrivedEventData {
	bonus := creditedGold - baseGold
	if bonus < 0 {
		bonus = 0
	}
	return &RechargeGoldArrivedEventData{
		Order:          order,
		Kind:           kind,
		BaseGold:       baseGold,
		CreditedGold:   creditedGold,
		BonusGold:      bonus,
		GoldBalance:    goldBalance,
		IsTierFirst:    isTierFirst,
		IsAccountFirst: isAccountFirst,
	}
}

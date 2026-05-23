package gameevent

import (
	"time"
	"xr-game-server/core/event"
)

const (
	// RechargeArrivedEvent 充值到账事件(订单完成且金币已发放)
	RechargeArrivedEvent event.Type = "RechargeArrivedEvent"
)

// RechargeArrivedEventData 充值到账事件数据
type RechargeArrivedEventData struct {
	OrderId      uint64  // 充值订单ID
	UserId       uint64  // 充值用户ID
	CfgId        uint64  // 充值档位ID(0=手动输入金额)
	Price        float64 // 实付金额(单位:分)
	Currency     string  // 结算货币
	Gold         float64 // 到账金币数
	Source       uint8   // 来源(1-App,2-后台手动)
	PayChannel   string  // 支付渠道
	ThirdOrderId string  // 第三方订单号
	AfterGold    float64 // 到账后玩家金币余额
	PaidAt       int64   // 到账时间(Unix秒)
}

func NewRechargeArrivedEventData(
	orderId, userId, cfgId uint64,
	price float64,
	currency string,
	gold float64,
	source uint8,
	payChannel, thirdOrderId string,
	afterGold float64,
	paidAt time.Time,
) *RechargeArrivedEventData {
	paidAtUnix := int64(0)
	if !paidAt.IsZero() {
		paidAtUnix = paidAt.Unix()
	}
	return &RechargeArrivedEventData{
		OrderId:      orderId,
		UserId:       userId,
		CfgId:        cfgId,
		Price:        price,
		Currency:     currency,
		Gold:         gold,
		Source:       source,
		PayChannel:   payChannel,
		ThirdOrderId: thirdOrderId,
		AfterGold:    afterGold,
		PaidAt:       paidAtUnix,
	}
}

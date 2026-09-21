package call

import (
	"strconv"
	"time"

	"xr-game-server/constants/currency"
	"xr-game-server/constants/liverevenue"
	"xr-game-server/core/event"
	"xr-game-server/core/math"
	"xr-game-server/dao/liveroomdao"
	callentity "xr-game-server/entity/call"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/gameevent"
	"xr-game-server/module/livecfg"
	"xr-game-server/module/liveroom"
	"xr-game-server/module/wallet"
)

// callDiamondBillingPartyIds 从订单保存的付费者确定主播与付费用户。
// payer_id=0 仅用于兼容字段上线前尚未结束的旧订单。
func callDiamondBillingPartyIds(order *callentity.CallOrder) (anchorId, payerId uint64, billable bool) {
	if order == nil {
		return 0, 0, false
	}
	switch order.Source {
	case callentity.CallOrderSourceLiveRoom:
		if order.CallerId == 0 || order.ReceiverId == 0 {
			return 0, 0, false
		}
		payerId = order.PayerId
		if payerId == 0 {
			payerId = order.CallerId
		}
		if payerId != order.CallerId {
			return 0, 0, false
		}
		return order.ReceiverId, payerId, true
	case callentity.CallOrderSourceOneToOneRoom:
		payerId = order.PayerId
		if payerId == 0 {
			resolvedAnchorId, resolvedPayerId, ok := resolveRoomCallPartyIds(order.CallerId, order.ReceiverId)
			if !ok {
				return 0, 0, false
			}
			return resolvedAnchorId, resolvedPayerId, true
		}
		switch payerId {
		case order.CallerId:
			anchorId = order.ReceiverId
		case order.ReceiverId:
			anchorId = order.CallerId
		default:
			return 0, 0, false
		}
		return anchorId, payerId, anchorId > 0
	default:
		return 0, 0, false
	}
}

func callDiamondBillingUserId(order *callentity.CallOrder) (uint64, bool) {
	if order == nil {
		return 0, false
	}
	_, payerId, billable := callDiamondBillingPartyIds(order)
	return payerId, billable
}

func liveRoomCallTicketPrice(order *callentity.CallOrder, enabled bool) float64 {
	if !enabled || order == nil || order.Source != callentity.CallOrderSourceLiveRoom || order.TicketPrice <= 0 {
		return 0
	}
	return order.TicketPrice
}

// checkLiveRoomCallDiamondOnAccept 接通时按来源校验付费用户是否可支付门票与首分钟费用。
func checkLiveRoomCallDiamondOnAccept(order *callentity.CallOrder) error {
	payerId, billable := callDiamondBillingUserId(order)
	if !billable {
		return nil
	}

	requiredDiamond := order.PricePerMinute
	requiredDiamond = math.AddFloat64(requiredDiamond, liveRoomCallTicketPrice(order, livecfg.IsVideoCallTicketEnabled()))
	if requiredDiamond <= 0 {
		return nil
	}
	return wallet.CanPayWithGoldExchange(payerId, requiredDiamond)
}

// chargeLiveRoomCallOnAccept 双方确认接通后扣费。
// source=1 按开关扣门票并扣首分钟；source=3 只扣首分钟。
func chargeLiveRoomCallOnAccept(order *callentity.CallOrder, now time.Time) error {
	if order == nil {
		return nil
	}
	anchorId, payerId, billable := callDiamondBillingPartyIds(order)
	if !billable {
		return nil
	}

	liveRecordId, _ := strconv.ParseUint(order.Params, 10, 64)
	var totalCost float64
	ticketPrice := liveRoomCallTicketPrice(order, livecfg.IsVideoCallTicketEnabled())

	if ticketPrice > 0 {
		if _, err := wallet.DiamondSubWithGoldExchange(payerId, ticketPrice, currency.ReasonLiveRoomVideoCallTicket); err != nil {
			return err
		}
		totalCost = math.AddFloat64(totalCost, ticketPrice)
		recordLiveRoomCallRevenue(anchorId, liveRecordId, payerId, order.ID, ticketPrice, liverevenue.LiveRoomVideoCallTicket)
	}

	if order.PricePerMinute > 0 {
		if _, err := wallet.DiamondSubWithGoldExchange(payerId, order.PricePerMinute, currency.ReasonLiveRoomVideoCallBilling); err != nil {
			return err
		}
		totalCost = math.AddFloat64(totalCost, order.PricePerMinute)
		recordLiveRoomCallRevenue(anchorId, liveRecordId, payerId, order.ID, order.PricePerMinute, liverevenue.LiveRoomVideoCallBilling)
		order.AddBillingDuration(1)
		nextCharge := now.Add(time.Minute)
		order.SetChargeTime(&nextCharge)
	}

	order.SetTotalCost(totalCost)
	return nil
}

// chargeLiveRoomCallBillingIfDue 通话按来源进行分钟续费(心跳触发,加锁避免双方重复扣费)
func chargeLiveRoomCallBillingIfDue(order *callentity.CallOrder, now time.Time) error {
	if order == nil {
		return nil
	}
	anchorId, payerId, billable := callDiamondBillingPartyIds(order)
	if !billable {
		return nil
	}
	if !order.IsCallStarted() || order.PricePerMinute <= 0 {
		return nil
	}
	if order.ChargeTime == nil || now.Before(*order.ChargeTime) {
		return nil
	}

	liveRecordId, _ := strconv.ParseUint(order.Params, 10, 64)
	if _, err := wallet.DiamondSubWithGoldExchange(payerId, order.PricePerMinute, currency.ReasonLiveRoomVideoCallBilling); err != nil {
		return err
	}
	recordLiveRoomCallRevenue(anchorId, liveRecordId, payerId, order.ID, order.PricePerMinute, liverevenue.LiveRoomVideoCallBilling)
	order.SetTotalCost(math.AddFloat64(order.TotalCost, order.PricePerMinute))
	order.AddBillingDuration(1)
	nextCharge := now.Add(time.Minute)
	order.SetChargeTime(&nextCharge)
	return nil
}

func recordLiveRoomCallRevenue(roomId, liveRecordId, callerId, orderId uint64, amount float64, revenueType liverevenue.Type) {
	applyLiveRoomCallRevenue(roomId, liveRecordId, callerId, orderId, amount, revenueType)
}

func applyLiveRoomCallRevenue(roomId, liveRecordId, callerId, orderId uint64, amount float64, revenueType liverevenue.Type) {
	if amount <= 0 || roomId == 0 {
		return
	}

	room := liveroomdao.GetRoomById(roomId)
	if room == nil {
		return
	}
	ticket := revenueType == liverevenue.LiveRoomVideoCallTicket
	billing := revenueType == liverevenue.LiveRoomVideoCallBilling
	if !ticket && !billing {
		return
	}
	if liveRecordId > 0 {
		if liveRecord := liveroomdao.GetLiveRecordById(liveRecordId); liveRecord != nil {
			liveRecord.ApplyVideoCallIncomeDelta(amount, ticket, billing)
			liveroomdao.PublishLiveRecord(liveRecord)
		}
	}
	applyRoomCallRevenueDelta(room, amount, ticket, billing)

	count := 1
	unitPrice := amount
	if ticket {
		count = 0
		unitPrice = 0
	}
	eventData := liveentity.NewLiveRevenueLogRecord(roomId, liveRecordId, callerId, orderId, count, unitPrice, amount, uint8(revenueType))
	event.Pub(gameevent.RevenueEventEvent, eventData)
	if liveRecordId > 0 {
		liveroom.NotifyLiveRecordTotalIncome(room)
	}
}

func refundLiveRoomCallRevenue(order *callentity.CallOrder, anchorId, payerId, liveRecordId uint64, refundAmount float64) {
	if order == nil || anchorId == 0 || payerId == 0 || refundAmount <= 0 {
		return
	}

	if log := liveroomdao.FindLatestUnrefundedVideoCallBillingLog(order.ID, payerId); log != nil {
		log.SetStatus(liveentity.LiveRevenueLogStatusRefunded)
		liveroomdao.PublishRevenueLog(log)
	}

	room := liveroomdao.GetRoomById(anchorId)
	if room == nil {
		return
	}
	if liveRecordId > 0 {
		if liveRecord := liveroomdao.GetLiveRecordById(liveRecordId); liveRecord != nil {
			liveRecord.ApplyVideoCallIncomeDelta(-refundAmount, false, true)
			liveroomdao.PublishLiveRecord(liveRecord)
		}
	}
	applyRoomCallRevenueDelta(room, -refundAmount, false, true)
	if liveRecordId > 0 {
		liveroom.NotifyLiveRecordTotalIncome(room)
	}
}

func applyRoomCallRevenueDelta(room *liveentity.LiveRoom, amount float64, ticket, billing bool) {
	unsettled := liveroomdao.GetLiveRoomIncomeUnsettled(room.ID)
	total := liveroomdao.GetLiveRoomIncomeTotal(room.ID)
	if unsettled == nil || total == nil {
		return
	}
	liveentity.ApplyVideoCallIncomeDelta(liveentity.TbLiveRoomIncomeUnsettled, unsettled.ID, &unsettled.LiveRoomIncomeAmounts, &unsettled.UpdatedAt, amount, ticket, billing)
	liveentity.ApplyVideoCallIncomeDelta(liveentity.TbLiveRoomIncomeTotal, total.ID, &total.LiveRoomIncomeAmounts, &total.UpdatedAt, amount, ticket, billing)
	liveroomdao.MirrorGuildVideoCallIncomeDelta(room.ID, amount, ticket, billing)
	liveroomdao.MirrorDailyAnchorVideoCallIncomeDelta(room.ID, time.Now(), amount, ticket, billing)
}

const callBillingRefundGrace = 30 * time.Second

// refundLiveRoomCallLastMinuteIfNeeded 结束通话时,若未超过ChargeTime 30秒则退回最后一次分钟扣费
func refundLiveRoomCallLastMinuteIfNeeded(order *callentity.CallOrder, endTime time.Time) error {
	if order == nil {
		return nil
	}
	anchorId, payerId, billable := callDiamondBillingPartyIds(order)
	if !billable {
		return nil
	}
	if order.ChargeTime == nil || order.PricePerMinute <= 0 || order.BillingDuration == 0 {
		return nil
	}
	if endTime.Add(callBillingRefundGrace).After(*order.ChargeTime) {
		return nil
	}

	refundAmount := order.PricePerMinute
	if _, err := wallet.DiamondAdd(payerId, refundAmount, currency.ReasonRefund); err != nil {
		return err
	}

	liveRecordId, _ := strconv.ParseUint(order.Params, 10, 64)
	refundLiveRoomCallRevenue(order, anchorId, payerId, liveRecordId, refundAmount)

	order.SetTotalCost(math.SubFloat64(order.TotalCost, refundAmount))
	order.SubBillingDuration(1)
	return nil
}

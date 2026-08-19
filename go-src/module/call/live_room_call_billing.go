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
	"xr-game-server/module/liveroom"
	"xr-game-server/module/wallet"
)

// checkLiveRoomCallDiamondOnAccept 接听时校验呼叫者是否可支付(钻石+按需兑换金币,门票+首分钟)
func checkLiveRoomCallDiamondOnAccept(order *callentity.CallOrder) error {
	if order == nil || order.Source != callentity.CallOrderSourceLiveRoom {
		return nil
	}

	requiredDiamond := order.TicketPrice
	if order.PricePerMinute > 0 {
		requiredDiamond += order.PricePerMinute
	}
	if requiredDiamond <= 0 {
		return nil
	}
	return wallet.CanPayWithGoldExchange(order.CallerId, requiredDiamond)
}

// chargeLiveRoomCallOnAccept 直播间来源通话接听后开始扣费(门票+首分钟)
func chargeLiveRoomCallOnAccept(order *callentity.CallOrder, now time.Time) error {
	if order == nil || order.Source != callentity.CallOrderSourceLiveRoom {
		return nil
	}

	liveRecordId, _ := strconv.ParseUint(order.Params, 10, 64)
	roomId := order.ReceiverId
	callerId := order.CallerId

	var totalCost float64

	if order.TicketPrice > 0 {
		if _, err := wallet.DiamondSubWithGoldExchange(callerId, order.TicketPrice, currency.ReasonLiveRoomVideoCallTicket); err != nil {
			return err
		}
		totalCost += order.TicketPrice
		recordLiveRoomCallRevenue(roomId, liveRecordId, callerId, order.ID, order.TicketPrice, liverevenue.LiveRoomVideoCallTicket)
	}

	if order.PricePerMinute > 0 {
		if _, err := wallet.DiamondSubWithGoldExchange(callerId, order.PricePerMinute, currency.ReasonLiveRoomVideoCallBilling); err != nil {
			return err
		}
		totalCost += order.PricePerMinute
		recordLiveRoomCallRevenue(roomId, liveRecordId, callerId, order.ID, order.PricePerMinute, liverevenue.LiveRoomVideoCallBilling)
		order.AddBillingDuration(1)
		nextCharge := now.Add(time.Minute)
		order.SetChargeTime(&nextCharge)
	}

	order.SetTotalCost(totalCost)
	return nil
}

// chargeLiveRoomCallBillingIfDue 直播间通话按分钟续费(心跳触发,加锁避免双方重复扣费)
func chargeLiveRoomCallBillingIfDue(order *callentity.CallOrder, now time.Time) error {
	if order == nil || order.Source != callentity.CallOrderSourceLiveRoom {
		return nil
	}
	if !order.IsCallStarted() || order.PricePerMinute <= 0 {
		return nil
	}
	if order.ChargeTime == nil || now.Before(*order.ChargeTime) {
		return nil
	}

	liveRecordId, _ := strconv.ParseUint(order.Params, 10, 64)
	if _, err := wallet.DiamondSubWithGoldExchange(order.CallerId, order.PricePerMinute, currency.ReasonLiveRoomVideoCallBilling); err != nil {
		return err
	}
	recordLiveRoomCallRevenue(order.ReceiverId, liveRecordId, order.CallerId, order.ID, order.PricePerMinute, liverevenue.LiveRoomVideoCallBilling)
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
	if amount <= 0 || liveRecordId == 0 {
		return
	}

	room := liveroomdao.GetRoomById(roomId)
	if room == nil {
		return
	}
	ticket := revenueType == liverevenue.LiveRoomVideoCallTicket
	billing := revenueType == liverevenue.LiveRoomVideoCallBilling
	if liveRecord := liveroomdao.GetLiveRecordById(liveRecordId); liveRecord != nil {
		liveRecord.ApplyVideoCallIncomeDelta(amount, ticket, billing)
	}
	applyRoomCallRevenueDelta(room, amount, ticket, billing)

	var eventData *liveentity.LiveRevenueLog
	switch revenueType {
	case liverevenue.LiveRoomVideoCallTicket:
		eventData = liveentity.NewLiveRevenueLogRecord(roomId, liveRecordId, callerId, roomId, orderId, 0, 0, amount, uint8(liverevenue.LiveRoomVideoCallTicket))
	case liverevenue.LiveRoomVideoCallBilling:
		eventData = liveentity.NewLiveRevenueLogRecord(roomId, liveRecordId, callerId, roomId, orderId, 1, amount, amount, uint8(liverevenue.LiveRoomVideoCallBilling))
	default:
		return
	}
	event.Pub(gameevent.RevenueEventEvent, eventData)
	liveroom.NotifyLiveRecordTotalIncome(room)
}

func refundLiveRoomCallRevenue(order *callentity.CallOrder, liveRecordId uint64, refundAmount float64) {
	if order == nil || refundAmount <= 0 || liveRecordId == 0 {
		return
	}

	if log := liveroomdao.FindLatestUnrefundedVideoCallBillingLog(order.ID, order.CallerId); log != nil {
		log.SetStatus(liveentity.LiveRevenueLogStatusRefunded)
	}

	room := liveroomdao.GetRoomById(order.ReceiverId)
	if room == nil {
		return
	}
	if liveRecord := liveroomdao.GetLiveRecordById(liveRecordId); liveRecord != nil {
		liveRecord.ApplyVideoCallIncomeDelta(-refundAmount, false, true)
	}
	applyRoomCallRevenueDelta(room, -refundAmount, false, true)
	liveroom.NotifyLiveRecordTotalIncome(room)
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
	if order == nil || order.Source != callentity.CallOrderSourceLiveRoom {
		return nil
	}
	if order.ChargeTime == nil || order.PricePerMinute <= 0 || order.BillingDuration == 0 {
		return nil
	}
	if endTime.Add(callBillingRefundGrace).After(*order.ChargeTime) {
		return nil
	}

	refundAmount := order.PricePerMinute
	if _, err := wallet.DiamondAdd(order.CallerId, refundAmount, currency.ReasonRefund); err != nil {
		return err
	}

	liveRecordId, _ := strconv.ParseUint(order.Params, 10, 64)
	refundLiveRoomCallRevenue(order, liveRecordId, refundAmount)

	order.SetTotalCost(math.SubFloat64(order.TotalCost, refundAmount))
	order.SubBillingDuration(1)
	return nil
}

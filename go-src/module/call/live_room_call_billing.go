package call

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/currency"
	"xr-game-server/constants/liverevenue"
	"xr-game-server/core/event"
	"xr-game-server/core/math"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
	"xr-game-server/module/wallet"
)

// checkLiveRoomCallDiamondOnAccept 接听时仅校验呼叫者钻石是否足够(门票+首分钟)
func checkLiveRoomCallDiamondOnAccept(order *entity.CallOrder) error {
	if order == nil || order.Source != entity.CallOrderSourceLiveRoom {
		return nil
	}

	requiredDiamond := order.TicketPrice
	if order.PricePerMinute > 0 {
		requiredDiamond += order.PricePerMinute
	}
	if requiredDiamond <= 0 {
		return nil
	}
	return wallet.DiamondNotEnough(order.CallerId, requiredDiamond)
}

// chargeLiveRoomCallOnAccept 直播间来源通话接听后开始扣费(门票+首分钟)
func chargeLiveRoomCallOnAccept(order *entity.CallOrder, now time.Time) error {
	if order == nil || order.Source != entity.CallOrderSourceLiveRoom {
		return nil
	}

	liveRecordId, _ := strconv.ParseUint(order.Params, 10, 64)
	roomId := order.ReceiverId
	callerId := order.CallerId

	var totalCost float64

	if order.TicketPrice > 0 {
		if _, err := wallet.DiamondSub(callerId, order.TicketPrice, currency.ReasonLiveRoomVideoCallTicket); err != nil {
			return err
		}
		totalCost += order.TicketPrice
		recordLiveRoomCallRevenue(roomId, liveRecordId, callerId, order.TicketPrice, liverevenue.LiveRoomVideoCallTicket)
	}

	if order.PricePerMinute > 0 {
		if _, err := wallet.DiamondSub(callerId, order.PricePerMinute, currency.ReasonLiveRoomVideoCallBilling); err != nil {
			return err
		}
		totalCost += order.PricePerMinute
		recordLiveRoomCallRevenue(roomId, liveRecordId, callerId, order.PricePerMinute, liverevenue.LiveRoomVideoCallBilling)
		order.AddBillingDuration(1)
		nextCharge := now.Add(time.Minute)
		order.SetChargeTime(&nextCharge)
	}

	order.SetTotalCost(totalCost)
	return nil
}

// chargeLiveRoomCallBillingIfDue 直播间通话按分钟续费(心跳触发,加锁避免双方重复扣费)
func chargeLiveRoomCallBillingIfDue(order *entity.CallOrder, now time.Time) error {
	if order == nil || order.Source != entity.CallOrderSourceLiveRoom {
		return nil
	}
	if !order.IsCallStarted() || order.PricePerMinute <= 0 {
		return nil
	}
	if order.ChargeTime == nil || now.Before(*order.ChargeTime) {
		return nil
	}

	liveRecordId, _ := strconv.ParseUint(order.Params, 10, 64)
	if _, err := wallet.DiamondSub(order.CallerId, order.PricePerMinute, currency.ReasonLiveRoomVideoCallBilling); err != nil {
		return err
	}
	recordLiveRoomCallRevenue(order.ReceiverId, liveRecordId, order.CallerId, order.PricePerMinute, liverevenue.LiveRoomVideoCallBilling)
	order.SetTotalCost(math.AddFloat64(order.TotalCost, order.PricePerMinute))
	order.AddBillingDuration(1)
	nextCharge := now.Add(time.Minute)
	order.SetChargeTime(&nextCharge)
	return nil
}

func recordLiveRoomCallRevenue(roomId, liveRecordId, callerId uint64, amount float64, revenueType liverevenue.Type) {
	if amount <= 0 || liveRecordId == 0 {
		return
	}

	lockKey := fmt.Sprintf("room_%d", roomId)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	room := liveroomdao.GetRoomById(roomId)
	if room == nil {
		return
	}
	if liveRecord := liveroomdao.GetLiveRecordById(liveRecordId); liveRecord != nil {
		liveRecord.AddTotalIncome(amount)
		liveRecord.AddTotalVideoCallIncome(amount)
		switch revenueType {
		case liverevenue.LiveRoomVideoCallTicket:
			liveRecord.AddTotalVideoCallTicketIncome(amount)
		case liverevenue.LiveRoomVideoCallBilling:
			liveRecord.AddTotalVideoCallBillingIncome(amount)
		}
	}
	room.AddTotalIncome(amount)
	switch revenueType {
	case liverevenue.LiveRoomVideoCallTicket:
		room.AddTotalVideoCallTicketIncome(amount)
	case liverevenue.LiveRoomVideoCallBilling:
		room.AddTotalVideoCallBillingIncome(amount)
	}

	var eventData *entity.LiveRevenueLog
	switch revenueType {
	case liverevenue.LiveRoomVideoCallTicket:
		eventData = entity.NewLiveRevenueLogRecord(roomId, liveRecordId, callerId, roomId, 0, 0, 0, amount, uint8(liverevenue.LiveRoomVideoCallTicket))
	case liverevenue.LiveRoomVideoCallBilling:
		eventData = entity.NewLiveRevenueLogRecord(roomId, liveRecordId, callerId, roomId, 0, 1, amount, amount, uint8(liverevenue.LiveRoomVideoCallBilling))
	default:
		return
	}
	event.Pub(gameevent.RevenueEventEvent, eventData)
}

const callBillingRefundGrace = 30 * time.Second

// refundLiveRoomCallLastMinuteIfNeeded 结束通话时,若未超过ChargeTime 30秒则退回最后一次分钟扣费
func refundLiveRoomCallLastMinuteIfNeeded(order *entity.CallOrder, endTime time.Time) error {
	if order == nil || order.Source != entity.CallOrderSourceLiveRoom {
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
	order.SetTotalCost(math.SubFloat64(order.TotalCost, refundAmount))
	order.SubBillingDuration(1)
	return nil
}

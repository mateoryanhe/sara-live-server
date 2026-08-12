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
	"xr-game-server/core/syndb"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
	"xr-game-server/module/wallet"
)

// checkLiveRoomCallDiamondOnAccept 接听时校验呼叫者是否可支付(钻石+按需兑换金币,门票+首分钟)
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
	return wallet.CanPayWithGoldExchange(order.CallerId, requiredDiamond)
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

	lockKey := fmt.Sprintf("room_%d", roomId)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	room := liveroomdao.GetRoomById(roomId)
	if room == nil {
		return
	}
	if liveRecord := liveroomdao.GetLiveRecordById(liveRecordId); liveRecord != nil {
		applyLiveRecordCallRevenueDelta(liveRecord, amount, revenueType)
	}
	applyRoomCallRevenueDelta(room, amount, revenueType)

	var eventData *entity.LiveRevenueLog
	switch revenueType {
	case liverevenue.LiveRoomVideoCallTicket:
		eventData = entity.NewLiveRevenueLogRecord(roomId, liveRecordId, callerId, roomId, orderId, 0, 0, amount, uint8(liverevenue.LiveRoomVideoCallTicket))
	case liverevenue.LiveRoomVideoCallBilling:
		eventData = entity.NewLiveRevenueLogRecord(roomId, liveRecordId, callerId, roomId, orderId, 1, amount, amount, uint8(liverevenue.LiveRoomVideoCallBilling))
	default:
		return
	}
	event.Pub(gameevent.RevenueEventEvent, eventData)
}

func refundLiveRoomCallRevenue(order *entity.CallOrder, liveRecordId uint64, refundAmount float64) {
	if order == nil || refundAmount <= 0 || liveRecordId == 0 {
		return
	}

	lockKey := fmt.Sprintf("room_%d", order.ReceiverId)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	if log := liveroomdao.FindLatestUnrefundedVideoCallBillingLog(order.ID, order.CallerId); log != nil {
		log.SetStatus(entity.LiveRevenueLogStatusRefunded)
	}

	room := liveroomdao.GetRoomById(order.ReceiverId)
	if room == nil {
		return
	}
	if liveRecord := liveroomdao.GetLiveRecordById(liveRecordId); liveRecord != nil {
		applyLiveRecordCallRevenueDelta(liveRecord, -refundAmount, liverevenue.LiveRoomVideoCallBilling)
	}
	applyRoomCallRevenueDelta(room, -refundAmount, liverevenue.LiveRoomVideoCallBilling)
}

func applyLiveRecordCallRevenueDelta(liveRecord *entity.LiveRecord, amount float64, revenueType liverevenue.Type) {
	liveRecord.TotalIncome = math.AddFloat64(liveRecord.TotalIncome, amount)
	syndb.AddData(entity.TbLiveRecord, entity.LiveRecordTotalIncome, &syndb.ColData{
		IdVal: liveRecord.ID, ColVal: liveRecord.TotalIncome,
	})
	liveRecord.TotalVideoCallIncome = math.AddFloat64(liveRecord.TotalVideoCallIncome, amount)
	syndb.AddData(entity.TbLiveRecord, entity.LiveRecordTotalVideoCallIncome, &syndb.ColData{
		IdVal: liveRecord.ID, ColVal: liveRecord.TotalVideoCallIncome,
	})
	switch revenueType {
	case liverevenue.LiveRoomVideoCallTicket:
		liveRecord.TotalVideoCallTicketIncome = math.AddFloat64(liveRecord.TotalVideoCallTicketIncome, amount)
		syndb.AddData(entity.TbLiveRecord, entity.LiveRecordTotalVideoCallTicketIncome, &syndb.ColData{
			IdVal: liveRecord.ID, ColVal: liveRecord.TotalVideoCallTicketIncome,
		})
	case liverevenue.LiveRoomVideoCallBilling:
		liveRecord.TotalVideoCallBillingIncome = math.AddFloat64(liveRecord.TotalVideoCallBillingIncome, amount)
		syndb.AddData(entity.TbLiveRecord, entity.LiveRecordTotalVideoCallBillingIncome, &syndb.ColData{
			IdVal: liveRecord.ID, ColVal: liveRecord.TotalVideoCallBillingIncome,
		})
	}
}

func applyRoomCallRevenueDelta(room *entity.LiveRoom, amount float64, revenueType liverevenue.Type) {
	room.TotalIncome = math.AddFloat64(room.TotalIncome, amount)
	syndb.AddData(entity.TbLiveRoom, entity.LiveRoomTotalIncome, &syndb.ColData{
		IdVal: room.ID, ColVal: room.TotalIncome,
	})
	room.TotalVideoCallIncome = math.AddFloat64(room.TotalVideoCallIncome, amount)
	syndb.AddData(entity.TbLiveRoom, entity.LiveRoomTotalVideoCallIncome, &syndb.ColData{
		IdVal: room.ID, ColVal: room.TotalVideoCallIncome,
	})
	switch revenueType {
	case liverevenue.LiveRoomVideoCallTicket:
		room.TotalVideoCallTicketIncome = math.AddFloat64(room.TotalVideoCallTicketIncome, amount)
		syndb.AddData(entity.TbLiveRoom, entity.LiveRoomTotalVideoCallTicketIncome, &syndb.ColData{
			IdVal: room.ID, ColVal: room.TotalVideoCallTicketIncome,
		})
	case liverevenue.LiveRoomVideoCallBilling:
		room.TotalVideoCallBillingIncome = math.AddFloat64(room.TotalVideoCallBillingIncome, amount)
		syndb.AddData(entity.TbLiveRoom, entity.LiveRoomTotalVideoCallBillingIncome, &syndb.ColData{
			IdVal: room.ID, ColVal: room.TotalVideoCallBillingIncome,
		})
	}
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

	liveRecordId, _ := strconv.ParseUint(order.Params, 10, 64)
	refundLiveRoomCallRevenue(order, liveRecordId, refundAmount)

	order.SetTotalCost(math.SubFloat64(order.TotalCost, refundAmount))
	order.SubBillingDuration(1)
	return nil
}

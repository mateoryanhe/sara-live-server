package liveroom

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gmlock"
	"time"
	"xr-game-server/constants/currency"
	"xr-game-server/constants/liverevenue"
	"xr-game-server/core/event"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/livecfg"
	"xr-game-server/module/wallet"
)

const (
	privateRoomMaxAudience = 1
	BillPeriodTime         = 5
)

// ensureCanJoinPrivateRoom 私密直播间仅允许1名观众(不含主播);已在房间内可重复进入
func ensureCanJoinPrivateRoom(userId uint64, room *entity.LiveRoom) error {
	if room == nil || room.Category != entity.LiveRoomCategoryPrivate {
		return nil
	}
	if userId == room.ID {
		return nil
	}
	if isUserInOnlineMap(userId, room.ID) {
		return nil
	}
	if countAudienceInRoom(room.ID) >= privateRoomMaxAudience {
		return errercode.CreateCode(errercode.LiveRoomPrivateAudienceFull)
	}
	return nil
}

// chargePrivateRoomTicketIfNeeded 私密直播间进房扣门票,24小时内同一用户同一房间只扣一次
func chargePrivateRoomTicketIfNeeded(userId uint64, room *entity.LiveRoom, now time.Time) (float64, error) {
	if room == nil || room.Category != entity.LiveRoomCategoryPrivate {
		return 0, nil
	}
	if userId == room.ID {
		return 0, nil
	}
	ticketPrice := room.Ticket
	if ticketPrice <= 0 {
		return 0, nil
	}

	pay := liveroomdao.GetLiveRoomBillingPay(userId, room.ID)
	if pay == nil {
		return 0, nil
	}
	if pay.IsValidWithin24h() {
		return 0, nil
	}

	if _, err := wallet.DiamondSub(userId, ticketPrice, currency.ReasonPrivateRoomTicket); err != nil {
		return 0, err
	}
	pay.SetLastTicketAt(now)
	pay.SetFreeTime(uint64(livecfg.GetPrivateRoomFreeWatchDuration().Seconds()))

	lockKey := fmt.Sprintf("%d", room.ID)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)
	//防止并发,主播可以收到多个人的礼物
	liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
	//添加本次直播收到的礼物总额
	liveRecord.AddTotalIncome(ticketPrice)
	liveRecord.AddTotalPrivateRoomIncome(ticketPrice)
	//记录主播总收益
	room.AddTotalIncome(ticketPrice)
	room.AddTotalPrivateRoomTicketIncome(ticketPrice)

	//记录直播收益流水(礼物)
	eventData := entity.NewLiveRevenueLogRecord(room.ID, room.LiveRecordId, 0, room.ID, 0, 0, 0, ticketPrice, uint8(liverevenue.Ticket))

	event.Pub(gameevent.RevenueEventEvent, eventData)
	return ticketPrice, nil
}

// chargePrivateRoomBillingIfNeeded 私密直播间按分钟扣观众钻石(每场直播独立计费)
func chargePrivateRoomBillingIfNeeded(userId uint64, room *entity.LiveRoom) (float64, error) {
	if room == nil || room.Category != entity.LiveRoomCategoryPrivate {
		return 0, nil
	}
	if userId == room.ID || room.LiveRecordId == 0 {
		return 0, nil
	}
	billingPrice := room.Billing / 60
	if billingPrice <= 0 {
		return 0, nil
	}

	pay := liveroomdao.GetLiveRoomBillingPay(userId, room.ID)
	if pay == nil {
		return 0, nil
	}
	now := time.Now()
	//开始判断免费时长
	if 0 >= pay.FreeTime {

		//开始预扣费
		if pay.LastPaidAt != nil && now.Before(*pay.LastPaidAt) {
			return 0, nil
		}
		if _, err := wallet.DiamondSub(userId, room.Billing, currency.ReasonPrivateRoomBilling); err != nil {
			return 0, err
		}
		recordPrivateRoomBillingRevenue(room, userId, room.Billing)
		pay.SetLastPaidAt(now.Add(time.Minute))
		return 0, nil

	} else {
		pay.SubFreeTime(PrivatePeriod)
	}
	return billingPrice, nil
}

// 清除私密房免费时间
func clearFreeTime(userId uint64, roomId uint64) {
	room := liveroomdao.GetRoomById(roomId)
	//开始全部计费清0
	if room == nil || room.Category != entity.LiveRoomCategoryPrivate {
		return
	}
	if userId == room.ID {
		return
	}
	pay := liveroomdao.GetLiveRoomBillingPay(userId, room.ID)
	if pay == nil {
		return
	}
	if 0 >= pay.FreeTime {
		return
	}
	onlineId := entity.BuildLiveRoomOnlineId(userId, roomId)
	existing := liveroomdao.GetOnlineById(onlineId, userId, roomId)

	now := time.Now()
	diff := now.Sub(*existing.JoinTime)
	pay.SubFreeTime(uint64(diff.Seconds()))
}

func joinChargePrivateRoom(userId uint64, roomId uint64) error {
	room := liveroomdao.GetRoomById(roomId)
	//开始全部计费清0
	if room == nil || room.Category != entity.LiveRoomCategoryPrivate {
		return nil
	}
	if userId == room.ID {
		return nil
	}
	now := time.Now()

	pay := liveroomdao.GetLiveRoomBillingPay(userId, room.ID)
	if pay == nil {
		return nil
	}
	if pay.FreeTime > 0 {
		return nil
	}
	if _, err := wallet.DiamondSub(userId, room.Billing, currency.ReasonPrivateRoomBilling); err != nil {
		return err
	}
	recordPrivateRoomBillingRevenue(room, userId, room.Billing)
	pay.SetLastPaidAt(now.Add(time.Minute))
	return nil
}

func recordPrivateRoomBillingRevenue(room *entity.LiveRoom, userId uint64, amount float64) {
	if amount <= 0 || room == nil || room.LiveRecordId == 0 {
		return
	}
	if liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId); liveRecord != nil {
		liveRecord.AddTotalIncome(amount)
		liveRecord.AddTotalPrivateRoomIncome(amount)
	}
	room.AddTotalIncome(amount)
	room.AddTotalPrivateRoomWatchIncome(amount)
	eventData := entity.NewLiveRevenueLogRecord(
		room.ID, room.LiveRecordId, userId, room.ID, 0, 1, amount, amount, uint8(liverevenue.PrivateRoom),
	)
	event.Pub(gameevent.RevenueEventEvent, eventData)
}

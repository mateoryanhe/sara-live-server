package minuteactivity

import (
	"time"

	"xr-game-server/gameevent"
)

func onShowcaseJoin(data any) {
	ev, ok := data.(*gameevent.ShowcaseLiveRoomJoinTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.RoomId == 0 || ev.UserId == ev.RoomId {
		return
	}
	schedule(metricJoin, ev.StatAt)
}

func onGameJoin(data any) {
	ev, ok := data.(*gameevent.GameLiveRoomJoinTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.RoomId == 0 || ev.UserId == ev.RoomId {
		return
	}
	schedule(metricJoin, ev.StatAt)
}

func onVoiceJoin(data any) {
	ev, ok := data.(*gameevent.VoiceRoomEnterTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.RoomId == 0 || ev.UserId == ev.RoomId {
		return
	}
	schedule(metricJoin, ev.StatAt)
}

func onShowcaseLeave(data any) {
	ev, ok := data.(*gameevent.ShowcaseLiveRoomLeaveTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.RoomId == 0 {
		return
	}
	schedule(metricLeave, ev.StatAt)
}

func onCall1v1Initiate(data any) {
	scheduleCall1v1InitiateMetric(data)
}

func onCall1v1RoomCall(data any) {
	scheduleCall1v1InitiateMetric(data)
}

func scheduleCall1v1InitiateMetric(data any) {
	switch ev := data.(type) {
	case *gameevent.Call1v1InitiateTrackingEventData:
		if ev == nil || ev.CallerId == 0 || ev.AnchorId == 0 || ev.CallerId == ev.AnchorId {
			return
		}
		schedule(metricCall1v1, ev.StatAt)
	case *gameevent.Call1v1RoomTrackingEventData:
		if ev == nil || ev.CallerId == 0 || ev.AnchorId == 0 || ev.CallerId == ev.AnchorId {
			return
		}
		schedule(metricCall1v1, ev.StatAt)
	}
}

func onCall1v1ConnectSuccess(data any) {
	ev, ok := data.(*gameevent.Call1v1ConnectSuccessTrackingEventData)
	if !ok || ev == nil || ev.CallerId == 0 || ev.ReceiverId == 0 || ev.CallerId == ev.ReceiverId {
		return
	}
	schedule(metricCall1v1ConnectSuccess, ev.StatAt)
}

func onCall1v1BillingSuccess(data any) {
	ev, ok := data.(*gameevent.Call1v1BillingSuccessTrackingEventData)
	if !ok || ev == nil || ev.OrderId == 0 {
		return
	}
	schedule(metricCall1v1Consume, ev.StatAt)
}

func onGameBet(data any) {
	ev, ok := data.(*gameevent.GameBetCreatedEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	schedule(metricBet, time.Now())
}

func onRechargeOrderCreated(data any) {
	ev, ok := data.(*gameevent.RechargeOrderCreatedEventData)
	if !ok || ev == nil || ev.OrderId == 0 {
		return
	}
	schedule(metricRechargeOrder, time.Now())
}

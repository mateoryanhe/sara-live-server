package minuteactivity

import (
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

const retentionDays = 3

func Init() {
	event.Sub(gameevent.DayEvent, onDayCleanup)
	event.Sub(gameevent.ShowcaseLiveRoomJoinTrackingEvent, onShowcaseJoin)
	event.Sub(gameevent.GameLiveRoomJoinTrackingEvent, onGameJoin)
	event.Sub(gameevent.VoiceRoomEnterTrackingEvent, onVoiceJoin)
	event.Sub(gameevent.ShowcaseLiveRoomLeaveTrackingEvent, onShowcaseLeave)
	event.Sub(gameevent.Call1v1InitiateTrackingEvent, onCall1v1Initiate)
	event.Sub(gameevent.GameBetCreatedEvent, onGameBet)
	event.Sub(gameevent.RechargeOrderCreatedEvent, onRechargeOrderCreated)
}

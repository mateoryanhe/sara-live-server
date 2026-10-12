package tracking

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/core/xrtimer"
	"xr-game-server/gameevent"
)

const trackingConsumeInterval = time.Second

var (
	trackingQueue                     = gqueue.NewTQueue[*gameevent.ShowcaseLiveRoomJoinTrackingEventData]()
	liveFirstFrameRenderTrackingQueue = gqueue.NewTQueue[*gameevent.LiveFirstFrameRenderedTrackingEventData]()
	showcaseLeaveTrackingQueue        = gqueue.NewTQueue[*gameevent.ShowcaseLiveRoomLeaveTrackingEventData]()
	call1v1InitiateTrackingQueue      = gqueue.NewTQueue[*gameevent.Call1v1InitiateTrackingEventData]()
	call1v1RoomTrackingQueue              = gqueue.NewTQueue[*gameevent.Call1v1RoomTrackingEventData]()
	call1v1ConnectSuccessTrackingQueue  = gqueue.NewTQueue[*gameevent.Call1v1ConnectSuccessTrackingEventData]()
	miniGameRoundStartTrackingQueue     = gqueue.NewTQueue[*gameevent.MiniGameRoundStartTrackingEventData]()
	miniGameRoundResultTrackingQueue    = gqueue.NewTQueue[*gameevent.MiniGameRoundResultTrackingEventData]()
	miniGameExposureTrackingQueue       = gqueue.NewTQueue[*gameevent.MiniGameExposureTrackingEventData]()
	gameLiveRoomJoinTrackingQueue       = gqueue.NewTQueue[*gameevent.GameLiveRoomJoinTrackingEventData]()
	miniGameWebViewLoadTrackingQueue    = gqueue.NewTQueue[*gameevent.MiniGameWebViewLoadTrackingEventData]()
	shortVideoPlayStartTrackingQueue    = gqueue.NewTQueue[*gameevent.ShortVideoPlayStartTrackingEventData]()
	shortVideoPaywallClickTrackingQueue = gqueue.NewTQueue[*gameevent.ShortVideoPaywallClickTrackingEventData]()
	shortVideoPaySuccessTrackingQueue   = gqueue.NewTQueue[*gameevent.ShortVideoPaySuccessTrackingEventData]()
	shortVideoWatchEndTrackingQueue     = gqueue.NewTQueue[*gameevent.ShortVideoWatchEndTrackingEventData]()
	voiceRoomEnterTrackingQueue         = gqueue.NewTQueue[*gameevent.VoiceRoomEnterTrackingEventData]()
	voiceRoomMicRequestTrackingQueue    = gqueue.NewTQueue[*gameevent.VoiceRoomMicRequestTrackingEventData]()
	voiceRoomGiftClickTrackingQueue     = gqueue.NewTQueue[*gameevent.VoiceRoomGiftClickTrackingEventData]()
	hotLiveRoomGiftClickTrackingQueue   = gqueue.NewTQueue[*gameevent.HotLiveRoomGiftClickTrackingEventData]()
	chargePanelEnterTrackingQueue       = gqueue.NewTQueue[*gameevent.ChargePanelEnterTrackingEventData]()
	chargeClickTrackingQueue            = gqueue.NewTQueue[*gameevent.ChargeClickTrackingEventData]()
	chargeResultServerTrackingQueue     = gqueue.NewTQueue[*gameevent.ChargeResultServerTrackingEventData]()
)

func initTrackingQueue() {
	event.Sub(event.HotStart, onTrackingHotStart)
	event.Sub(event.PrepareRestart, onTrackingPrepareRestart)
	event.Sub(gameevent.ShowcaseLiveRoomJoinTrackingEvent, onShowcaseLiveRoomJoinTrackingEvent)
	event.Sub(gameevent.LiveFirstFrameRenderedTrackingEvent, onLiveFirstFrameRenderedTrackingEvent)
	event.Sub(gameevent.ShowcaseLiveRoomLeaveTrackingEvent, onShowcaseLiveRoomLeaveTrackingEvent)
	event.Sub(gameevent.Call1v1InitiateTrackingEvent, onCall1v1InitiateTrackingEvent)
	event.Sub(gameevent.Call1v1RoomTrackingEvent, onCall1v1RoomTrackingEvent)
	event.Sub(gameevent.Call1v1ConnectSuccessTrackingEvent, onCall1v1ConnectSuccessTrackingEvent)
	event.Sub(gameevent.MiniGameRoundStartTrackingEvent, onMiniGameRoundStartTrackingEvent)
	event.Sub(gameevent.MiniGameRoundResultTrackingEvent, onMiniGameRoundResultTrackingEvent)
	event.Sub(gameevent.MiniGameExposureTrackingEvent, onMiniGameExposureTrackingEvent)
	event.Sub(gameevent.GameLiveRoomJoinTrackingEvent, onGameLiveRoomJoinTrackingEvent)
	event.Sub(gameevent.MiniGameWebViewLoadTrackingEvent, onMiniGameWebViewLoadTrackingEvent)
	event.Sub(gameevent.ShortVideoPlayStartTrackingEvent, onShortVideoPlayStartTrackingEvent)
	event.Sub(gameevent.ShortVideoPaywallClickTrackingEvent, onShortVideoPaywallClickTrackingEvent)
	event.Sub(gameevent.ShortVideoPaySuccessTrackingEvent, onShortVideoPaySuccessTrackingEvent)
	event.Sub(gameevent.ShortVideoWatchEndTrackingEvent, onShortVideoWatchEndTrackingEvent)
	event.Sub(gameevent.VoiceRoomEnterTrackingEvent, onVoiceRoomEnterTrackingEvent)
	event.Sub(gameevent.VoiceRoomMicRequestTrackingEvent, onVoiceRoomMicRequestTrackingEvent)
	event.Sub(gameevent.VoiceRoomGiftClickTrackingEvent, onVoiceRoomGiftClickTrackingEvent)
	event.Sub(gameevent.HotLiveRoomGiftClickTrackingEvent, onHotLiveRoomGiftClickTrackingEvent)
	event.Sub(gameevent.ChargePanelEnterTrackingEvent, onChargePanelEnterTrackingEvent)
	event.Sub(gameevent.ChargeClickTrackingEvent, onChargeClickTrackingEvent)
	event.Sub(gameevent.ChargeResultServerTrackingEvent, onChargeResultServerTrackingEvent)
}

func onShowcaseLiveRoomJoinTrackingEvent(data any) {
	ev, ok := data.(*gameevent.ShowcaseLiveRoomJoinTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.RoomId == 0 {
		return
	}
	trackingQueue.Push(ev)
}

func onLiveFirstFrameRenderedTrackingEvent(data any) {
	ev, ok := data.(*gameevent.LiveFirstFrameRenderedTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.RoomId == 0 {
		return
	}
	liveFirstFrameRenderTrackingQueue.Push(ev)
}

func onShowcaseLiveRoomLeaveTrackingEvent(data any) {
	ev, ok := data.(*gameevent.ShowcaseLiveRoomLeaveTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.RoomId == 0 {
		return
	}
	showcaseLeaveTrackingQueue.Push(ev)
}

func onCall1v1InitiateTrackingEvent(data any) {
	ev, ok := data.(*gameevent.Call1v1InitiateTrackingEventData)
	if !ok || ev == nil || ev.CallerId == 0 || ev.AnchorId == 0 {
		return
	}
	call1v1InitiateTrackingQueue.Push(ev)
}

func onCall1v1RoomTrackingEvent(data any) {
	ev, ok := data.(*gameevent.Call1v1RoomTrackingEventData)
	if !ok || ev == nil || ev.CallerId == 0 || ev.AnchorId == 0 {
		return
	}
	call1v1RoomTrackingQueue.Push(ev)
}

func onCall1v1ConnectSuccessTrackingEvent(data any) {
	ev, ok := data.(*gameevent.Call1v1ConnectSuccessTrackingEventData)
	if !ok || ev == nil || ev.CallerId == 0 || ev.ReceiverId == 0 {
		return
	}
	call1v1ConnectSuccessTrackingQueue.Push(ev)
}

func onMiniGameRoundStartTrackingEvent(data any) {
	ev, ok := data.(*gameevent.MiniGameRoundStartTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	miniGameRoundStartTrackingQueue.Push(ev)
}

func onMiniGameRoundResultTrackingEvent(data any) {
	ev, ok := data.(*gameevent.MiniGameRoundResultTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	miniGameRoundResultTrackingQueue.Push(ev)
}

func onMiniGameExposureTrackingEvent(data any) {
	ev, ok := data.(*gameevent.MiniGameExposureTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	miniGameExposureTrackingQueue.Push(ev)
}

func onGameLiveRoomJoinTrackingEvent(data any) {
	ev, ok := data.(*gameevent.GameLiveRoomJoinTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.RoomId == 0 {
		return
	}
	gameLiveRoomJoinTrackingQueue.Push(ev)
}

func onMiniGameWebViewLoadTrackingEvent(data any) {
	ev, ok := data.(*gameevent.MiniGameWebViewLoadTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	miniGameWebViewLoadTrackingQueue.Push(ev)
}

func onShortVideoPlayStartTrackingEvent(data any) {
	ev, ok := data.(*gameevent.ShortVideoPlayStartTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	shortVideoPlayStartTrackingQueue.Push(ev)
}

func onShortVideoPaywallClickTrackingEvent(data any) {
	ev, ok := data.(*gameevent.ShortVideoPaywallClickTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	shortVideoPaywallClickTrackingQueue.Push(ev)
}

func onShortVideoPaySuccessTrackingEvent(data any) {
	ev, ok := data.(*gameevent.ShortVideoPaySuccessTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	shortVideoPaySuccessTrackingQueue.Push(ev)
}

func onShortVideoWatchEndTrackingEvent(data any) {
	ev, ok := data.(*gameevent.ShortVideoWatchEndTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	shortVideoWatchEndTrackingQueue.Push(ev)
}

func onVoiceRoomEnterTrackingEvent(data any) {
	ev, ok := data.(*gameevent.VoiceRoomEnterTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 || ev.RoomId == 0 {
		return
	}
	voiceRoomEnterTrackingQueue.Push(ev)
}

func onVoiceRoomMicRequestTrackingEvent(data any) {
	ev, ok := data.(*gameevent.VoiceRoomMicRequestTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	voiceRoomMicRequestTrackingQueue.Push(ev)
}

func onVoiceRoomGiftClickTrackingEvent(data any) {
	ev, ok := data.(*gameevent.VoiceRoomGiftClickTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	voiceRoomGiftClickTrackingQueue.Push(ev)
}

func onHotLiveRoomGiftClickTrackingEvent(data any) {
	ev, ok := data.(*gameevent.HotLiveRoomGiftClickTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	hotLiveRoomGiftClickTrackingQueue.Push(ev)
}

func onChargePanelEnterTrackingEvent(data any) {
	ev, ok := data.(*gameevent.ChargePanelEnterTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	chargePanelEnterTrackingQueue.Push(ev)
}

func onChargeClickTrackingEvent(data any) {
	ev, ok := data.(*gameevent.ChargeClickTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	chargeClickTrackingQueue.Push(ev)
}

func onChargeResultServerTrackingEvent(data any) {
	ev, ok := data.(*gameevent.ChargeResultServerTrackingEventData)
	if !ok || ev == nil || ev.UserId == 0 {
		return
	}
	chargeResultServerTrackingQueue.Push(ev)
}

func onTrackingHotStart(_ any) {
	xrtimer.AddSingleton(gctx.New(), trackingConsumeInterval, func(ctx context.Context) {
		drainTrackingQueue()
	})
}

func onTrackingPrepareRestart(_ any) {
	drainTrackingQueue()
}

func trackingQueueIdle() bool {
	return trackingQueue.Len() == 0 &&
		liveFirstFrameRenderTrackingQueue.Len() == 0 &&
		showcaseLeaveTrackingQueue.Len() == 0 &&
		call1v1InitiateTrackingQueue.Len() == 0 &&
		call1v1RoomTrackingQueue.Len() == 0 &&
		call1v1ConnectSuccessTrackingQueue.Len() == 0 &&
		miniGameRoundStartTrackingQueue.Len() == 0 &&
		miniGameRoundResultTrackingQueue.Len() == 0 &&
		miniGameExposureTrackingQueue.Len() == 0 &&
		gameLiveRoomJoinTrackingQueue.Len() == 0 &&
		miniGameWebViewLoadTrackingQueue.Len() == 0 &&
		shortVideoPlayStartTrackingQueue.Len() == 0 &&
		shortVideoPaywallClickTrackingQueue.Len() == 0 &&
		shortVideoPaySuccessTrackingQueue.Len() == 0 &&
		shortVideoWatchEndTrackingQueue.Len() == 0 &&
		voiceRoomEnterTrackingQueue.Len() == 0 &&
		voiceRoomMicRequestTrackingQueue.Len() == 0 &&
		voiceRoomGiftClickTrackingQueue.Len() == 0 &&
		hotLiveRoomGiftClickTrackingQueue.Len() == 0 &&
		chargePanelEnterTrackingQueue.Len() == 0 &&
		chargeClickTrackingQueue.Len() == 0 &&
		chargeResultServerTrackingQueue.Len() == 0
}

func drainTrackingTQueue[T any](q *gqueue.TQueue[T], handle func(T)) {
	for {
		select {
		case ev := <-q.C:
			handle(ev)
		default:
			return
		}
	}
}

func drainTrackingQueue() {
	drainTrackingTQueue(trackingQueue, func(ev *gameevent.ShowcaseLiveRoomJoinTrackingEventData) {
		recordShowcaseLiveRoomJoin(ev.UserId, ev.RoomId, ev.StatAt)
	})
	drainTrackingTQueue(liveFirstFrameRenderTrackingQueue, func(ev *gameevent.LiveFirstFrameRenderedTrackingEventData) {
		recordLiveFirstFrameRendered(ev.UserId, ev.RoomId, ev.StatAt)
		recordGameLiveWatchStream(ev.UserId, ev.RoomId, ev.StatAt)
	})
	drainTrackingTQueue(showcaseLeaveTrackingQueue, func(ev *gameevent.ShowcaseLiveRoomLeaveTrackingEventData) {
		recordShowcaseLiveRoomLeave(ev.UserId, ev.RoomId, ev.StatAt)
	})
	drainTrackingTQueue(call1v1InitiateTrackingQueue, func(ev *gameevent.Call1v1InitiateTrackingEventData) {
		recordCall1v1Initiate(ev.CallerId, ev.AnchorId, ev.StatAt)
	})
	drainTrackingTQueue(call1v1RoomTrackingQueue, func(ev *gameevent.Call1v1RoomTrackingEventData) {
		recordCall1v1Room(ev.CallerId, ev.AnchorId, ev.StatAt)
	})
	drainTrackingTQueue(call1v1ConnectSuccessTrackingQueue, func(ev *gameevent.Call1v1ConnectSuccessTrackingEventData) {
		recordCall1v1ConnectSuccess(ev.Source, ev.CallType, ev.CallerId, ev.ReceiverId, ev.StatAt)
	})
	drainTrackingTQueue(miniGameRoundStartTrackingQueue, func(ev *gameevent.MiniGameRoundStartTrackingEventData) {
		recordMiniGameRoundStart(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(miniGameRoundResultTrackingQueue, func(ev *gameevent.MiniGameRoundResultTrackingEventData) {
		recordMiniGameRoundResult(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(miniGameExposureTrackingQueue, func(ev *gameevent.MiniGameExposureTrackingEventData) {
		recordMiniGameExposure(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(gameLiveRoomJoinTrackingQueue, func(ev *gameevent.GameLiveRoomJoinTrackingEventData) {
		recordGameLiveRoomJoin(ev.UserId, ev.RoomId, ev.StatAt)
	})
	drainTrackingTQueue(miniGameWebViewLoadTrackingQueue, func(ev *gameevent.MiniGameWebViewLoadTrackingEventData) {
		recordMiniGameWebViewLoad(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(shortVideoPlayStartTrackingQueue, func(ev *gameevent.ShortVideoPlayStartTrackingEventData) {
		recordShortVideoPlayStart(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(shortVideoPaywallClickTrackingQueue, func(ev *gameevent.ShortVideoPaywallClickTrackingEventData) {
		recordShortVideoPaywallClick(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(shortVideoPaySuccessTrackingQueue, func(ev *gameevent.ShortVideoPaySuccessTrackingEventData) {
		recordShortVideoPaySuccess(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(shortVideoWatchEndTrackingQueue, func(ev *gameevent.ShortVideoWatchEndTrackingEventData) {
		recordShortVideoWatchEnd(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(voiceRoomEnterTrackingQueue, func(ev *gameevent.VoiceRoomEnterTrackingEventData) {
		recordVoiceRoomEnter(ev.UserId, ev.RoomId, ev.StatAt)
	})
	drainTrackingTQueue(voiceRoomMicRequestTrackingQueue, func(ev *gameevent.VoiceRoomMicRequestTrackingEventData) {
		recordVoiceRoomMicRequest(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(voiceRoomGiftClickTrackingQueue, func(ev *gameevent.VoiceRoomGiftClickTrackingEventData) {
		recordVoiceRoomGiftClick(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(hotLiveRoomGiftClickTrackingQueue, func(ev *gameevent.HotLiveRoomGiftClickTrackingEventData) {
		recordHotLiveRoomGiftClick(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(chargePanelEnterTrackingQueue, func(ev *gameevent.ChargePanelEnterTrackingEventData) {
		recordChargePanelEnter(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(chargeClickTrackingQueue, func(ev *gameevent.ChargeClickTrackingEventData) {
		recordChargeClick(ev.UserId, ev.StatAt)
	})
	drainTrackingTQueue(chargeResultServerTrackingQueue, func(ev *gameevent.ChargeResultServerTrackingEventData) {
		recordChargeResultServer(ev.UserId, ev.StatAt)
	})
}

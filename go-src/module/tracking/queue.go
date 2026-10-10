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
		miniGameWebViewLoadTrackingQueue.Len() == 0
}

func drainTrackingQueue() {
	for {
		select {
		case ev := <-trackingQueue.C:
			if ev != nil {
				recordShowcaseLiveRoomJoin(ev.UserId, ev.RoomId, ev.StatAt)
			}
		default:
			break
		}
	}
	for {
		select {
		case ev := <-liveFirstFrameRenderTrackingQueue.C:
			if ev != nil {
				recordLiveFirstFrameRendered(ev.UserId, ev.RoomId, ev.StatAt)
			}
		default:
			break
		}
	}
	for {
		select {
		case ev := <-showcaseLeaveTrackingQueue.C:
			if ev != nil {
				recordShowcaseLiveRoomLeave(ev.UserId, ev.RoomId, ev.StatAt)
			}
		default:
			break
		}
	}
	for {
		select {
		case ev := <-call1v1InitiateTrackingQueue.C:
			if ev != nil {
				recordCall1v1Initiate(ev.CallerId, ev.AnchorId, ev.StatAt)
			}
		default:
			break
		}
	}
	for {
		select {
		case ev := <-call1v1RoomTrackingQueue.C:
			if ev != nil {
				recordCall1v1Room(ev.CallerId, ev.AnchorId, ev.StatAt)
			}
		default:
			break
		}
	}
	for {
		select {
		case ev := <-call1v1ConnectSuccessTrackingQueue.C:
			if ev != nil {
				recordCall1v1ConnectSuccess(ev.Source, ev.CallType, ev.CallerId, ev.ReceiverId, ev.StatAt)
			}
		default:
			break
		}
	}
	for {
		select {
		case ev := <-miniGameRoundStartTrackingQueue.C:
			if ev != nil {
				recordMiniGameRoundStart(ev.UserId, ev.StatAt)
			}
		default:
			break
		}
	}
	for {
		select {
		case ev := <-miniGameRoundResultTrackingQueue.C:
			if ev != nil {
				recordMiniGameRoundResult(ev.UserId, ev.StatAt)
			}
		default:
			break
		}
	}
	for {
		select {
		case ev := <-miniGameExposureTrackingQueue.C:
			if ev != nil {
				recordMiniGameExposure(ev.UserId, ev.StatAt)
			}
		default:
			break
		}
	}
	for {
		select {
		case ev := <-gameLiveRoomJoinTrackingQueue.C:
			if ev != nil {
				recordGameLiveRoomJoin(ev.UserId, ev.RoomId, ev.StatAt)
			}
		default:
			break
		}
	}
	for {
		select {
		case ev := <-miniGameWebViewLoadTrackingQueue.C:
			if ev != nil {
				recordMiniGameWebViewLoad(ev.UserId, ev.StatAt)
			}
		default:
			return
		}
	}
}

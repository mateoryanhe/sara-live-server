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
)

func initTrackingQueue() {
	event.Sub(event.HotStart, onTrackingHotStart)
	event.Sub(event.PrepareRestart, onTrackingPrepareRestart)
	event.Sub(gameevent.ShowcaseLiveRoomJoinTrackingEvent, onShowcaseLiveRoomJoinTrackingEvent)
	event.Sub(gameevent.LiveFirstFrameRenderedTrackingEvent, onLiveFirstFrameRenderedTrackingEvent)
	event.Sub(gameevent.ShowcaseLiveRoomLeaveTrackingEvent, onShowcaseLiveRoomLeaveTrackingEvent)
	event.Sub(gameevent.Call1v1InitiateTrackingEvent, onCall1v1InitiateTrackingEvent)
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
		call1v1InitiateTrackingQueue.Len() == 0
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
			return
		}
	}
}

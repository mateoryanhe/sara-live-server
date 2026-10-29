package stat

import (
	"time"

	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

func initPromotionChannelH5DownloadEvent() {
	event.Sub(gameevent.H5DownloadTrackingEvent, onH5DownloadForPromotionChannel)
}

func onH5DownloadForPromotionChannel(data any) {
	enqueue(&statJob{Kind: jobPromotionChannelH5Download, Payload: data})
}

func consumePromotionChannelH5DownloadJob(job *statJob) {
	ev, ok := job.Payload.(*gameevent.H5DownloadTrackingEventData)
	if !ok || ev == nil || ev.ChannelCode <= 0 {
		return
	}
	at := ev.StatAt
	if at.IsZero() {
		at = time.Now()
	}
	recordDailyPromotionChannelH5DownloadApk(at, ev.ChannelCode)
}

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

// EnqueuePromotionChannelH5DownloadApk lookupApkDownload 等命中推广渠道包时计入 H5 下载 APK 按钮点击.
func EnqueuePromotionChannelH5DownloadApk(channelCode int, at time.Time) {
	if channelCode <= 0 {
		return
	}
	enqueue(&statJob{Kind: jobPromotionChannelH5Download, Payload: gameevent.NewH5DownloadTrackingEventData(at, channelCode)})
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

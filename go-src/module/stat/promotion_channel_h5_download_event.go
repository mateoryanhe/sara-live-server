package stat

import (
	"time"

	"xr-game-server/gameevent"
)

func initPromotionChannelH5DownloadEvent() {}

// EnqueuePromotionChannelH5DownloadApk lookupApkDownload 命中推广渠道包时计入 H5推广日统计 h5_download_apk_count.
// 不与 reportH5Download 重复：埋点走 H5DownloadTrackingEvent 只计 tracking 表，渠道日统计仅本入口 +1.
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

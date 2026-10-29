package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const H5DownloadTrackingEvent event.Type = "H5DownloadTrackingEvent"
const OfficialSiteDownloadTrackingEvent event.Type = "OfficialSiteDownloadTrackingEvent"
const ShareShortLinkClickTrackingEvent event.Type = "ShareShortLinkClickTrackingEvent"

type H5DownloadTrackingEventData struct {
	StatAt      time.Time
	ChannelCode int // 推广渠道编码(>0 时计入渠道日统计)
}

type OfficialSiteDownloadTrackingEventData struct {
	StatAt time.Time
}

type ShareShortLinkClickTrackingEventData struct {
	StatAt time.Time
}

func NewH5DownloadTrackingEventData(statAt time.Time, channelCode int) *H5DownloadTrackingEventData {
	return &H5DownloadTrackingEventData{StatAt: statAt, ChannelCode: channelCode}
}

func NewOfficialSiteDownloadTrackingEventData(statAt time.Time) *OfficialSiteDownloadTrackingEventData {
	return &OfficialSiteDownloadTrackingEventData{StatAt: statAt}
}

func NewShareShortLinkClickTrackingEventData(statAt time.Time) *ShareShortLinkClickTrackingEventData {
	return &ShareShortLinkClickTrackingEventData{StatAt: statAt}
}

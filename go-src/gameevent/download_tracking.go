package gameevent

import (
	"time"

	"xr-game-server/core/event"
)

const H5DownloadTrackingEvent event.Type = "H5DownloadTrackingEvent"
const OfficialSiteDownloadTrackingEvent event.Type = "OfficialSiteDownloadTrackingEvent"

type H5DownloadTrackingEventData struct {
	StatAt time.Time
}

type OfficialSiteDownloadTrackingEventData struct {
	StatAt time.Time
}

func NewH5DownloadTrackingEventData(statAt time.Time) *H5DownloadTrackingEventData {
	return &H5DownloadTrackingEventData{StatAt: statAt}
}

func NewOfficialSiteDownloadTrackingEventData(statAt time.Time) *OfficialSiteDownloadTrackingEventData {
	return &OfficialSiteDownloadTrackingEventData{StatAt: statAt}
}

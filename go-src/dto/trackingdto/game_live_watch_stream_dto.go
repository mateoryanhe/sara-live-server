package trackingdto

import "github.com/gogf/gf/v2/frame/g"

type CMSGameLiveWatchStreamTrendReq struct {
	g.Meta `path:"/getGameLiveWatchStreamTrend" method:"post" summary:"CMS游戏直播流首帧加载次数趋势" tags:"埋点"`
}

type CMSGameLiveWatchStreamTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

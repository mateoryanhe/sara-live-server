package trackingdto

import "github.com/gogf/gf/v2/frame/g"

type CMSShortVideoPlayStartTrendReq struct {
	g.Meta `path:"/getShortVideoPlayStartTrend" method:"post" summary:"CMS短视频起播次数趋势" tags:"埋点"`
}

type CMSShortVideoPlayStartTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

type CMSShortVideoPaywallClickTrendReq struct {
	g.Meta `path:"/getShortVideoPaywallClickTrend" method:"post" summary:"CMS短视频付费解锁点击次数趋势" tags:"埋点"`
}

type CMSShortVideoPaywallClickTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

type CMSShortVideoPaySuccessTrendReq struct {
	g.Meta `path:"/getShortVideoPaySuccessTrend" method:"post" summary:"CMS短视频付费成功次数趋势" tags:"埋点"`
}

type CMSShortVideoPaySuccessTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

type CMSShortVideoWatchEndTrendReq struct {
	g.Meta `path:"/getShortVideoWatchEndTrend" method:"post" summary:"CMS短视频播放结束次数趋势" tags:"埋点"`
}

type CMSShortVideoWatchEndTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

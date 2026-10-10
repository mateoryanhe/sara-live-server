package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// ReportMiniGameExposureReq App 上报半屏游戏窗口曝光(暂不要求 body 字段)
type ReportMiniGameExposureReq struct {
	g.Meta `path:"/reportMiniGameExposure" method:"post" summary:"上报半屏游戏窗口曝光" tags:"直播间"`
}

type ReportMiniGameExposureRes struct {
	Success bool `json:"success"`
}

// CMSMiniGameExposureTrendReq CMS 半屏游戏曝光次数趋势
type CMSMiniGameExposureTrendReq struct {
	g.Meta `path:"/getMiniGameExposureTrend" method:"post" summary:"CMS半屏游戏曝光次数趋势" tags:"埋点"`
}

type CMSMiniGameExposureTrendRes struct {
	TodayCount uint64                        `json:"todayCount" dc:"今日次数"`
	WeekCount  uint64                        `json:"weekCount" dc:"本周次数"`
	MonthCount uint64                        `json:"monthCount" dc:"本月次数"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日(最近30天)"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly" dc:"按周(最近12周)"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly" dc:"按月(最近12月)"`
}

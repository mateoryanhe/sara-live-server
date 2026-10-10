package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// ReportMiniGameWebViewLoadSuccessReq App 上报游戏 WebView 成功加载(暂不要求 body 字段)
type ReportMiniGameWebViewLoadSuccessReq struct {
	g.Meta `path:"/reportMiniGameWebViewLoadSuccess" method:"post" summary:"上报游戏WebView成功加载" tags:"直播间"`
}

type ReportMiniGameWebViewLoadSuccessRes struct {
	Success bool `json:"success"`
}

// CMSMiniGameWebViewLoadTrendReq CMS 游戏 WebView 加载成功次数趋势
type CMSMiniGameWebViewLoadTrendReq struct {
	g.Meta `path:"/getMiniGameWebViewLoadTrend" method:"post" summary:"CMS游戏WebView加载成功次数趋势" tags:"埋点"`
}

type CMSMiniGameWebViewLoadTrendRes struct {
	TodayCount uint64                        `json:"todayCount" dc:"今日次数"`
	WeekCount  uint64                        `json:"weekCount" dc:"本周次数"`
	MonthCount uint64                        `json:"monthCount" dc:"本月次数"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日(最近30天)"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly" dc:"按周(最近12周)"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly" dc:"按月(最近12月)"`
}

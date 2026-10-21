package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// ReportMiniGamePlusClickReq App 上报游戏+(暂不要求 body 字段)
type ReportMiniGamePlusClickReq struct {
	g.Meta `path:"/reportMiniGamePlusClick" method:"post" summary:"上报游戏+" tags:"直播间"`
}

type ReportMiniGamePlusClickRes struct {
	Success bool `json:"success"`
}

// CMSMiniGamePlusClickTrendReq CMS 游戏+ 点击次数趋势
type CMSMiniGamePlusClickTrendReq struct {
	g.Meta `path:"/getMiniGamePlusClickTrend" method:"post" summary:"CMS游戏+点击次数趋势" tags:"埋点"`
}

type CMSMiniGamePlusClickTrendRes struct {
	TodayCount uint64                        `json:"todayCount" dc:"今日次数"`
	WeekCount  uint64                        `json:"weekCount" dc:"本周次数"`
	MonthCount uint64                        `json:"monthCount" dc:"本月次数"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日(最近30天)"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly" dc:"按周(最近12周)"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly" dc:"按月(最近12月)"`
}

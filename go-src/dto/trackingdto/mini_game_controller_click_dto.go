package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// ReportMiniGameControllerClickReq App 上报点击游戏手柄(暂不要求 body 字段)
type ReportMiniGameControllerClickReq struct {
	g.Meta `path:"/reportMiniGameControllerClick" method:"post" summary:"上报点击游戏手柄" tags:"直播间"`
}

type ReportMiniGameControllerClickRes struct {
	Success bool `json:"success"`
}

// CMSMiniGameControllerClickTrendReq CMS 点击游戏手柄次数趋势
type CMSMiniGameControllerClickTrendReq struct {
	g.Meta `path:"/getMiniGameControllerClickTrend" method:"post" summary:"CMS点击游戏手柄次数趋势" tags:"埋点"`
}

type CMSMiniGameControllerClickTrendRes struct {
	TodayCount uint64                        `json:"todayCount" dc:"今日次数"`
	WeekCount  uint64                        `json:"weekCount" dc:"本周次数"`
	MonthCount uint64                        `json:"monthCount" dc:"本月次数"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日(最近30天)"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly" dc:"按周(最近12周)"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly" dc:"按月(最近12月)"`
}

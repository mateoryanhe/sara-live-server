package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// CMSMiniGameRoundResultTrendReq CMS 单局游戏结算次数趋势
type CMSMiniGameRoundResultTrendReq struct {
	g.Meta `path:"/getMiniGameRoundResultTrend" method:"post" summary:"CMS单局游戏结算次数趋势" tags:"埋点"`
}

type CMSMiniGameRoundResultTrendRes struct {
	TodayCount uint64                        `json:"todayCount" dc:"今日次数"`
	WeekCount  uint64                        `json:"weekCount" dc:"本周次数"`
	MonthCount uint64                        `json:"monthCount" dc:"本月次数"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日(最近30天)"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly" dc:"按周(最近12周)"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly" dc:"按月(最近12月)"`
}

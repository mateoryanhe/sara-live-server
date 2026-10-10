package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// CMSCall1v1RoomCallTrendReq CMS 1v1 房间视频通话发起次数趋势(source=3)
type CMSCall1v1RoomCallTrendReq struct {
	g.Meta `path:"/getCall1v1RoomCallTrend" method:"post" summary:"CMS 1v1房间视频通话次数趋势(source=3)" tags:"埋点"`
}

type CMSCall1v1RoomCallTrendRes struct {
	TodayCount uint64                        `json:"todayCount" dc:"今日次数"`
	WeekCount  uint64                        `json:"weekCount" dc:"本周次数"`
	MonthCount uint64                        `json:"monthCount" dc:"本月次数"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日(最近30天)"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly" dc:"按周(最近12周)"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly" dc:"按月(最近12月)"`
}

package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// CMSTrackingEventTrendPoint 埋点事件趋势点(次数)
type CMSTrackingEventTrendPoint struct {
	Time  string `json:"time" dc:"时间桶(日/周/月标识)"`
	Count uint64 `json:"count" dc:"次数"`
}

// CMSHotLiveRoomJoinTrendReq CMS 进入 Hot 直播间次数趋势
type CMSHotLiveRoomJoinTrendReq struct {
	g.Meta `path:"/getHotLiveRoomJoinTrend" method:"post" summary:"CMS进入秀场直播间次数趋势" tags:"埋点"`
}

// CMSHotLiveRoomJoinTrendRes 进入 Hot 直播间次数(日/周/月)
type CMSHotLiveRoomJoinTrendRes struct {
	TodayCount uint64                        `json:"todayCount" dc:"今日次数"`
	WeekCount  uint64                        `json:"weekCount" dc:"本周次数"`
	MonthCount uint64                        `json:"monthCount" dc:"本月次数"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日(最近30天)"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly" dc:"按周(最近12周)"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly" dc:"按月(最近12月)"`
}

package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// CMSGameLiveRoomJoinTrendReq CMS 进入游戏直播类直播间次数趋势
type CMSGameLiveRoomJoinTrendReq struct {
	g.Meta `path:"/getGameLiveRoomJoinTrend" method:"post" summary:"CMS进入游戏直播类直播间次数趋势" tags:"埋点"`
}

// CMSGameLiveRoomJoinTrendRes 进入游戏直播类直播间次数(日/周/月)
type CMSGameLiveRoomJoinTrendRes struct {
	TodayCount uint64                        `json:"todayCount" dc:"今日次数"`
	WeekCount  uint64                        `json:"weekCount" dc:"本周次数"`
	MonthCount uint64                        `json:"monthCount" dc:"本月次数"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日(最近30天)"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly" dc:"按周(最近12周)"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly" dc:"按月(最近12月)"`
}

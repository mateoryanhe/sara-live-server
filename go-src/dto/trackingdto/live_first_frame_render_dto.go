package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// ReportLiveFirstFrameRenderedReq App 上报直播首帧画面渲染完成
type ReportLiveFirstFrameRenderedReq struct {
	g.Meta `path:"/reportLiveFirstFrameRendered" method:"post" summary:"上报直播首帧画面渲染完成" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
}

type ReportLiveFirstFrameRenderedRes struct {
	Success bool `json:"success"`
}

// CMSLiveFirstFrameRenderTrendReq CMS 直播首帧渲染完成次数趋势
type CMSLiveFirstFrameRenderTrendReq struct {
	g.Meta `path:"/getLiveFirstFrameRenderTrend" method:"post" summary:"CMS直播首帧画面渲染完成次数趋势" tags:"埋点"`
}

// CMSLiveFirstFrameRenderTrendRes 首帧渲染完成次数(日/周/月)
type CMSLiveFirstFrameRenderTrendRes struct {
	TodayCount uint64                        `json:"todayCount" dc:"今日次数"`
	WeekCount  uint64                        `json:"weekCount" dc:"本周次数"`
	MonthCount uint64                        `json:"monthCount" dc:"本月次数"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日(最近30天)"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly" dc:"按周(最近12周)"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly" dc:"按月(最近12月)"`
}

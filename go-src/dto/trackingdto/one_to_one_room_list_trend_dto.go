package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// CMSOneToOneRoomListTrendReq CMS 1v1 房间列表查询次数趋势
type CMSOneToOneRoomListTrendReq struct {
	g.Meta `path:"/getOneToOneRoomListTrend" method:"post" summary:"CMS 1v1房间列表查询次数趋势" tags:"埋点"`
}

type CMSOneToOneRoomListTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// ReportVoiceRoomGiftClickReq App 上报语聊房礼物面板点击(暂不要求 body)
type ReportVoiceRoomGiftClickReq struct {
	g.Meta `path:"/reportVoiceRoomGiftClick" method:"post" summary:"上报语聊房礼物面板点击" tags:"直播间"`
}

type ReportVoiceRoomGiftClickRes struct {
	Success bool `json:"success"`
}

type CMSVoiceRoomEnterTrendReq struct {
	g.Meta `path:"/getVoiceRoomEnterTrend" method:"post" summary:"CMS进入语聊房次数趋势" tags:"埋点"`
}

type CMSVoiceRoomEnterTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

type CMSVoiceRoomMicRequestTrendReq struct {
	g.Meta `path:"/getVoiceRoomMicRequestTrend" method:"post" summary:"CMS语聊房申请上麦次数趋势" tags:"埋点"`
}

type CMSVoiceRoomMicRequestTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

type CMSVoiceRoomGiftClickTrendReq struct {
	g.Meta `path:"/getVoiceRoomGiftClickTrend" method:"post" summary:"CMS语聊房礼物面板点击次数趋势" tags:"埋点"`
}

type CMSVoiceRoomGiftClickTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

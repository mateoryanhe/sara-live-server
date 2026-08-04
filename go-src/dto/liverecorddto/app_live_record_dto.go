package liverecorddto

import "github.com/gogf/gf/v2/frame/g"

// AppLiveRecordListReq App端分页查询主播直播数据
type AppLiveRecordListReq struct {
	g.Meta   `path:"/appLiveRecordList" method:"post" summary:"App查询主播直播数据" tags:"直播记录"`
	Page     int `json:"page"     v:"min:1#页码从1开始" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" v:"max:100#单页最多100条" dc:"每页数量(默认20,最大100)"`
}

// AppLiveRecordItem App端直播记录条目
type AppLiveRecordItem struct {
	Id                           string  `json:"id"                dc:"直播记录ID"`
	StartTime                    int64   `json:"startTime"         dc:"开播时间(秒)"`
	EndTime                      int64   `json:"endTime"           dc:"下播时间(秒,0表示未结束)"`
	TotalAudience                uint64  `json:"totalAudience"     dc:"累计观众人数(去重)"`
	TotalLiveDuration            float64 `json:"totalLiveDuration" dc:"累计直播时长(秒)"`
	TotalIncome                  float64 `json:"totalIncome"       dc:"总收益(钻石)"`
	TotalGiftIncome              float64 `json:"totalGiftIncome"          dc:"礼物收入(钻石)"`
	TotalPaidDanmakuIncome       float64 `json:"totalPaidDanmakuIncome"   dc:"付费弹幕收入(钻石)"`
	TotalPrivateRoomIncome       float64 `json:"totalPrivateRoomIncome"       dc:"私密直播间收入(钻石)"`
	TotalPrivateRoomTicketIncome float64 `json:"totalPrivateRoomTicketIncome" dc:"私密直播间门票收入(钻石)"`
	TotalPrivateRoomWatchIncome  float64 `json:"totalPrivateRoomWatchIncome"  dc:"私密房观看收入(钻石)"`
	TotalVideoCallIncome         float64 `json:"totalVideoCallIncome"         dc:"直播间视频通话收入(钻石)"`
	TotalVideoCallTicketIncome   float64 `json:"totalVideoCallTicketIncome"   dc:"直播间视频通话门票收入(钻石)"`
	TotalVideoCallBillingIncome  float64 `json:"totalVideoCallBillingIncome"  dc:"直播间视频通话计费收入(钻石)"`
	TotalGameBet                 float64 `json:"totalGameBet"                 dc:"游戏下注总金额"`
	TotalGiftSender              uint64  `json:"totalGiftSender"   dc:"送礼人数(去重)"`
	TotalNewFollower             uint64  `json:"totalNewFollower"  dc:"新加粉丝数(去重)"`
}

// AppLiveRecordListRes App端直播记录分页响应
type AppLiveRecordListRes struct {
	Page     int                  `json:"page"     dc:"当前页码"`
	PageSize int                  `json:"pageSize" dc:"每页数量"`
	List     []*AppLiveRecordItem `json:"list"     dc:"直播记录列表"`
}

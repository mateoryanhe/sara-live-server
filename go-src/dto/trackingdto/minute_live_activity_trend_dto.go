package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// CMSMinuteLiveActivityTrendReq CMS 直播间分钟级活跃趋势(保留约 3 天)
type CMSMinuteLiveActivityTrendReq struct {
	g.Meta    `path:"/getMinuteLiveActivityTrend" method:"post" summary:"CMS直播间分钟活跃趋势" tags:"埋点"`
	StartTime int64 `json:"startTime" dc:"开始时间 Unix 秒(UTC)"`
	EndTime   int64 `json:"endTime" dc:"结束时间 Unix 秒(UTC)"`
	Limit                int   `json:"limit" dc:"最多返回点数,默认4320"`
	BypassDashboardCache bool  `json:"bypassDashboardCache" dc:"true 时跳过仪表盘读缓存, 重新查库"`
}

type CMSMinuteLiveActivityTrendPoint struct {
	Time         string `json:"time"`
	JoinCount    uint64 `json:"joinCount"`
	LeaveCount   uint64 `json:"leaveCount"`
	GiftCount    uint64 `json:"giftCount"`
	ChatCount    uint64 `json:"chatCount"`
	BetCount           uint64 `json:"betCount"`
	Call1v1Count       uint64 `json:"call1v1Count"`
	RechargeOrderCount uint64 `json:"rechargeOrderCount"`
}

type CMSMinuteLiveActivityTrendRes struct {
	Points         []*CMSMinuteLiveActivityTrendPoint `json:"points"`
	RangeTruncated bool                             `json:"rangeTruncated"`
	TrendMaxPoints int                              `json:"trendMaxPoints"`
}

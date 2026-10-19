package statdto

import "github.com/gogf/gf/v2/frame/g"

// CMSGameBetUserTrendPoint 游戏下注人数趋势点
type CMSGameBetUserTrendPoint struct {
	Time        string `json:"time" dc:"时间桶(日/周/月标识)"`
	Count       uint64 `json:"count" dc:"下注人数(去重)"`
	OnlineCount uint64 `json:"onlineCount" dc:"在线人数峰值(该桶内)"`
}

// CMSGameBetUserTrendReq CMS 游戏下注人数趋势
type CMSGameBetUserTrendReq struct {
	g.Meta      `path:"/getGameBetUserTrend" method:"post" summary:"CMS游戏下注人数趋势" tags:"仪表盘"`
	Granularity string `json:"granularity" dc:"折线粒度 daily|weekly|monthly(传 startDate/endDate 时必填)"`
	StartDate   string `json:"startDate" dc:"开始日期 YYYY-MM-DD(可选)"`
	EndDate     string `json:"endDate" dc:"结束日期 YYYY-MM-DD(可选, 默认今天)"`
}

// CMSGameBetUserTrendRes 游戏下注人数(日/周/月)
type CMSGameBetUserTrendRes struct {
	TodayCount          uint64                      `json:"todayCount" dc:"今日下注人数"`
	WeekCount           uint64                      `json:"weekCount" dc:"本周下注人数"`
	MonthCount          uint64                      `json:"monthCount" dc:"本月下注人数"`
	CurrentOnlineCount  uint64                      `json:"currentOnlineCount" dc:"当前游戏在线人数"`
	TodayOnlinePeak     uint64                      `json:"todayOnlinePeak" dc:"今日游戏在线峰值"`
	WeekOnlinePeak      uint64                      `json:"weekOnlinePeak" dc:"本周游戏在线峰值"`
	MonthOnlinePeak     uint64                      `json:"monthOnlinePeak" dc:"本月游戏在线峰值"`
	Daily               []*CMSGameBetUserTrendPoint `json:"daily" dc:"按日(最近30天)"`
	Weekly              []*CMSGameBetUserTrendPoint `json:"weekly" dc:"按周(最近12周)"`
	Monthly             []*CMSGameBetUserTrendPoint `json:"monthly" dc:"按月(最近12月)"`
	LastHour            []*CMSGameBetUserTrendPoint `json:"lastHour" dc:"最近1小时(分钟桶, 默认折线)"`
	RangeTruncated      bool                        `json:"rangeTruncated" dc:"时间范围超出上限, 已截断为最近 N 个点"`
	TrendMaxPoints      int                         `json:"trendMaxPoints" dc:"当前粒度允许的最大点数"`
}

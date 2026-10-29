package statdto

import "github.com/gogf/gf/v2/frame/g"

// CMSUserStatTrendReq CMS获取用户数据趋势
type CMSUserStatTrendReq struct {
	g.Meta      `path:"/getUserStatTrend" method:"post" summary:"CMS获取用户数据趋势" tags:"仪表盘"`
	Granularity string `json:"granularity" dc:"折线粒度 daily|weekly|monthly(传 startDate/endDate 时必填)"`
	StartDate            string `json:"startDate" dc:"开始日期 YYYY-MM-DD(可选)"`
	EndDate              string `json:"endDate" dc:"结束日期 YYYY-MM-DD(可选, 默认今天)"`
	BypassDashboardCache bool   `json:"bypassDashboardCache" dc:"true 时跳过仪表盘读缓存, 重新查库/组装"`
}

// CMSUserStatTrendPoint 用户数据趋势点(折线图 + 柱形图共用时间轴)
type CMSUserStatTrendPoint struct {
	Time                    string            `json:"time"                    dc:"时间标识"`
	ActiveUserCount         uint64            `json:"activeUserCount"         dc:"活跃用户数(去重登录)"`
	RegisterUserCount       uint64            `json:"registerUserCount"       dc:"新注册用户数"`
	AppDeviceRechargeAmount float64           `json:"appDeviceRechargeAmount" dc:"App设备码用户充值(USD)"`
	H5DeviceRechargeAmount  float64           `json:"h5DeviceRechargeAmount"  dc:"H5设备码用户充值(USD)"`
	AppDeviceLoginCount     uint64            `json:"appDeviceLoginCount"     dc:"App设备码用户登录人数(去重)"`
	H5DeviceLoginCount      uint64            `json:"h5DeviceLoginCount"      dc:"H5设备码用户登录人数(去重)"`
	BarMetrics              map[string]uint64 `json:"barMetrics"              dc:"柱形图指标(见 UserStatBarMetric*)"`
}

// CMSUserStatTrendRes CMS用户数据趋势(日/周/月)
type CMSUserStatTrendRes struct {
	Daily          []*CMSUserStatTrendPoint `json:"daily"   dc:"按日统计(最近30天)"`
	Weekly         []*CMSUserStatTrendPoint `json:"weekly"  dc:"按周统计(最近12周)"`
	Monthly        []*CMSUserStatTrendPoint `json:"monthly" dc:"按月统计(最近12月)"`
	RangeTruncated bool                     `json:"rangeTruncated" dc:"时间范围超出上限, 已截断"`
	TrendMaxPoints int                      `json:"trendMaxPoints" dc:"当前粒度最大点数"`
}

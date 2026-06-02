package statdto

import "github.com/gogf/gf/v2/frame/g"

// CMSUserStatTrendReq CMS获取用户数据趋势
type CMSUserStatTrendReq struct {
	g.Meta `path:"/getUserStatTrend" method:"post" summary:"CMS获取用户数据趋势" tags:"仪表盘"`
}

// CMSUserStatTrendPoint 用户数据趋势点(折线图 + 柱形图共用时间轴)
type CMSUserStatTrendPoint struct {
	Time              string            `json:"time"              dc:"时间标识"`
	ActiveUserCount   uint64            `json:"activeUserCount"   dc:"活跃用户数(去重登录)"`
	RegisterUserCount uint64            `json:"registerUserCount" dc:"新注册用户数"`
	BarMetrics        map[string]uint64 `json:"barMetrics"        dc:"柱形图指标(见 UserStatBarMetric*)"`
}

// CMSUserStatTrendRes CMS用户数据趋势(日/周/月)
type CMSUserStatTrendRes struct {
	Daily   []*CMSUserStatTrendPoint `json:"daily"   dc:"按日统计(最近30天)"`
	Weekly  []*CMSUserStatTrendPoint `json:"weekly"  dc:"按周统计(最近12周)"`
	Monthly []*CMSUserStatTrendPoint `json:"monthly" dc:"按月统计(最近12月)"`
}

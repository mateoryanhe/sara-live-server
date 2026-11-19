package trackingdto

import "github.com/gogf/gf/v2/frame/g"

// CMSEmailAccountTrendReq CMS 邮箱绑定/注册日趋势(仅日维度)
type CMSEmailAccountTrendReq struct {
	g.Meta `path:"/getEmailAccountTrend" method:"post" summary:"CMS邮箱绑定与注册日趋势" tags:"埋点"`
}

type CMSEmailAccountTrendRes struct {
	TodayVal uint64                        `json:"todayVal" dc:"今日人数"`
	Daily    []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日折线(Count=val)"`
}
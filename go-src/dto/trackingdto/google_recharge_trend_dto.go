package trackingdto

import "github.com/gogf/gf/v2/frame/g"

type CMSGoogleRechargeDailyTrendRes struct {
	TodayCount uint64                        `json:"todayCount" dc:"今日次数"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily" dc:"按日折线"`
}

type CMSGoogleRechargeListTrendReq struct {
	g.Meta `path:"/getGoogleRechargeListTrend" method:"post" summary:"CMS Google Play充值档位列表拉取日趋势" tags:"埋点"`
}

type CMSGoogleRechargeOrderTrendReq struct {
	g.Meta `path:"/getGoogleRechargeOrderTrend" method:"post" summary:"CMS Google Play充值下单日趋势" tags:"埋点"`
}

type CMSGoogleRechargeResultTrendReq struct {
	g.Meta `path:"/getGoogleRechargeResultTrend" method:"post" summary:"CMS Google Play充值到账日趋势" tags:"埋点"`
}
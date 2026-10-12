package trackingdto

import "github.com/gogf/gf/v2/frame/g"

type ReportChargePanelEnterReq struct {
	g.Meta `path:"/reportChargePanelEnter" method:"post" summary:"上报充值面板曝光" tags:"充值订单"`
}

type ReportChargePanelEnterRes struct {
	Success bool `json:"success"`
}

type CMSChargePanelEnterTrendReq struct {
	g.Meta `path:"/getChargePanelEnterTrend" method:"post" summary:"CMS充值面板曝光次数趋势" tags:"埋点"`
}

type CMSChargePanelEnterTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

type CMSChargeClickTrendReq struct {
	g.Meta `path:"/getChargeClickTrend" method:"post" summary:"CMS充值档位点击建单次数趋势" tags:"埋点"`
}

type CMSChargeClickTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

type CMSChargeResultServerTrendReq struct {
	g.Meta `path:"/getChargeResultServerTrend" method:"post" summary:"CMS充值到账成功次数趋势" tags:"埋点"`
}

type CMSChargeResultServerTrendRes struct {
	TodayCount uint64                        `json:"todayCount"`
	WeekCount  uint64                        `json:"weekCount"`
	MonthCount uint64                        `json:"monthCount"`
	Daily      []*CMSTrackingEventTrendPoint `json:"daily"`
	Weekly     []*CMSTrackingEventTrendPoint `json:"weekly"`
	Monthly    []*CMSTrackingEventTrendPoint `json:"monthly"`
}

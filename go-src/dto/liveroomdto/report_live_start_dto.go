package liveroomdto

import "github.com/gogf/gf/v2/frame/g"

// ReportLiveStartStatusReq 主播上报开播状态(无需传参,主播ID由鉴权中间件提供)
type ReportLiveStartStatusReq struct {
	g.Meta `path:"/reportLiveStartStatus" method:"post" summary:"主播上报开播状态" tags:"直播间"`
}

type ReportLiveStartStatusRes struct {
	Success bool `json:"success"`
}

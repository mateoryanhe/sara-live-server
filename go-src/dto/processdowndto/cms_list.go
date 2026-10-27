package processdowndto

import "github.com/gogf/gf/v2/frame/g"

type CMSListProcessDownLogReq struct {
	g.Meta    `path:"/listProcessDownLog" method:"post" summary:"CMS进程挂掉记录列表" tags:"运维监控"`
	Host      string `json:"host" dc:"主机名模糊匹配,可选"`
	StartTime int64  `json:"startTime" dc:"判定挂掉时间起(unix秒),可选"`
	EndTime   int64  `json:"endTime" dc:"判定挂掉时间止(unix秒),可选"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

type CMSProcessDownLogItem struct {
	ID           uint64 `json:"id"`
	Host         string `json:"host"`
	AppName      string `json:"appName"`
	MissingSince string `json:"missingSince"`
	DetectedAt   string `json:"detectedAt"`
	ColdStartOK  bool   `json:"coldStartOK"`
	NewPID       int    `json:"newPid"`
	Detail       string `json:"detail"`
	CreatedAt    string `json:"createdAt"`
}

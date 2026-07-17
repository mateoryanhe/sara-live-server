package logquerydto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type CMSGetLogPathsReq struct {
	g.Meta `path:"/getLogPaths" method:"post" summary:"CMS获取日志目录配置" tags:"日志查询"`
}

type CMSGetLogPathsRes struct {
	ServerTime       string `json:"serverTime"`
	DetailLogDir     string `json:"detailLogDir"`
	DetailLogPattern string `json:"detailLogPattern"`
	AccessLogDir     string `json:"accessLogDir"`
	AccessLogPattern string `json:"accessLogPattern"`
	ErrorLogDir      string `json:"errorLogDir"`
	ErrorLogPattern  string `json:"errorLogPattern"`
}

type CMSQueryDetailLogsReq struct {
	g.Meta `path:"/queryDetailLogs" method:"post" summary:"CMS查询详情日志" tags:"日志查询"`
	httpserver.CMSQueryReq
	StartDate string `json:"startDate" v:"required|date-format:Y-m-d#开始日期不能为空|开始日期格式应为Y-m-d"`
	EndDate   string `json:"endDate" v:"required|date-format:Y-m-d#结束日期不能为空|结束日期格式应为Y-m-d"`
	TraceId   string `json:"traceId"`
	ReqId     string `json:"reqId"`
	AuthId    string `json:"authId"`
	Url       string `json:"url"`
	Keyword   string `json:"keyword"`
}

type CMSQueryAccessLogsReq struct {
	g.Meta `path:"/queryAccessLogs" method:"post" summary:"CMS查询Access日志" tags:"日志查询"`
	httpserver.CMSQueryReq
	StartDate    string   `json:"startDate" v:"required|date-format:Y-m-d#开始日期不能为空|开始日期格式应为Y-m-d"`
	EndDate      string   `json:"endDate" v:"required|date-format:Y-m-d#结束日期不能为空|结束日期格式应为Y-m-d"`
	TraceId      string   `json:"traceId"`
	Url          string   `json:"url"`
	Ip           string   `json:"ip"`
	StatusCode   int      `json:"statusCode"`
	MinHandlerMs *float64 `json:"minHandlerMs"`
	MaxHandlerMs *float64 `json:"maxHandlerMs"`
}

type CMSQueryErrorLogsReq struct {
	g.Meta `path:"/queryErrorLogs" method:"post" summary:"CMS查询Error日志" tags:"日志查询"`
	httpserver.CMSQueryReq
	StartDate  string `json:"startDate" v:"required|date-format:Y-m-d#开始日期不能为空|开始日期格式应为Y-m-d"`
	EndDate    string `json:"endDate" v:"required|date-format:Y-m-d#结束日期不能为空|结束日期格式应为Y-m-d"`
	TraceId    string `json:"traceId"`
	Url        string `json:"url"`
	Ip         string `json:"ip"`
	StatusCode int    `json:"statusCode"`
	Keyword    string `json:"keyword"`
}

type CMSGetTraceLogsReq struct {
	g.Meta    `path:"/getTraceLogs" method:"post" summary:"CMS按TraceId查询日志详情" tags:"日志查询"`
	TraceId   string `json:"traceId" v:"required#TraceId不能为空"`
	StartDate string `json:"startDate" v:"required|date-format:Y-m-d#开始日期不能为空|开始日期格式应为Y-m-d"`
	EndDate   string `json:"endDate" v:"required|date-format:Y-m-d#结束日期不能为空|结束日期格式应为Y-m-d"`
}

type CMSGetTraceLogsRes struct {
	TraceId    string `json:"traceId"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
	DetailLogs any    `json:"detailLogs"`
	AccessLogs any    `json:"accessLogs"`
	ErrorLogs  any    `json:"errorLogs"`
}

type CMSGetAccessStatsReq struct {
	g.Meta    `path:"/getAccessStats" method:"post" summary:"CMS获取Access访问统计TopN" tags:"日志查询"`
	StartDate string `json:"startDate" v:"required|date-format:Y-m-d#开始日期不能为空|开始日期格式应为Y-m-d"`
	EndDate   string `json:"endDate" v:"required|date-format:Y-m-d#结束日期不能为空|结束日期格式应为Y-m-d"`
	TopN      int    `json:"topN"`
}

type TopStatItem struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

type CMSGetAccessStatsRes struct {
	UrlTop []TopStatItem `json:"urlTop"`
	IpTop  []TopStatItem `json:"ipTop"`
}

type CMSGetAccessTrendReq struct {
	g.Meta          `path:"/getAccessTrend" method:"post" summary:"CMS获取Access访问趋势" tags:"日志查询"`
	StartDate       string   `json:"startDate" v:"required|date-format:Y-m-d#开始日期不能为空|开始日期格式应为Y-m-d"`
	EndDate         string   `json:"endDate" v:"required|date-format:Y-m-d#结束日期不能为空|结束日期格式应为Y-m-d"`
	TraceId         string   `json:"traceId"`
	Url             string   `json:"url"`
	Ip              string   `json:"ip"`
	StatusCode      int      `json:"statusCode"`
	MinHandlerMs    *float64 `json:"minHandlerMs"`
	MaxHandlerMs    *float64 `json:"maxHandlerMs"`
	IntervalMinutes int      `json:"intervalMinutes" dc:"聚合粒度(分钟),0表示自动:1天=1分钟,3天=5分钟,7天=15分钟"`
}

type AccessTrendPoint struct {
	Time  string `json:"time"`
	Count int64  `json:"count"`
}

type CMSGetAccessTrendRes struct {
	IntervalMinutes int                `json:"intervalMinutes"`
	Points          []AccessTrendPoint `json:"points"`
	TotalCount      int64              `json:"totalCount"`
	PeakTime        string             `json:"peakTime"`
	PeakCount       int64              `json:"peakCount"`
}

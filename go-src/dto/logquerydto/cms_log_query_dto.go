package logquerydto

import (
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type CMSGetLogPathsReq struct {
	g.Meta `path:"/getLogPaths" method:"post" summary:"CMS获取日志目录配置" tags:"日志查询"`
}

type CMSGetLogPathsRes struct {
	ServerTime      string `json:"serverTime"`
	LogDir          string `json:"logDir"`
	AccessPrefix    string `json:"accessPrefix"`
	DetailPrefix    string `json:"detailPrefix"`
	ErrorPrefix     string `json:"errorPrefix"`
	ExportSubDir    string `json:"exportSubDir"`
	ExportURLPrefix string `json:"exportUrlPrefix"`
	LinuxOnly       bool   `json:"linuxOnly"`
}

type CMSLogQueryExportRes struct {
	ExportId  string `json:"exportId"`
	FileName  string `json:"fileName"`
	FileUrl   string `json:"fileUrl"`
	Total     int    `json:"total"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

type CMSDeleteLogQueryExportReq struct {
	g.Meta   `path:"/deleteExport" method:"post" summary:"CMS删除日志查询导出文件" tags:"日志查询"`
	ExportId string `json:"exportId" v:"required#导出ID不能为空"`
}

type CMSDeleteLogQueryExportRes struct {
	Success bool `json:"success"`
}

type CMSQueryDetailLogsReq struct {
	g.Meta `path:"/queryDetailLogs" method:"post" summary:"CMS查询详情日志" tags:"日志查询"`
	httpserver.CMSQueryReq
	StartDate string `json:"startDate" v:"required#开始时间不能为空" dc:"开始时间(Y-m-d或Y-m-d H:i:s)"`
	EndDate   string `json:"endDate" v:"required#结束时间不能为空" dc:"结束时间(Y-m-d或Y-m-d H:i:s)"`
	TraceId   string `json:"traceId"`
	ReqId     string `json:"reqId"`
	AuthId    string `json:"authId"`
	Url       string `json:"url"`
	Keyword   string `json:"keyword"`
}

type CMSQueryAccessLogsReq struct {
	g.Meta `path:"/queryAccessLogs" method:"post" summary:"CMS查询Access日志" tags:"日志查询"`
	httpserver.CMSQueryReq
	StartDate    string   `json:"startDate" v:"required#开始时间不能为空" dc:"开始时间(Y-m-d或Y-m-d H:i:s)"`
	EndDate      string   `json:"endDate" v:"required#结束时间不能为空" dc:"结束时间(Y-m-d或Y-m-d H:i:s)"`
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
	StartDate  string `json:"startDate" v:"required#开始时间不能为空" dc:"开始时间(Y-m-d或Y-m-d H:i:s)"`
	EndDate    string `json:"endDate" v:"required#结束时间不能为空" dc:"结束时间(Y-m-d或Y-m-d H:i:s)"`
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

type CMSGetAccessStatsReq struct {
	g.Meta    `path:"/getAccessStats" method:"post" summary:"CMS获取Access访问统计TopN" tags:"日志查询"`
	StartDate string `json:"startDate" v:"required#开始时间不能为空" dc:"开始时间(Y-m-d或Y-m-d H:i:s)"`
	EndDate   string `json:"endDate" v:"required#结束时间不能为空" dc:"结束时间(Y-m-d或Y-m-d H:i:s)"`
	TopN      int    `json:"topN"`
}

type CMSGetAccessTrendReq struct {
	g.Meta          `path:"/getAccessTrend" method:"post" summary:"CMS获取Access访问趋势" tags:"日志查询"`
	StartDate       string   `json:"startDate" v:"required#开始时间不能为空" dc:"开始时间(Y-m-d或Y-m-d H:i:s)"`
	EndDate         string   `json:"endDate" v:"required#结束时间不能为空" dc:"结束时间(Y-m-d或Y-m-d H:i:s)"`
	TraceId         string   `json:"traceId"`
	Url             string   `json:"url"`
	Ip              string   `json:"ip"`
	StatusCode      int      `json:"statusCode"`
	MinHandlerMs    *float64 `json:"minHandlerMs"`
	MaxHandlerMs    *float64 `json:"maxHandlerMs"`
	IntervalMinutes int      `json:"intervalMinutes" dc:"聚合粒度(分钟),0表示自动:1天=1分钟,3天=5分钟,7天=15分钟"`
}

type CMSSubmitLogQueryJobReq struct {
	g.Meta    `path:"/submitJob" method:"post" summary:"CMS提交异步日志查询任务" tags:"日志查询"`
	QueryType string          `json:"queryType" v:"required#查询类型不能为空"`
	Payload   json.RawMessage `json:"payload" v:"required#查询参数不能为空"`
}

type CMSSubmitLogQueryJobRes struct {
	JobId         string `json:"jobId"`
	QueuePosition int    `json:"queuePosition"`
}

type CMSGetLogQueryJobReq struct {
	g.Meta `path:"/getJob" method:"post" summary:"CMS获取异步日志查询结果" tags:"日志查询"`
	JobId  string `json:"jobId" v:"required#任务ID不能为空"`
}

type CMSGetLogQueryJobRes struct {
	JobId         string `json:"jobId"`
	QueryType     string `json:"queryType"`
	Status        string `json:"status"`
	QueuePosition int    `json:"queuePosition"`
	ErrorMessage  string `json:"errorMessage"`
	Result        any    `json:"result"`
}

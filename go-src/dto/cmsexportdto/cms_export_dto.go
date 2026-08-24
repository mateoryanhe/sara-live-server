package cmsexportdto

import (
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	ExportTypeLiveRecord                       = "liveRecord"
	ExportTypeLiveRevenueLog                   = "liveRevenueLog"
	ExportTypeVideoCallLog                     = "videoCallLog"
	ExportTypeAnchorIncomeSettlementLog        = "anchorIncomeSettlementLog"
	ExportTypeGuildIncomeSettlementLog         = "guildIncomeSettlementLog"
	ExportTypeGuildAnchorIncomeSettlementLog   = "guildAnchorIncomeSettlementLog"
	ExportTypeMyGuildAnchorIncomeSettlementLog = "myGuildAnchorIncomeSettlementLog"
	ExportTypeAnchorDailyEffectiveLive         = "anchorDailyEffectiveLive"
	ExportTypeGuildDailyEffectiveLive          = "guildDailyEffectiveLive"
	ExportTypeGuildAnchorDailyEffectiveLive    = "guildAnchorDailyEffectiveLive"
	ExportTypeMyGuildAnchorDailyEffectiveLive  = "myGuildAnchorDailyEffectiveLive"
	ExportTypeLiveDailyEffectiveLive             = "liveDailyEffectiveLive"
	ExportTypeCurrencyLog                      = "currencyLog"
)

type CMSExportResult struct {
	ExportId  string `json:"exportId"`
	FileName  string `json:"fileName"`
	FileUrl   string `json:"fileUrl"`
	Total     int    `json:"total"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

type CMSExportJobProgress struct {
	ExportedRows int `json:"exportedRows"`
	TotalRows    int `json:"totalRows"`
}

type CMSSubmitExportJobReq struct {
	g.Meta     `path:"/submitJob" method:"post" summary:"CMS提交异步导出任务" tags:"CMS导出"`
	ExportType string          `json:"exportType" v:"required#导出类型不能为空"`
	Payload    json.RawMessage `json:"payload" v:"required#导出参数不能为空"`
}

type CMSSubmitExportJobRes struct {
	JobId         string `json:"jobId"`
	QueuePosition int    `json:"queuePosition"`
}

type CMSGetExportJobReq struct {
	g.Meta `path:"/getJob" method:"post" summary:"CMS获取异步导出任务状态" tags:"CMS导出"`
	JobId  string `json:"jobId" v:"required#任务ID不能为空"`
}

type CMSGetExportJobRes struct {
	JobId         string                `json:"jobId"`
	ExportType    string                `json:"exportType"`
	Status        string                `json:"status"`
	QueuePosition int                   `json:"queuePosition"`
	ErrorMessage  string                `json:"errorMessage"`
	Progress      *CMSExportJobProgress `json:"progress"`
	Result        any                   `json:"result"`
}

type CMSDeleteExportReq struct {
	g.Meta   `path:"/deleteExport" method:"post" summary:"CMS删除导出文件" tags:"CMS导出"`
	ExportId string `json:"exportId" v:"required#导出ID不能为空"`
}

type CMSDeleteExportRes struct {
	Success bool `json:"success"`
}

type CMSExportHeadersPayload struct {
	Headers []string `json:"headers" v:"required#表头不能为空"`
}

type CMSExportLiveRecordPayload struct {
	CMSExportHeadersPayload
	AnchorId         string   `json:"anchorId"`
	PlatformAnchorId string   `json:"platformAnchorId"`
	GuildAnchorId    string   `json:"guildAnchorId"`
	AnchorIds        []string `json:"anchorIds"`
	StartTime        int64    `json:"startTime"`
	EndTime          int64    `json:"endTime"`
}

type CMSExportLiveRevenueLogPayload struct {
	CMSExportHeadersPayload
	ReceiverId       string   `json:"receiverId"`
	PlatformAnchorId string   `json:"platformAnchorId"`
	GuildAnchorId    string   `json:"guildAnchorId"`
	ReceiverIds      []string `json:"receiverIds"`
	LiveRecordId     string   `json:"liveRecordId"`
	RevenueType      uint8    `json:"revenueType"`
	StartTime        int64    `json:"startTime"`
	EndTime          int64    `json:"endTime"`
}

type CMSExportVideoCallLogPayload struct {
	CMSExportHeadersPayload
	CallerId   string `json:"callerId"`
	ReceiverId string `json:"receiverId"`
	Source     uint8  `json:"source"`
	Status     uint8  `json:"status"`
	StartTime  int64  `json:"startTime"`
	EndTime    int64  `json:"endTime"`
}

type CMSExportAnchorIncomeSettlementLogPayload struct {
	CMSExportHeadersPayload
	RoomId    string `json:"roomId"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
}

type CMSExportGuildIncomeSettlementLogPayload struct {
	CMSExportHeadersPayload
	GuildId   string `json:"guildId"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
}

type CMSExportGuildAnchorIncomeSettlementLogPayload struct {
	CMSExportHeadersPayload
	GuildId   string `json:"guildId"`
	RoomId    string `json:"roomId"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
}

type CMSExportMyGuildAnchorIncomeSettlementLogPayload struct {
	CMSExportHeadersPayload
	GuildId   string `json:"guildId"`
	RoomId    string `json:"roomId"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
}

type CMSExportDailyEffectiveLiveLabels struct {
	SettledYesText string `json:"settledYesText"`
	SettledNoText  string `json:"settledNoText"`
}

type CMSExportAnchorDailyEffectiveLivePayload struct {
	CMSExportHeadersPayload
	CMSExportDailyEffectiveLiveLabels
	AnchorId      uint64 `json:"anchorId"`
	LiveDateStart string `json:"liveDateStart"`
	LiveDateEnd   string `json:"liveDateEnd"`
	Settled       int8   `json:"settled"`
}

type CMSExportGuildDailyEffectiveLivePayload struct {
	CMSExportHeadersPayload
	CMSExportDailyEffectiveLiveLabels
	GuildId       uint64 `json:"guildId"`
	LiveDateStart string `json:"liveDateStart"`
	LiveDateEnd   string `json:"liveDateEnd"`
	Settled       int8   `json:"settled"`
}

type CMSExportGuildAnchorDailyEffectiveLivePayload struct {
	CMSExportHeadersPayload
	CMSExportDailyEffectiveLiveLabels
	GuildId       string `json:"guildId"`
	RoomId        string `json:"roomId"`
	LiveDateStart string `json:"liveDateStart"`
	LiveDateEnd   string `json:"liveDateEnd"`
	Settled       int8   `json:"settled"`
}

type CMSExportMyGuildAnchorDailyEffectiveLivePayload struct {
	CMSExportHeadersPayload
	CMSExportDailyEffectiveLiveLabels
	RoomId        string   `json:"roomId"`
	RoomIds       []string `json:"roomIds"`
	LiveDateStart string   `json:"liveDateStart"`
	LiveDateEnd   string   `json:"liveDateEnd"`
	Settled       int8     `json:"settled"`
}

type CMSExportLiveDailyEffectiveLivePayload struct {
	CMSExportHeadersPayload
	CMSExportDailyEffectiveLiveLabels
	AnchorId         string   `json:"anchorId"`
	PlatformAnchorId string   `json:"platformAnchorId"`
	GuildAnchorId    string   `json:"guildAnchorId"`
	AnchorIds        []string `json:"anchorIds"`
	LiveDateStart    string   `json:"liveDateStart"`
	LiveDateEnd      string   `json:"liveDateEnd"`
	Settled          int8     `json:"settled"`
}

type CMSExportCurrencyLogPayload struct {
	CMSExportHeadersPayload
	UserId       string `json:"userId"`
	CurrencyType uint8  `json:"currencyType"`
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
}

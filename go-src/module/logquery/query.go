package logquery

import (
	"context"
	"time"

	"xr-game-server/dto/logquerydto"
	"xr-game-server/errercode"
)

func GetLogPaths(_ context.Context, _ *logquerydto.CMSGetLogPathsReq) (*logquerydto.CMSGetLogPathsRes, error) {
	cfg := loadLogQueryConfig().normalized()
	return &logquerydto.CMSGetLogPathsRes{
		ServerTime:      time.Now().Format("2006-01-02 15:04:05.000"),
		LogDir:          cfg.LogDir,
		AccessPrefix:    cfg.AccessPrefix,
		DetailPrefix:    cfg.DetailPrefix,
		ErrorPrefix:     cfg.ErrorPrefix,
		ExportSubDir:    cfg.ExportSubDir,
		ExportURLPrefix: cfg.exportURLPrefix(),
		LinuxOnly:       true,
	}, nil
}

func QueryDetailLogs(_ context.Context, req *logquerydto.CMSQueryDetailLogsReq) (*logquerydto.CMSLogQueryExportRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	patterns := buildDetailPatterns(req.TraceId, req.ReqId, req.AuthId, req.Url, req.Keyword)
	result, err := createShellExport(logTypeDetail, patterns, req.StartDate, req.EndDate, req.PageIndex, req.PageSize)
	return toExportRes(result), err
}

func QueryAccessLogs(_ context.Context, req *logquerydto.CMSQueryAccessLogsReq) (*logquerydto.CMSLogQueryExportRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	patterns := buildAccessPatterns(req.TraceId, req.Url, req.Ip, req.StatusCode)
	result, err := createShellExport(logTypeAccess, patterns, req.StartDate, req.EndDate, req.PageIndex, req.PageSize)
	return toExportRes(result), err
}

func QueryErrorLogs(_ context.Context, req *logquerydto.CMSQueryErrorLogsReq) (*logquerydto.CMSLogQueryExportRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	patterns := buildErrorPatterns(req.TraceId, req.Url, req.Ip, req.Keyword, req.StatusCode)
	result, err := createShellExport(logTypeError, patterns, req.StartDate, req.EndDate, req.PageIndex, req.PageSize)
	return toExportRes(result), err
}

func GetTraceLogs(_ context.Context, req *logquerydto.CMSGetTraceLogsReq) (*logquerydto.CMSLogQueryExportRes, error) {
	if req == nil || req.TraceId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	result, err := createTraceShellExport(req.TraceId, req.StartDate, req.EndDate)
	return toExportRes(result), err
}

func GetAccessStats(_ context.Context, req *logquerydto.CMSGetAccessStatsReq) (*logquerydto.CMSLogQueryExportRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	result, err := createAccessStatsExport(req.StartDate, req.EndDate, req.TopN)
	return toExportRes(result), err
}

func GetAccessTrend(_ context.Context, req *logquerydto.CMSGetAccessTrendReq) (*logquerydto.CMSLogQueryExportRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	result, err := createAccessTrendExport(req)
	return toExportRes(result), err
}

func DeleteLogQueryExport(_ context.Context, req *logquerydto.CMSDeleteLogQueryExportReq) (*logquerydto.CMSDeleteLogQueryExportRes, error) {
	if req == nil || req.ExportId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := deleteExport(req.ExportId); err != nil {
		return nil, err
	}
	return &logquerydto.CMSDeleteLogQueryExportRes{Success: true}, nil
}

func toExportRes(result *shellExportResult) *logquerydto.CMSLogQueryExportRes {
	if result == nil {
		return nil
	}
	return &logquerydto.CMSLogQueryExportRes{
		ExportId:  result.ExportID,
		FileName:  result.FileName,
		FileUrl:   result.FileUrl,
		Total:     result.Total,
		PageIndex: result.PageIndex,
		PageSize:  result.PageSize,
	}
}

package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/logquerydto"
	"xr-game-server/module/logquery"
)

const LogQueryUrl = "/logQuery"

type LogQueryController struct{}

func initLogQueryController() {
	httpserver.RegCMS(LogQueryUrl, &LogQueryController{})
}

func (c *LogQueryController) GetLogPaths(ctx context.Context, req *logquerydto.CMSGetLogPathsReq) (res *logquerydto.CMSGetLogPathsRes, err error) {
	return logquery.GetLogPaths(ctx, req)
}

func (c *LogQueryController) QueryDetailLogs(ctx context.Context, req *logquerydto.CMSQueryDetailLogsReq) (res *httpserver.CMSQueryResp, err error) {
	return logquery.QueryDetailLogs(ctx, req)
}

func (c *LogQueryController) QueryAccessLogs(ctx context.Context, req *logquerydto.CMSQueryAccessLogsReq) (res *httpserver.CMSQueryResp, err error) {
	return logquery.QueryAccessLogs(ctx, req)
}

func (c *LogQueryController) QueryErrorLogs(ctx context.Context, req *logquerydto.CMSQueryErrorLogsReq) (res *httpserver.CMSQueryResp, err error) {
	return logquery.QueryErrorLogs(ctx, req)
}

func (c *LogQueryController) GetTraceLogs(ctx context.Context, req *logquerydto.CMSGetTraceLogsReq) (res *logquerydto.CMSGetTraceLogsRes, err error) {
	return logquery.GetTraceLogs(ctx, req)
}

func (c *LogQueryController) GetAccessStats(ctx context.Context, req *logquerydto.CMSGetAccessStatsReq) (res *logquerydto.CMSGetAccessStatsRes, err error) {
	return logquery.GetAccessStats(ctx, req)
}

func (c *LogQueryController) GetAccessTrend(ctx context.Context, req *logquerydto.CMSGetAccessTrendReq) (res *logquerydto.CMSGetAccessTrendRes, err error) {
	return logquery.GetAccessTrend(ctx, req)
}

func (c *LogQueryController) SubmitLogQueryJob(ctx context.Context, req *logquerydto.CMSSubmitLogQueryJobReq) (res *logquerydto.CMSSubmitLogQueryJobRes, err error) {
	return logquery.SubmitLogQueryJob(ctx, req)
}

func (c *LogQueryController) GetLogQueryJob(ctx context.Context, req *logquerydto.CMSGetLogQueryJobReq) (res *logquerydto.CMSGetLogQueryJobRes, err error) {
	return logquery.GetLogQueryJob(ctx, req)
}

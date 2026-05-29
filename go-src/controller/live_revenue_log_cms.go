package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liverevenuedto"
	"xr-game-server/module/liverevenue"
)

const (
	LiveRevenueLogUrl = "/liveRevenueLog"
)

type LiveRevenueLogController struct{}

func initLiveRevenueLogController() {
	httpserver.RegCMS(LiveRevenueLogUrl, &LiveRevenueLogController{})
}

// CMSLiveRevenueLogList CMS分页查询直播收益流水
func (c *LiveRevenueLogController) CMSLiveRevenueLogList(ctx context.Context, req *liverevenuedto.CMSLiveRevenueLogListReq) (res *httpserver.CMSQueryResp, err error) {
	return liverevenue.GetCMSList(ctx, req)
}

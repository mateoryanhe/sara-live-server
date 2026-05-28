package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/livegiftlogdto"
	"xr-game-server/module/livegiftlog"
)

const (
	LiveGiftLogUrl = "/liveGiftLog"
)

type LiveGiftLogController struct{}

func initLiveGiftLogController() {
	httpserver.RegCMS(LiveGiftLogUrl, &LiveGiftLogController{})
}

// CMSLiveGiftLogList CMS分页查询礼物流水
func (c *LiveGiftLogController) CMSLiveGiftLogList(ctx context.Context, req *livegiftlogdto.CMSLiveGiftLogListReq) (res *httpserver.CMSQueryResp, err error) {
	return livegiftlog.GetCMSList(ctx, req)
}

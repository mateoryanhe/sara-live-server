package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liverecorddto"
	"xr-game-server/module/liverecord"
)

const (
	LiveRecordCMSUrl = "/liveRecord"
)

type LiveRecordCMSController struct{}

func initLiveRecordCMSController() {
	httpserver.RegCMS(LiveRecordCMSUrl, &LiveRecordCMSController{})
}

// CMSLiveRecordList CMS分页查询直播记录
func (c *LiveRecordCMSController) CMSLiveRecordList(ctx context.Context, req *liverecorddto.CMSLiveRecordListReq) (res *httpserver.CMSQueryResp, err error) {
	return liverecord.GetCMSList(ctx, req)
}

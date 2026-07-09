package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/calldto"
	"xr-game-server/module/call"
)

const CallCMSUrl = "/call"

type CallCMSController struct{}

func initCallCMSController() {
	httpserver.RegCMS(CallCMSUrl, &CallCMSController{})
}

// CMSVideoCallLogList CMS分页查询视频通话日志
func (c *CallCMSController) CMSVideoCallLogList(ctx context.Context, req *calldto.CMSVideoCallLogListReq) (res *httpserver.CMSQueryResp, err error) {
	return call.GetCMSVideoCallLogList(ctx, req)
}

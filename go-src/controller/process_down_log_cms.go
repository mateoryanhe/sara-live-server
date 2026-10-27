package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/processdowndto"
	"xr-game-server/module/processdown"
)

const ProcessDownLogCMSUrl = "/processDownLog"

type ProcessDownLogCMSController struct{}

func initProcessDownLogCMSController() {
	httpserver.RegCMS(ProcessDownLogCMSUrl, &ProcessDownLogCMSController{})
}

func (c *ProcessDownLogCMSController) ListProcessDownLog(ctx context.Context, req *processdowndto.CMSListProcessDownLogReq) (*httpserver.CMSQueryResp, error) {
	return processdown.GetCMSList(ctx, req)
}

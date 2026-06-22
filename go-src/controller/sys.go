package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/sysdto"
	"xr-game-server/module/sys"
)

const SysInfoUrl = "/sysInfo"

type SysInfoController struct{}

func initSysController() {
	httpserver.RegNonAuthAPI(SysInfoUrl, &SysInfoController{})
}

func (SysInfoController) GetInfo(ctx context.Context, req *sysdto.SysCfgReq) (*sysdto.SysCfgResp, error) {
	return sys.GetSysCfg(ctx, req)
}

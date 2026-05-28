package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/statdto"
	"xr-game-server/module/stat"
)

const SysStatUrl = "/sysStat"

type SysStatController struct{}

func initSysStatController() {
	httpserver.RegCMS(SysStatUrl, &SysStatController{})
}

// GetSysStat CMS获取系统总数据
func (c *SysStatController) GetSysStat(ctx context.Context, req *statdto.CMSSysStatReq) (res *statdto.CMSSysStatRes, err error) {
	return stat.GetCMSSysStat(ctx, req)
}

// GetUserStatTrend CMS获取用户数据趋势
func (c *SysStatController) GetUserStatTrend(ctx context.Context, req *statdto.CMSUserStatTrendReq) (res *statdto.CMSUserStatTrendRes, err error) {
	return stat.GetCMSUserStatTrend(ctx, req)
}

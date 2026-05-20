package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/rechargecfgdto"
	"xr-game-server/module/rechargecfg"
)

const (
	RechargeCfgUrl = "/rechargeCfg"
)

type RechargeCfgController struct{}

func initRechargeCfgController() {
	httpserver.RegCMS(RechargeCfgUrl, &RechargeCfgController{})
}

// RechargeCfgList 分页获取充值配置列表
func (c *RechargeCfgController) RechargeCfgList(ctx context.Context, req *rechargecfgdto.RechargeCfgListReq) (res *httpserver.CMSQueryResp, err error) {
	return rechargecfg.GetList(ctx, req)
}

// CreateRechargeCfg 新建
func (c *RechargeCfgController) CreateRechargeCfg(ctx context.Context, req *rechargecfgdto.CreateRechargeCfgReq) (res *rechargecfgdto.CreateRechargeCfgRes, err error) {
	return rechargecfg.Create(ctx, req)
}

// UpdateRechargeCfg 修改
func (c *RechargeCfgController) UpdateRechargeCfg(ctx context.Context, req *rechargecfgdto.UpdateRechargeCfgReq) (res *rechargecfgdto.UpdateRechargeCfgRes, err error) {
	return rechargecfg.Update(ctx, req)
}

// DeleteRechargeCfg 删除
func (c *RechargeCfgController) DeleteRechargeCfg(ctx context.Context, req *rechargecfgdto.DeleteRechargeCfgReq) (res *rechargecfgdto.DeleteRechargeCfgRes, err error) {
	return rechargecfg.Delete(ctx, req)
}

// OnShelfRechargeCfg 上架
func (c *RechargeCfgController) OnShelfRechargeCfg(ctx context.Context, req *rechargecfgdto.OnShelfRechargeCfgReq) (res *rechargecfgdto.OnShelfRechargeCfgRes, err error) {
	return rechargecfg.OnShelf(ctx, req)
}

// OffShelfRechargeCfg 下架
func (c *RechargeCfgController) OffShelfRechargeCfg(ctx context.Context, req *rechargecfgdto.OffShelfRechargeCfgReq) (res *rechargecfgdto.OffShelfRechargeCfgRes, err error) {
	return rechargecfg.OffShelf(ctx, req)
}

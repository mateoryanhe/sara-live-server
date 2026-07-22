package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/apppkgdto"
	"xr-game-server/module/apppkg"
)

const AppPkgUrl = "/appPkg"

type AppPkgController struct{}

func initAppPkgController() {
	httpserver.RegCMS(AppPkgUrl, &AppPkgController{})
}

func (c *AppPkgController) AppPkgList(ctx context.Context, req *apppkgdto.AppPkgListReq) (*httpserver.CMSQueryResp, error) {
	return apppkg.GetList(ctx, req)
}

func (c *AppPkgController) CreateAppPkg(ctx context.Context, req *apppkgdto.CreateAppPkgReq) (*apppkgdto.CreateAppPkgRes, error) {
	return apppkg.Create(ctx, req)
}

func (c *AppPkgController) UpdateAppPkg(ctx context.Context, req *apppkgdto.UpdateAppPkgReq) (*apppkgdto.UpdateAppPkgRes, error) {
	return apppkg.Update(ctx, req)
}

func (c *AppPkgController) DeleteAppPkg(ctx context.Context, req *apppkgdto.DeleteAppPkgReq) (*apppkgdto.DeleteAppPkgRes, error) {
	return apppkg.Delete(ctx, req)
}

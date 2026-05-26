package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/diamonddto"
	"xr-game-server/module/userinfo"
)

const DiamondUrl = "/diamond"

type DiamondController struct{}

func initDiamondController() {
	httpserver.RegCMS(DiamondUrl, &DiamondController{})
}

func (c *DiamondController) Add(ctx context.Context, req *diamonddto.CMSAddDiamondReq) (*diamonddto.CMSAddDiamondRes, error) {
	return userinfo.DiamondCMSAdd(ctx, req)
}

func (c *DiamondController) Sub(ctx context.Context, req *diamonddto.CMSSubDiamondReq) (*diamonddto.CMSSubDiamondRes, error) {
	return userinfo.DiamondCMSSub(ctx, req)
}

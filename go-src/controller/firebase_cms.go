package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/firebasedto"
	"xr-game-server/module/auth"
)

const FirebaseCMSUrl = "/firebase"

type FirebaseCMSController struct{}

func initFirebaseCMSController() {
	httpserver.RegCMS(FirebaseCMSUrl, &FirebaseCMSController{})
}

func (c *FirebaseCMSController) GetFirebaseCfg(ctx context.Context, req *firebasedto.GetFirebaseCfgReq) (*firebasedto.GetFirebaseCfgRes, error) {
	return auth.GetFirebaseCfg(ctx, req)
}

func (c *FirebaseCMSController) SaveFirebaseCfg(ctx context.Context, req *firebasedto.SaveFirebaseCfgReq) (*firebasedto.SaveFirebaseCfgRes, error) {
	return auth.SaveFirebaseCfg(ctx, req)
}

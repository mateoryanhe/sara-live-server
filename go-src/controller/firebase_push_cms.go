package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/firebasepushdto"
	"xr-game-server/module/firebasepush"
)

const FirebasePushCMSUrl = "/firebasePush"

type FirebasePushCMSController struct{}

func initFirebasePushCMSController() {
	httpserver.RegCMS(FirebasePushCMSUrl, &FirebasePushCMSController{})
}

func (c *FirebasePushCMSController) GetFirebasePushCfg(ctx context.Context, req *firebasepushdto.GetFirebasePushCfgReq) (*firebasepushdto.GetFirebasePushCfgRes, error) {
	return firebasepush.GetFirebasePushCfg(ctx, req)
}

func (c *FirebasePushCMSController) SaveFirebasePushCfg(ctx context.Context, req *firebasepushdto.SaveFirebasePushCfgReq) (*firebasepushdto.SaveFirebasePushCfgRes, error) {
	return firebasepush.SaveFirebasePushCfg(ctx, req)
}

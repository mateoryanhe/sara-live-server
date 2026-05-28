package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/apptokendto"
	"xr-game-server/module/auth"
)

type AppTokenController struct {
}

func initAppTokenController() {
	httpserver.RegCMS("/appToken", &AppTokenController{})
}

func (a *AppTokenController) GetAppToken(ctx context.Context, req *apptokendto.GetAppTokenReq) (*httpserver.CMSQueryResp, error) {
	return auth.GetAppToken(ctx, req)
}

func (a *AppTokenController) SaveAppToken(ctx context.Context, req *apptokendto.SaveAppTokenReq) (bool, error) {
	return auth.SaveAppToken(ctx, req)
}

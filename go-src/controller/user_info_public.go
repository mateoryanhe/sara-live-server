package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/module/userinfo"
)

type UserInfoPublicController struct{}

func initUserInfoPublicController() {
	httpserver.RegNonAuthAPI(UserInfoUrl, &UserInfoPublicController{})
}

func (c *UserInfoPublicController) CancelAccountByCode(ctx context.Context, req *userinfodto.CancelAccountByCodeReq) (*userinfodto.CancelAccountByCodeRes, error) {
	return userinfo.CancelAccountByCode(ctx, req)
}

func (c *UserInfoPublicController) CancelAccountByPhone(ctx context.Context, req *userinfodto.CancelAccountByPhoneReq) (*userinfodto.CancelAccountByPhoneRes, error) {
	return userinfo.CancelAccountByPhone(ctx, req)
}

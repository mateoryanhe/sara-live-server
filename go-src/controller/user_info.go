package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/module/currencylog"
	"xr-game-server/module/userinfo"
)

const (
	UserInfoUrl = "/userInfo"
)

type UserInfoController struct {
}

func initUserInfoController() {
	httpserver.RegAPI(UserInfoUrl, &UserInfoController{})
}

func (c *UserInfoController) Get(ctx context.Context, req *userinfodto.GetUserInfoReq) (res *userinfodto.GetUserInfoRes, err error) {
	return userinfo.GetUserInfo(ctx, req)
}

func (c *UserInfoController) UpdateNickname(ctx context.Context, req *userinfodto.UpdateNicknameReq) (res *userinfodto.UpdateNicknameRes, err error) {
	return userinfo.UpdateNickname(ctx, req)
}

func (c *UserInfoController) GetCurrencyLog(ctx context.Context, req *userinfodto.GetCurrencyLogReq) (res *userinfodto.GetCurrencyLogRes, err error) {
	return currencylog.GetByUserId(ctx, req)
}

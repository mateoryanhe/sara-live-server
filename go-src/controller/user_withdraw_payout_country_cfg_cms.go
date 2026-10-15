package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/haipaydto"
	"xr-game-server/module/userwithdraw"
)

const UserWithdrawPayoutCountryCfgCMSURL = "/userWithdrawPayoutCountryCfg"

type UserWithdrawPayoutCountryCfgCMSController struct{}

func initUserWithdrawPayoutCountryCfgCMSController() {
	httpserver.RegCMS(UserWithdrawPayoutCountryCfgCMSURL, &UserWithdrawPayoutCountryCfgCMSController{})
}

func (c *UserWithdrawPayoutCountryCfgCMSController) GetPayoutCountryCfg(ctx context.Context, req *haipaydto.GetUserWithdrawPayoutCountryCfgReq) (*haipaydto.GetCollectionCountryCfgRes, error) {
	return userwithdraw.GetHaiPayUserWithdrawPayoutCountryCfg(ctx, req)
}

func (c *UserWithdrawPayoutCountryCfgCMSController) SavePayoutCountryCfg(ctx context.Context, req *haipaydto.SaveUserWithdrawPayoutCountryCfgReq) (*haipaydto.SaveCollectionCountryCfgRes, error) {
	return userwithdraw.SaveHaiPayUserWithdrawPayoutCountryCfg(ctx, req)
}

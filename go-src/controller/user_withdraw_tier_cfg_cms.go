package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/userwithdrawtiercfgdto"
	"xr-game-server/module/userwithdraw"
)

const UserWithdrawTierCfgCMSUrl = "/userWithdrawTierCfg"

type UserWithdrawTierCfgCMSController struct{}

func initUserWithdrawTierCfgCMSController() {
	httpserver.RegCMS(UserWithdrawTierCfgCMSUrl, &UserWithdrawTierCfgCMSController{})
}

func (c *UserWithdrawTierCfgCMSController) UserWithdrawTierCfgList(ctx context.Context, req *userwithdrawtiercfgdto.UserWithdrawTierCfgListReq) (*httpserver.CMSQueryResp, error) {
	return userwithdraw.ListUserWithdrawTierCfgs(ctx, req)
}

func (c *UserWithdrawTierCfgCMSController) CreateUserWithdrawTierCfg(ctx context.Context, req *userwithdrawtiercfgdto.CreateUserWithdrawTierCfgReq) (*userwithdrawtiercfgdto.CreateUserWithdrawTierCfgRes, error) {
	return userwithdraw.CreateUserWithdrawTierCfg(ctx, req)
}

func (c *UserWithdrawTierCfgCMSController) UpdateUserWithdrawTierCfg(ctx context.Context, req *userwithdrawtiercfgdto.UpdateUserWithdrawTierCfgReq) (*userwithdrawtiercfgdto.UpdateUserWithdrawTierCfgRes, error) {
	return userwithdraw.UpdateUserWithdrawTierCfg(ctx, req)
}

func (c *UserWithdrawTierCfgCMSController) DeleteUserWithdrawTierCfg(ctx context.Context, req *userwithdrawtiercfgdto.DeleteUserWithdrawTierCfgReq) (*userwithdrawtiercfgdto.DeleteUserWithdrawTierCfgRes, error) {
	return userwithdraw.DeleteUserWithdrawTierCfg(ctx, req)
}

func (c *UserWithdrawTierCfgCMSController) OnShelfUserWithdrawTierCfg(ctx context.Context, req *userwithdrawtiercfgdto.OnShelfUserWithdrawTierCfgReq) (*userwithdrawtiercfgdto.OnShelfUserWithdrawTierCfgRes, error) {
	return userwithdraw.OnShelfUserWithdrawTierCfg(ctx, req)
}

func (c *UserWithdrawTierCfgCMSController) OffShelfUserWithdrawTierCfg(ctx context.Context, req *userwithdrawtiercfgdto.OffShelfUserWithdrawTierCfgReq) (*userwithdrawtiercfgdto.OffShelfUserWithdrawTierCfgRes, error) {
	return userwithdraw.OffShelfUserWithdrawTierCfg(ctx, req)
}

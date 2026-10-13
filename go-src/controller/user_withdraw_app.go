package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/userwithdrawdto"
	"xr-game-server/module/userwithdraw"
)

const UserWithdrawAppUrl = "/userWithdraw"

type UserWithdrawAppController struct{}

func initUserWithdrawAppController() {
	httpserver.RegAPI(UserWithdrawAppUrl, &UserWithdrawAppController{})
}

func (c *UserWithdrawAppController) GetWithdrawCfg(ctx context.Context, req *userwithdrawdto.AppGetWithdrawCfgReq) (*userwithdrawdto.AppGetWithdrawCfgRes, error) {
	return userwithdraw.GetWithdrawCfg(ctx, req)
}

func (c *UserWithdrawAppController) GetPayoutInfo(ctx context.Context, req *userwithdrawdto.AppGetPayoutInfoReq) (*userwithdrawdto.AppGetPayoutInfoRes, error) {
	return userwithdraw.AppGetPayoutInfo(ctx, req)
}

func (c *UserWithdrawAppController) SavePayoutInfo(ctx context.Context, req *userwithdrawdto.AppSavePayoutInfoReq) (*userwithdrawdto.AppSavePayoutInfoRes, error) {
	return userwithdraw.AppSavePayoutInfo(ctx, req)
}

func (c *UserWithdrawAppController) ApplyWithdraw(ctx context.Context, req *userwithdrawdto.AppApplyWithdrawReq) (*userwithdrawdto.AppApplyWithdrawRes, error) {
	return userwithdraw.AppApplyWithdraw(ctx, req)
}

func (c *UserWithdrawAppController) ListWithdrawLogs(ctx context.Context, req *userwithdrawdto.AppWithdrawLogListReq) (*userwithdrawdto.AppWithdrawLogListRes, error) {
	return userwithdraw.AppListWithdrawLogs(ctx, req)
}

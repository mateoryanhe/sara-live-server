package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/userwithdrawdto"
	"xr-game-server/module/userwithdraw"
)

const UserWithdrawCMSUrl = "/userWithdrawCms"

type UserWithdrawCMSController struct{}

func initUserWithdrawCMSController() {
	httpserver.RegCMS(UserWithdrawCMSUrl, &UserWithdrawCMSController{})
}

func (c *UserWithdrawCMSController) ListWithdrawLogs(ctx context.Context, req *userwithdrawdto.CMSListWithdrawLogsReq) (*userwithdrawdto.CMSListWithdrawLogsRes, error) {
	return userwithdraw.CMSListWithdrawLogs(ctx, req)
}

func (c *UserWithdrawCMSController) BatchApproveWithdraw(ctx context.Context, req *userwithdrawdto.CMSBatchApproveWithdrawReq) (*userwithdrawdto.CMSBatchApproveWithdrawRes, error) {
	return userwithdraw.CMSBatchApproveWithdraw(ctx, req)
}

func (c *UserWithdrawCMSController) BatchRejectWithdraw(ctx context.Context, req *userwithdrawdto.CMSBatchRejectWithdrawReq) (*userwithdrawdto.CMSBatchRejectWithdrawRes, error) {
	return userwithdraw.CMSBatchRejectWithdraw(ctx, req)
}

func (c *UserWithdrawCMSController) BatchRevokeWithdrawApproval(ctx context.Context, req *userwithdrawdto.CMSBatchRevokeWithdrawApprovalReq) (*userwithdrawdto.CMSBatchRevokeWithdrawApprovalRes, error) {
	return userwithdraw.CMSBatchRevokeWithdrawApproval(ctx, req)
}

func (c *UserWithdrawCMSController) BatchSubmitPayout(ctx context.Context, req *userwithdrawdto.CMSBatchSubmitPayoutReq) (*userwithdrawdto.CMSBatchSubmitPayoutRes, error) {
	return userwithdraw.CMSBatchSubmitPayout(ctx, req)
}

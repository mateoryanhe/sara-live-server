package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/userwithdrawdeploydto"
	"xr-game-server/module/userwithdrawdeploy"
)

const UserWithdrawDeployCMSUrl = "/userWithdrawDeploy"

type UserWithdrawDeployCMSController struct{}

func initUserWithdrawDeployCMSController() {
	httpserver.RegCMSHandler(UserWithdrawDeployCMSUrl, "/deployZip", handleUserWithdrawDeployZip)
	httpserver.RegCMS(UserWithdrawDeployCMSUrl, &UserWithdrawDeployCMSController{})
}

func handleUserWithdrawDeployZip(r *ghttp.Request) {
	res, err := userwithdrawdeploy.DeployZipFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *UserWithdrawDeployCMSController) GetUserWithdrawDeployInfo(ctx context.Context, req *userwithdrawdeploydto.GetUserWithdrawDeployInfoReq) (*userwithdrawdeploydto.GetUserWithdrawDeployInfoRes, error) {
	return userwithdrawdeploy.GetUserWithdrawDeployInfo(ctx, req)
}

func (c *UserWithdrawDeployCMSController) SaveUserWithdrawDeployCfg(ctx context.Context, req *userwithdrawdeploydto.SaveUserWithdrawDeployCfgReq) (*userwithdrawdeploydto.SaveUserWithdrawDeployCfgRes, error) {
	return userwithdrawdeploy.SaveUserWithdrawDeployCfg(ctx, req)
}

package userwithdrawdeploy

import (
	"context"
	"strings"

	"xr-game-server/core/cfg"
	"xr-game-server/dto/userwithdrawdeploydto"
	"xr-game-server/errercode"
	"xr-game-server/module/domainsite"
)

func GetUserWithdrawDeployInfo(ctx context.Context, _ *userwithdrawdeploydto.GetUserWithdrawDeployInfoReq) (*userwithdrawdeploydto.GetUserWithdrawDeployInfoRes, error) {
	deployPath, err := getDeployDir()
	if err != nil {
		return nil, err
	}
	return &userwithdrawdeploydto.GetUserWithdrawDeployInfoRes{
		Info: &userwithdrawdeploydto.UserWithdrawDeployInfoItem{
			Domain:       cfg.GetStaticSiteDomains(userwithdrawdeploydto.UserWithdrawStaticPrefix),
			UrlPrefix:    userwithdrawdeploydto.UserWithdrawStaticPrefix,
			DeployPath:   deployPath,
			AcceptExt:    ".zip",
			LastUploadAt: domainsite.GetLastUploadAt(ctx, userwithdrawdeploydto.UserWithdrawSiteKey),
		},
	}, nil
}

func SaveUserWithdrawDeployCfg(ctx context.Context, req *userwithdrawdeploydto.SaveUserWithdrawDeployCfgReq) (*userwithdrawdeploydto.SaveUserWithdrawDeployCfgRes, error) {
	domain := strings.TrimSpace(req.Domain)
	deployPath := strings.TrimSpace(req.DeployPath)
	if domain == "" || deployPath == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if _, err := domainsite.SaveStaticSiteMapping(
		ctx,
		userwithdrawdeploydto.UserWithdrawSiteKey,
		userwithdrawdeploydto.UserWithdrawStaticPrefix,
		domain,
		deployPath,
	); err != nil {
		return nil, err
	}
	return &userwithdrawdeploydto.SaveUserWithdrawDeployCfgRes{Success: true}, nil
}

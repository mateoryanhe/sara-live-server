package thirdpaydeploy

import (
	"context"
	"strings"

	"xr-game-server/core/cfg"
	"xr-game-server/dto/thirdpaydeploydto"
	"xr-game-server/errercode"
	"xr-game-server/module/domainsite"
)

func GetThirdPayDeployInfo(ctx context.Context, _ *thirdpaydeploydto.GetThirdPayDeployInfoReq) (*thirdpaydeploydto.GetThirdPayDeployInfoRes, error) {
	deployPath, err := getDeployDir()
	if err != nil {
		return nil, err
	}
	return &thirdpaydeploydto.GetThirdPayDeployInfoRes{
		Info: &thirdpaydeploydto.ThirdPayDeployInfoItem{
			Domain:       cfg.GetStaticSiteDomains(thirdpaydeploydto.ThirdPayStaticPrefix),
			UrlPrefix:    thirdpaydeploydto.ThirdPayStaticPrefix,
			DeployPath:   deployPath,
			AcceptExt:    ".zip",
			LastUploadAt: domainsite.GetLastUploadAt(ctx, thirdpaydeploydto.ThirdPaySiteKey),
		},
	}, nil
}

func SaveThirdPayDeployCfg(ctx context.Context, req *thirdpaydeploydto.SaveThirdPayDeployCfgReq) (*thirdpaydeploydto.SaveThirdPayDeployCfgRes, error) {
	domain := strings.TrimSpace(req.Domain)
	deployPath := strings.TrimSpace(req.DeployPath)
	if domain == "" || deployPath == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if _, err := domainsite.SaveStaticSiteMapping(
		ctx,
		thirdpaydeploydto.ThirdPaySiteKey,
		thirdpaydeploydto.ThirdPayStaticPrefix,
		domain,
		deployPath,
	); err != nil {
		return nil, err
	}
	return &thirdpaydeploydto.SaveThirdPayDeployCfgRes{Success: true}, nil
}

package thirdpaydeploy

import (
	"context"

	"xr-game-server/dto/thirdpaydeploydto"
)

func GetThirdPayDeployInfo(_ context.Context, _ *thirdpaydeploydto.GetThirdPayDeployInfoReq) (*thirdpaydeploydto.GetThirdPayDeployInfoRes, error) {
	deployPath, err := getDeployDir()
	if err != nil {
		return nil, err
	}
	return &thirdpaydeploydto.GetThirdPayDeployInfoRes{
		Info: &thirdpaydeploydto.ThirdPayDeployInfoItem{
			UrlPrefix:  thirdpaydeploydto.ThirdPayStaticPrefix,
			DeployPath: deployPath,
			AcceptExt:  ".zip",
		},
	}, nil
}

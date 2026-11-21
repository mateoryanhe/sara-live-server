package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/audienceonlinegreetdto"
	"xr-game-server/module/audienceonlinegreet"
)

const AudienceOnlineAnchorGreetCMSUrl = "/audienceOnlineAnchorGreet"

type AudienceOnlineAnchorGreetCMSController struct{}

func initAudienceOnlineAnchorGreetCMSController() {
	httpserver.RegCMS(AudienceOnlineAnchorGreetCMSUrl, &AudienceOnlineAnchorGreetCMSController{})
}

func (c *AudienceOnlineAnchorGreetCMSController) GetAudienceOnlineAnchorGreetCfg(ctx context.Context, req *audienceonlinegreetdto.GetAudienceOnlineAnchorGreetCfgReq) (*audienceonlinegreetdto.GetAudienceOnlineAnchorGreetCfgRes, error) {
	return audienceonlinegreet.GetCfg(ctx, req)
}

func (c *AudienceOnlineAnchorGreetCMSController) SaveAudienceOnlineAnchorGreetCfg(ctx context.Context, req *audienceonlinegreetdto.SaveAudienceOnlineAnchorGreetCfgReq) (*audienceonlinegreetdto.SaveAudienceOnlineAnchorGreetCfgRes, error) {
	return audienceonlinegreet.SaveCfg(ctx, req)
}
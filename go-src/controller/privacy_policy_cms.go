package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/privacypolicydto"
	"xr-game-server/module/privacypolicy"
)

const PrivacyPolicyCMSUrl = "/privacyPolicy"

type PrivacyPolicyCMSController struct{}

func initPrivacyPolicyCMSController() {
	httpserver.RegCMS(PrivacyPolicyCMSUrl, &PrivacyPolicyCMSController{})
}

func (c *PrivacyPolicyCMSController) GetPrivacyPolicyCfg(ctx context.Context, req *privacypolicydto.GetPrivacyPolicyCfgReq) (*privacypolicydto.GetPrivacyPolicyCfgRes, error) {
	return privacypolicy.GetPrivacyPolicyCfg(ctx, req)
}

func (c *PrivacyPolicyCMSController) SavePrivacyPolicyCfg(ctx context.Context, req *privacypolicydto.SavePrivacyPolicyCfgReq) (*privacypolicydto.SavePrivacyPolicyCfgRes, error) {
	return privacypolicy.SavePrivacyPolicyCfg(ctx, req)
}

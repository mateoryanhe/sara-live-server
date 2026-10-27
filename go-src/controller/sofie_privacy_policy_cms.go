package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/privacypolicydto"
	"xr-game-server/module/privacypolicy"
)

const SofiePrivacyPolicyCMSUrl = "/sofiePrivacyPolicy"

type SofiePrivacyPolicyCMSController struct{}

func initSofiePrivacyPolicyCMSController() {
	httpserver.RegCMS(SofiePrivacyPolicyCMSUrl, &SofiePrivacyPolicyCMSController{})
}

func (c *SofiePrivacyPolicyCMSController) GetSofiePrivacyPolicyCfg(ctx context.Context, req *privacypolicydto.GetSofiePrivacyPolicyCfgReq) (*privacypolicydto.GetSofiePrivacyPolicyCfgRes, error) {
	return privacypolicy.GetSofiePrivacyPolicyCfg(ctx, req)
}

func (c *SofiePrivacyPolicyCMSController) SaveSofiePrivacyPolicyCfg(ctx context.Context, req *privacypolicydto.SaveSofiePrivacyPolicyCfgReq) (*privacypolicydto.SaveSofiePrivacyPolicyCfgRes, error) {
	return privacypolicy.SaveSofiePrivacyPolicyCfg(ctx, req)
}

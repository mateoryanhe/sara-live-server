package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/firebaseanalyticsdto"
	"xr-game-server/module/firebaseanalytics"
)

type FirebaseAnalyticsAppController struct{}

func initFirebaseAnalyticsAppController() {
	httpserver.RegNonAuthAPI(FirebaseAnalyticsCMSUrl, &FirebaseAnalyticsAppController{})
}

func (c *FirebaseAnalyticsAppController) GetClientCfgForApp(ctx context.Context, req *firebaseanalyticsdto.GetFirebaseAnalyticsClientCfgForAppReq) (*firebaseanalyticsdto.GetFirebaseAnalyticsClientCfgForAppRes, error) {
	return firebaseanalytics.GetClientCfgForApp(ctx, req)
}

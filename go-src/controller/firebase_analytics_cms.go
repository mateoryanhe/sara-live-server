package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/firebaseanalyticsdto"
	"xr-game-server/module/firebaseanalytics"
)

const FirebaseAnalyticsCMSUrl = "/firebaseAnalytics"

type FirebaseAnalyticsCMSController struct{}

func initFirebaseAnalyticsCMSController() {
	httpserver.RegCMS(FirebaseAnalyticsCMSUrl, &FirebaseAnalyticsCMSController{})
}

func (c *FirebaseAnalyticsCMSController) GetFirebaseAnalyticsCfg(ctx context.Context, req *firebaseanalyticsdto.GetFirebaseAnalyticsCfgReq) (*firebaseanalyticsdto.GetFirebaseAnalyticsCfgRes, error) {
	return firebaseanalytics.GetFirebaseAnalyticsCfg(ctx, req)
}

func (c *FirebaseAnalyticsCMSController) SaveFirebaseAnalyticsCfg(ctx context.Context, req *firebaseanalyticsdto.SaveFirebaseAnalyticsCfgReq) (*firebaseanalyticsdto.SaveFirebaseAnalyticsCfgRes, error) {
	return firebaseanalytics.SaveFirebaseAnalyticsCfg(ctx, req)
}

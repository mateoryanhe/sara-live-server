package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/firebasedto"
	"xr-game-server/module/auth"
)

type FirebaseAppController struct{}

func initFirebaseAppController() {
	httpserver.RegNonAuthAPI(FirebaseCMSUrl, &FirebaseAppController{})
}

func (c *FirebaseAppController) GetClientCfgForApp(ctx context.Context, req *firebasedto.GetFirebaseClientCfgForAppReq) (*firebasedto.GetFirebaseClientCfgForAppRes, error) {
	return auth.GetFirebaseClientCfgForApp(ctx, req)
}

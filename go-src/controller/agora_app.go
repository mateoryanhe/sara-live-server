package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/agoradto"
	"xr-game-server/module/agora"
)

const AgoraAppUrl = "/agora"

type AgoraAppController struct{}

func initAgoraAppController() {
	httpserver.RegAPI(AgoraAppUrl, &AgoraAppController{})
}

// LiveRoomToken App端上报频道名获取声网Token
func (c *AgoraAppController) LiveRoomToken(ctx context.Context, req *agoradto.GetLiveRoomTokenReq) (*agoradto.GetLiveRoomTokenRes, error) {
	return agora.GetLiveRoomToken(ctx, req)
}

// AppId App端获取声网AppId
func (c *AgoraAppController) AppId(ctx context.Context, req *agoradto.GetAppIdReq) (*agoradto.GetAppIdRes, error) {
	return agora.GetAppId(ctx, req)
}

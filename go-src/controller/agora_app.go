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

// LiveRoomToken App端获取进直播间声网Token
func (c *AgoraAppController) LiveRoomToken(ctx context.Context, req *agoradto.GetLiveRoomTokenReq) (*agoradto.GetLiveRoomTokenRes, error) {
	return agora.GetLiveRoomToken(ctx, req)
}

// CheckUserOnline 查询用户是否在直播间(声网)
func (c *AgoraAppController) CheckUserOnline(ctx context.Context, req *agoradto.CheckUserOnlineReq) (*agoradto.CheckUserOnlineRes, error) {
	return agora.CheckUserOnline(ctx, req)
}

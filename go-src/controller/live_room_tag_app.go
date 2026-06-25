package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/module/liveroom"
)

const LiveRoomTagAppUrl = "/liveRoomTag"

type LiveRoomTagAppController struct{}

func initLiveRoomTagAppController() {
	httpserver.RegAPI(LiveRoomTagAppUrl, &LiveRoomTagAppController{})
}

func (c *LiveRoomTagAppController) AppLiveRoomTagList(ctx context.Context, req *liveroomdto.AppLiveRoomTagListReq) (*liveroomdto.AppLiveRoomTagListRes, error) {
	return liveroom.GetAppLiveRoomTagList(ctx, req)
}

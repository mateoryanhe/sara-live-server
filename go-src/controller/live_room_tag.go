package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/module/liveroom"
)

const LiveRoomTagUrl = "/liveRoomTag"

type LiveRoomTagController struct{}

func initLiveRoomTagController() {
	httpserver.RegCMS(LiveRoomTagUrl, &LiveRoomTagController{})
}

func (c *LiveRoomTagController) LiveRoomTagList(ctx context.Context, req *liveroomdto.LiveRoomTagListReq) (*httpserver.CMSQueryResp, error) {
	return liveroom.GetLiveRoomTagList(ctx, req)
}

func (c *LiveRoomTagController) CreateLiveRoomTag(ctx context.Context, req *liveroomdto.CreateLiveRoomTagReq) (*liveroomdto.CreateLiveRoomTagRes, error) {
	return liveroom.CreateLiveRoomTag(ctx, req)
}

func (c *LiveRoomTagController) UpdateLiveRoomTag(ctx context.Context, req *liveroomdto.UpdateLiveRoomTagReq) (*liveroomdto.UpdateLiveRoomTagRes, error) {
	return liveroom.UpdateLiveRoomTag(ctx, req)
}

func (c *LiveRoomTagController) DeleteLiveRoomTag(ctx context.Context, req *liveroomdto.DeleteLiveRoomTagReq) (*liveroomdto.DeleteLiveRoomTagRes, error) {
	return liveroom.DeleteLiveRoomTag(ctx, req)
}

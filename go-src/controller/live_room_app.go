package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/module/liveroom"
)

const (
	LiveRoomAppUrl = "/liveRoom"
)

type LiveRoomAppController struct{}

func initLiveRoomAppController() {
	httpserver.RegAPI(LiveRoomAppUrl, &LiveRoomAppController{})
}

// CreateRoom 创建直播间
func (c *LiveRoomAppController) CreateRoom(ctx context.Context, req *liveroomdto.CreateLiveRoomReq) (res *liveroomdto.CreateLiveRoomRes, err error) {
	return liveroom.CreateRoom(ctx, req)
}

// StartLive 开播
func (c *LiveRoomAppController) StartLive(ctx context.Context, req *liveroomdto.StartLiveReq) (res *liveroomdto.StartLiveRes, err error) {
	return liveroom.StartLive(ctx, req)
}

// StopLive 下播
func (c *LiveRoomAppController) StopLive(ctx context.Context, req *liveroomdto.StopLiveReq) (res *liveroomdto.StopLiveRes, err error) {
	return liveroom.StopLive(ctx, req)
}

// UpdateCover 修改封面
func (c *LiveRoomAppController) UpdateCover(ctx context.Context, req *liveroomdto.UpdateCoverReq) (res *liveroomdto.UpdateCoverRes, err error) {
	return liveroom.UpdateCover(ctx, req)
}

// UpdateNotice 修改公告
func (c *LiveRoomAppController) UpdateNotice(ctx context.Context, req *liveroomdto.UpdateNoticeReq) (res *liveroomdto.UpdateNoticeRes, err error) {
	return liveroom.UpdateNotice(ctx, req)
}

// GetRoom 查询直播间
func (c *LiveRoomAppController) GetRoom(ctx context.Context, req *liveroomdto.GetLiveRoomReq) (res *liveroomdto.GetLiveRoomRes, err error) {
	return liveroom.GetRoom(ctx, req)
}

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

// JoinRoom 加入直播间
func (c *LiveRoomAppController) JoinRoom(ctx context.Context, req *liveroomdto.JoinRoomReq) (res *liveroomdto.JoinRoomRes, err error) {
	return liveroom.JoinRoom(ctx, req)
}

// LeaveRoom 离开直播间
func (c *LiveRoomAppController) LeaveRoom(ctx context.Context, req *liveroomdto.LeaveRoomReq) (res *liveroomdto.LeaveRoomRes, err error) {
	return liveroom.LeaveRoom(ctx, req)
}

// GetOnlineUserList 分页查询直播间在线玩家
func (c *LiveRoomAppController) GetOnlineUserList(ctx context.Context, req *liveroomdto.GetOnlineUserListReq) (res *liveroomdto.GetOnlineUserListRes, err error) {
	return liveroom.GetOnlineUserList(ctx, req)
}

// GetRoom 查询直播间
func (c *LiveRoomAppController) GetRoom(ctx context.Context, req *liveroomdto.GetLiveRoomReq) (res *liveroomdto.GetLiveRoomRes, err error) {
	return liveroom.GetRoom(ctx, req)
}

// SendGift 直播间送礼
func (c *LiveRoomAppController) SendGift(ctx context.Context, req *liveroomdto.SendGiftReq) (res *liveroomdto.SendGiftRes, err error) {
	return liveroom.SendGift(ctx, req)
}

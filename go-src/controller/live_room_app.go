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

// GetVipOnlineUserList 分页查询直播间VIP在线玩家
func (c *LiveRoomAppController) GetVipOnlineUserList(ctx context.Context, req *liveroomdto.GetVipOnlineUserListReq) (res *liveroomdto.GetVipOnlineUserListRes, err error) {
	return liveroom.GetVipOnlineUserList(ctx, req)
}

// GetContributionRank 查询直播间观众贡献榜
func (c *LiveRoomAppController) GetContributionRank(ctx context.Context, req *liveroomdto.GetContributionRankReq) (res *liveroomdto.GetContributionRankRes, err error) {
	return liveroom.GetContributionRank(ctx, req)
}

// GetRoom 查询直播间
func (c *LiveRoomAppController) GetRoom(ctx context.Context, req *liveroomdto.GetLiveRoomReq) (res *liveroomdto.GetLiveRoomRes, err error) {
	return liveroom.GetRoom(ctx, req)
}

// RoomList 分页查询直播间列表
func (c *LiveRoomAppController) RoomList(ctx context.Context, req *liveroomdto.GetLiveRoomListReq) (res *liveroomdto.GetLiveRoomListRes, err error) {
	return liveroom.GetRoomList(ctx, req)
}

// NearbyRoomList 以当前直播间为锚点查询相邻直播中直播间
func (c *LiveRoomAppController) NearbyRoomList(ctx context.Context, req *liveroomdto.GetNearbyLiveRoomListReq) (res *liveroomdto.GetNearbyLiveRoomListRes, err error) {
	return liveroom.GetNearbyLiveRoomList(ctx, req)
}

// HotRoomList 分页查询 Hot 分类直播中房间列表
func (c *LiveRoomAppController) HotRoomList(ctx context.Context, req *liveroomdto.GetHotLiveRoomListReq) (res *liveroomdto.GetHotLiveRoomListRes, err error) {
	return liveroom.GetHotLiveRoomList(ctx, req)
}

// FollowedRoomList 分页查询我关注的直播间列表
func (c *LiveRoomAppController) FollowedRoomList(ctx context.Context, req *liveroomdto.GetFollowedLiveRoomListReq) (res *liveroomdto.GetLiveRoomListRes, err error) {
	return liveroom.GetFollowedRoomList(ctx, req)
}

// SendGift 直播间送礼
func (c *LiveRoomAppController) SendGift(ctx context.Context, req *liveroomdto.SendGiftReq) (res *liveroomdto.SendGiftRes, err error) {
	return liveroom.SendGift(ctx, req)
}

// SendGiftToAnchor 给指定主播送礼
func (c *LiveRoomAppController) SendGiftToAnchor(ctx context.Context, req *liveroomdto.SendGiftToAnchorReq) (res *liveroomdto.SendGiftToAnchorRes, err error) {
	return liveroom.SendGiftToAnchor(ctx, req)
}

// SendChat 直播间文字消息
func (c *LiveRoomAppController) SendChat(ctx context.Context, req *liveroomdto.SendChatReq) (res *liveroomdto.SendChatRes, err error) {
	return liveroom.SendChat(ctx, req)
}

// SendPrivateRoomChat 私密房文字消息
func (c *LiveRoomAppController) SendPrivateRoomChat(ctx context.Context, req *liveroomdto.SendPrivateRoomChatReq) (res *liveroomdto.SendPrivateRoomChatRes, err error) {
	return liveroom.SendPrivateRoomChat(ctx, req)
}

// SendPaidDanmaku 直播间付费弹幕
func (c *LiveRoomAppController) SendPaidDanmaku(ctx context.Context, req *liveroomdto.SendPaidDanmakuReq) (res *liveroomdto.SendPaidDanmakuRes, err error) {
	return liveroom.SendPaidDanmaku(ctx, req)
}

// SetAudienceMute 主播对指定观众禁言/解禁
func (c *LiveRoomAppController) SetAudienceMute(ctx context.Context, req *liveroomdto.SetAudienceMuteReq) (res *liveroomdto.SetAudienceMuteRes, err error) {
	return liveroom.SetAudienceMute(ctx, req)
}

// CancelAudienceMute 主播取消观众禁言
func (c *LiveRoomAppController) CancelAudienceMute(ctx context.Context, req *liveroomdto.CancelAudienceMuteReq) (res *liveroomdto.CancelAudienceMuteRes, err error) {
	return liveroom.CancelAudienceMute(ctx, req)
}

// GetAudienceRestrictStatus 查询观众禁言与被踢状态
func (c *LiveRoomAppController) GetAudienceRestrictStatus(ctx context.Context, req *liveroomdto.GetAudienceRestrictStatusReq) (res *liveroomdto.GetAudienceRestrictStatusRes, err error) {
	return liveroom.GetAudienceRestrictStatus(ctx, req)
}

// KickAudience 主播踢出指定观众
func (c *LiveRoomAppController) KickAudience(ctx context.Context, req *liveroomdto.KickAudienceReq) (res *liveroomdto.KickAudienceRes, err error) {
	return liveroom.KickAudience(ctx, req)
}

// CancelKickBan 主播取消观众进入限制
func (c *LiveRoomAppController) CancelKickBan(ctx context.Context, req *liveroomdto.CancelKickBanReq) (res *liveroomdto.CancelKickBanRes, err error) {
	return liveroom.CancelKickBan(ctx, req)
}

// ReportLiveStartStatus 主播上报开播状态
func (c *LiveRoomAppController) ReportLiveStartStatus(ctx context.Context, req *liveroomdto.ReportLiveStartStatusReq) (res *liveroomdto.ReportLiveStartStatusRes, err error) {
	return liveroom.ReportLiveStartStatus(ctx, req)
}

// GameRecommendList 查询直播间推荐游戏列表
func (c *LiveRoomAppController) GameRecommendList(ctx context.Context, req *liveroomdto.GetLiveRoomGameRecommendListReq) (*liveroomdto.GetLiveRoomGameRecommendListRes, error) {
	return liveroom.GetLiveRoomGameRecommendList(ctx, req)
}

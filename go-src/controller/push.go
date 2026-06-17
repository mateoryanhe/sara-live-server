package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/pushdto"
)

type PushController struct{}

func initPushController() {
	httpserver.RegNonAuthAPI("/push", new(PushController))
}

func (c *PushController) Enter(_ context.Context, _ *pushdto.EnterPushReq) (*pushdto.EnterPushResp, error) {
	return nil, nil
}

func (c *PushController) Heart(_ context.Context, _ *pushdto.HeartPushReq) (*pushdto.HeartPushResp, error) {
	return nil, nil
}

func (c *PushController) Kick(_ context.Context, _ *pushdto.KickPushReq) (*pushdto.KickPushResp, error) {
	return nil, nil
}

func (c *PushController) CloseServer(_ context.Context, _ *pushdto.CloseServerPushReq) (*pushdto.CloseServerPushResp, error) {
	return nil, nil
}

func (c *PushController) RepeatLogin(_ context.Context, _ *pushdto.RepeatLoginPushReq) (*pushdto.RepeatLoginPushResp, error) {
	return nil, nil
}

func (c *PushController) Error(_ context.Context, _ *pushdto.ErrorPushReq) (*pushdto.ErrorPushResp, error) {
	return nil, nil
}

func (c *PushController) ErrorParam(_ context.Context, _ *pushdto.ErrorParamPushReq) (*pushdto.ErrorParamPushResp, error) {
	return nil, nil
}

func (c *PushController) LiveRoomGift(_ context.Context, _ *pushdto.LiveRoomGiftPushReq) (*pushdto.LiveRoomGiftPushResp, error) {
	return nil, nil
}

func (c *PushController) LiveRoomChat(_ context.Context, _ *pushdto.LiveRoomChatPushReq) (*pushdto.LiveRoomChatPushResp, error) {
	return nil, nil
}

func (c *PushController) Diamond(_ context.Context, _ *pushdto.DiamondPushReq) (*pushdto.DiamondPushResp, error) {
	return nil, nil
}

func (c *PushController) Gold(_ context.Context, _ *pushdto.GoldPushReq) (*pushdto.GoldPushResp, error) {
	return nil, nil
}

func (c *PushController) VipLevel(_ context.Context, _ *pushdto.VipLevelPushReq) (*pushdto.VipLevelPushResp, error) {
	return nil, nil
}

func (c *PushController) PrivateMessage(_ context.Context, _ *pushdto.PrivateMessagePushReq) (*pushdto.PrivateMessagePushResp, error) {
	return nil, nil
}

func (c *PushController) SystemMessage(_ context.Context, _ *pushdto.SystemMessagePushReq) (*pushdto.SystemMessagePushResp, error) {
	return nil, nil
}

func (c *PushController) LiveRoomAnchorBan(_ context.Context, _ *pushdto.LiveRoomAnchorBanPushReq) (*pushdto.LiveRoomAnchorBanPushResp, error) {
	return nil, nil
}

func (c *PushController) LiveRoomAudienceMute(_ context.Context, _ *pushdto.LiveRoomAudienceMutePushReq) (*pushdto.LiveRoomAudienceMutePushResp, error) {
	return nil, nil
}

func (c *PushController) LiveRoomAudienceKick(_ context.Context, _ *pushdto.LiveRoomAudienceKickPushReq) (*pushdto.LiveRoomAudienceKickPushResp, error) {
	return nil, nil
}

func (c *PushController) LiveRoomAudienceKickCancel(_ context.Context, _ *pushdto.LiveRoomAudienceKickCancelPushReq) (*pushdto.LiveRoomAudienceKickCancelPushResp, error) {
	return nil, nil
}

func (c *PushController) LiveRoomAudienceJoin(_ context.Context, _ *pushdto.LiveRoomAudienceJoinPushReq) (*pushdto.LiveRoomAudienceJoinPushResp, error) {
	return nil, nil
}

func (c *PushController) LiveRoomAudienceLeave(_ context.Context, _ *pushdto.LiveRoomAudienceLeavePushReq) (*pushdto.LiveRoomAudienceLeavePushResp, error) {
	return nil, nil
}

func (c *PushController) LiveRoomPaidDanmaku(_ context.Context, _ *pushdto.LiveRoomPaidDanmakuPushReq) (*pushdto.LiveRoomPaidDanmakuPushResp, error) {
	return nil, nil
}

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

// 以下 handler 仅用于 OpenAPI/Swagger 文档展示,按业务模块划分,模块内 cmd 降序.

// --- 推送-通话 ---

// LiveRoomCallAnchorAcceptedAudience cmd=33 主播开始接听视频通话(推送给直播间观众)
func (c *PushController) LiveRoomCallAnchorAcceptedAudience(_ context.Context, _ *pushdto.LiveRoomCallAnchorAcceptedAudiencePushReq) (*pushdto.LiveRoomCallAnchorAcceptedAudiencePushResp, error) {
	return nil, nil
}

// LiveRoomCallTimeout cmd=32 直播间通话呼叫超时
func (c *PushController) LiveRoomCallTimeout(_ context.Context, _ *pushdto.LiveRoomCallTimeoutPushReq) (*pushdto.LiveRoomCallTimeoutPushResp, error) {
	return nil, nil
}

// LiveRoomCallStarted cmd=31 直播间通话开始
func (c *PushController) LiveRoomCallStarted(_ context.Context, _ *pushdto.LiveRoomCallStartedPushReq) (*pushdto.LiveRoomCallStartedPushResp, error) {
	return nil, nil
}

// LiveRoomCallEnded cmd=28 直播间通话结束
func (c *PushController) LiveRoomCallEnded(_ context.Context, _ *pushdto.LiveRoomCallEndedPushReq) (*pushdto.LiveRoomCallEndedPushResp, error) {
	return nil, nil
}

// LiveRoomCallAccepted cmd=27 直播间通话被接听
func (c *PushController) LiveRoomCallAccepted(_ context.Context, _ *pushdto.LiveRoomCallAcceptedPushReq) (*pushdto.LiveRoomCallAcceptedPushResp, error) {
	return nil, nil
}

// LiveRoomCallRejected cmd=26 直播间通话被拒接
func (c *PushController) LiveRoomCallRejected(_ context.Context, _ *pushdto.LiveRoomCallRejectedPushReq) (*pushdto.LiveRoomCallRejectedPushResp, error) {
	return nil, nil
}

// LiveRoomCallRequest cmd=25 直播间通话请求
func (c *PushController) LiveRoomCallRequest(_ context.Context, _ *pushdto.LiveRoomCallRequestPushReq) (*pushdto.LiveRoomCallRequestPushResp, error) {
	return nil, nil
}

// --- 推送-私密房 ---

// LiveRoomPrivateGift cmd=30 给指定主播送礼
func (c *PushController) LiveRoomPrivateGift(_ context.Context, _ *pushdto.LiveRoomPrivateGiftPushReq) (*pushdto.LiveRoomPrivateGiftPushResp, error) {
	return nil, nil
}

// LiveRoomPrivateChat cmd=29 私密房文字消息
func (c *PushController) LiveRoomPrivateChat(_ context.Context, _ *pushdto.LiveRoomPrivateChatPushReq) (*pushdto.LiveRoomPrivateChatPushResp, error) {
	return nil, nil
}

// --- 推送-关注 ---

// LiveFollowCount cmd=35 关注数/粉丝数
func (c *PushController) LiveFollowCount(_ context.Context, _ *pushdto.LiveFollowCountPushReq) (*pushdto.LiveFollowCountPushResp, error) {
	return nil, nil
}

// --- 推送-直播间 ---

// LiveRoomStartLive cmd=34 主播开播(全服广播)
func (c *PushController) LiveRoomStartLive(_ context.Context, _ *pushdto.LiveRoomStartLivePushReq) (*pushdto.LiveRoomStartLivePushResp, error) {
	return nil, nil
}

// LiveRoomTotalIncome cmd=37 本场直播总收益
func (c *PushController) LiveRoomTotalIncome(_ context.Context, _ *pushdto.LiveRoomTotalIncomePushReq) (*pushdto.LiveRoomTotalIncomePushResp, error) {
	return nil, nil
}

// LiveRoomAudienceListRefresh cmd=24 观众列表刷新
func (c *PushController) LiveRoomAudienceListRefresh(_ context.Context, _ *pushdto.LiveRoomAudienceListRefreshPushReq) (*pushdto.LiveRoomAudienceListRefreshPushResp, error) {
	return nil, nil
}

// LiveRoomStopLive cmd=23 主播下播(全服广播)
func (c *PushController) LiveRoomStopLive(_ context.Context, _ *pushdto.LiveRoomStopLivePushReq) (*pushdto.LiveRoomStopLivePushResp, error) {
	return nil, nil
}

// LiveRoomPaidDanmaku cmd=22 付费弹幕
func (c *PushController) LiveRoomPaidDanmaku(_ context.Context, _ *pushdto.LiveRoomPaidDanmakuPushReq) (*pushdto.LiveRoomPaidDanmakuPushResp, error) {
	return nil, nil
}

// LiveRoomAudienceLeave cmd=21 观众离房
func (c *PushController) LiveRoomAudienceLeave(_ context.Context, _ *pushdto.LiveRoomAudienceLeavePushReq) (*pushdto.LiveRoomAudienceLeavePushResp, error) {
	return nil, nil
}

// LiveRoomAudienceJoin cmd=20 观众进房
func (c *PushController) LiveRoomAudienceJoin(_ context.Context, _ *pushdto.LiveRoomAudienceJoinPushReq) (*pushdto.LiveRoomAudienceJoinPushResp, error) {
	return nil, nil
}

// LiveRoomAudienceKickCancel cmd=19 取消进入限制
func (c *PushController) LiveRoomAudienceKickCancel(_ context.Context, _ *pushdto.LiveRoomAudienceKickCancelPushReq) (*pushdto.LiveRoomAudienceKickCancelPushResp, error) {
	return nil, nil
}

// LiveRoomAudienceKick cmd=18 观众被踢
func (c *PushController) LiveRoomAudienceKick(_ context.Context, _ *pushdto.LiveRoomAudienceKickPushReq) (*pushdto.LiveRoomAudienceKickPushResp, error) {
	return nil, nil
}

// LiveRoomAudienceMute cmd=17 观众禁言
func (c *PushController) LiveRoomAudienceMute(_ context.Context, _ *pushdto.LiveRoomAudienceMutePushReq) (*pushdto.LiveRoomAudienceMutePushResp, error) {
	return nil, nil
}

// LiveRoomAnchorBan cmd=15 主播封禁
func (c *PushController) LiveRoomAnchorBan(_ context.Context, _ *pushdto.LiveRoomAnchorBanPushReq) (*pushdto.LiveRoomAnchorBanPushResp, error) {
	return nil, nil
}

// LiveRoomChat cmd=9 直播间文字消息
func (c *PushController) LiveRoomChat(_ context.Context, _ *pushdto.LiveRoomChatPushReq) (*pushdto.LiveRoomChatPushResp, error) {
	return nil, nil
}

// LiveRoomGift cmd=8 直播间送礼
func (c *PushController) LiveRoomGift(_ context.Context, _ *pushdto.LiveRoomGiftPushReq) (*pushdto.LiveRoomGiftPushResp, error) {
	return nil, nil
}

// --- 推送-消息 ---

// ActivityMessage cmd=36 活动消息发布(全服广播)
func (c *PushController) ActivityMessage(_ context.Context, _ *pushdto.ActivityMessagePushReq) (*pushdto.ActivityMessagePushResp, error) {
	return nil, nil
}

// SystemMessage cmd=14 系统消息
func (c *PushController) SystemMessage(_ context.Context, _ *pushdto.SystemMessagePushReq) (*pushdto.SystemMessagePushResp, error) {
	return nil, nil
}

// --- 推送-私信 ---

// PrivateMessage cmd=13 私信
func (c *PushController) PrivateMessage(_ context.Context, _ *pushdto.PrivateMessagePushReq) (*pushdto.PrivateMessagePushResp, error) {
	return nil, nil
}

// --- 推送-钱包 ---

// VipLevel cmd=12 VIP等级
func (c *PushController) VipLevel(_ context.Context, _ *pushdto.VipLevelPushReq) (*pushdto.VipLevelPushResp, error) {
	return nil, nil
}

// Gold cmd=11 金币余额
func (c *PushController) Gold(_ context.Context, _ *pushdto.GoldPushReq) (*pushdto.GoldPushResp, error) {
	return nil, nil
}

// FirstRechargeSuccess cmd=38 首充成功
func (c *PushController) FirstRechargeSuccess(_ context.Context, _ *pushdto.FirstRechargeSuccessPushReq) (*pushdto.FirstRechargeSuccessPushResp, error) {
	return nil, nil
}

// Diamond cmd=10 钻石余额
func (c *PushController) Diamond(_ context.Context, _ *pushdto.DiamondPushReq) (*pushdto.DiamondPushResp, error) {
	return nil, nil
}

// --- 推送-错误 ---

// ErrorParam cmd=7 带参数错误
func (c *PushController) ErrorParam(_ context.Context, _ *pushdto.ErrorParamPushReq) (*pushdto.ErrorParamPushResp, error) {
	return nil, nil
}

// Error cmd=6 无参数错误
func (c *PushController) Error(_ context.Context, _ *pushdto.ErrorPushReq) (*pushdto.ErrorPushResp, error) {
	return nil, nil
}

// --- 推送-连接 ---

// RepeatLogin cmd=5 重复登录
func (c *PushController) RepeatLogin(_ context.Context, _ *pushdto.RepeatLoginPushReq) (*pushdto.RepeatLoginPushResp, error) {
	return nil, nil
}

// CloseServer cmd=4 关服通知
func (c *PushController) CloseServer(_ context.Context, _ *pushdto.CloseServerPushReq) (*pushdto.CloseServerPushResp, error) {
	return nil, nil
}

// Kick cmd=3 踢下线
func (c *PushController) Kick(_ context.Context, _ *pushdto.KickPushReq) (*pushdto.KickPushResp, error) {
	return nil, nil
}

// Heart cmd=2 服务端心跳
func (c *PushController) Heart(_ context.Context, _ *pushdto.HeartPushReq) (*pushdto.HeartPushResp, error) {
	return nil, nil
}

// Enter cmd=1 连接鉴权失败
func (c *PushController) Enter(_ context.Context, _ *pushdto.EnterPushReq) (*pushdto.EnterPushResp, error) {
	return nil, nil
}

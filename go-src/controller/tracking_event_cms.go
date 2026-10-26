package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/trackingdto"
	"xr-game-server/module/tracking"
)

const TrackingEventCMSUrl = "/trackingEvent"

type TrackingEventCMSController struct{}

func initTrackingEventCMSController() {
	httpserver.RegCMS(TrackingEventCMSUrl, &TrackingEventCMSController{})
}

func (c *TrackingEventCMSController) GetHotLiveRoomJoinTrend(ctx context.Context, req *trackingdto.CMSHotLiveRoomJoinTrendReq) (*trackingdto.CMSHotLiveRoomJoinTrendRes, error) {
	return tracking.GetCMSHotLiveRoomJoinTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetGameLiveRoomJoinTrend(ctx context.Context, req *trackingdto.CMSGameLiveRoomJoinTrendReq) (*trackingdto.CMSGameLiveRoomJoinTrendRes, error) {
	return tracking.GetCMSGameLiveRoomJoinTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetHotLiveRoomGiftClickTrend(ctx context.Context, req *trackingdto.CMSHotLiveRoomGiftClickTrendReq) (*trackingdto.CMSHotLiveRoomGiftClickTrendRes, error) {
	return tracking.GetCMSHotLiveRoomGiftClickTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetLiveFirstFrameRenderTrend(ctx context.Context, req *trackingdto.CMSLiveFirstFrameRenderTrendReq) (*trackingdto.CMSLiveFirstFrameRenderTrendRes, error) {
	return tracking.GetCMSLiveFirstFrameRenderTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetGameLiveWatchStreamTrend(ctx context.Context, req *trackingdto.CMSGameLiveWatchStreamTrendReq) (*trackingdto.CMSGameLiveWatchStreamTrendRes, error) {
	return tracking.GetCMSGameLiveWatchStreamTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetHotLiveRoomLeaveTrend(ctx context.Context, req *trackingdto.CMSHotLiveRoomLeaveTrendReq) (*trackingdto.CMSHotLiveRoomLeaveTrendRes, error) {
	return tracking.GetCMSHotLiveRoomLeaveTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetCall1v1InitiateTrend(ctx context.Context, req *trackingdto.CMSCall1v1InitiateTrendReq) (*trackingdto.CMSCall1v1InitiateTrendRes, error) {
	return tracking.GetCMSCall1v1InitiateTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetCall1v1RoomCallTrend(ctx context.Context, req *trackingdto.CMSCall1v1RoomCallTrendReq) (*trackingdto.CMSCall1v1RoomCallTrendRes, error) {
	return tracking.GetCMSCall1v1RoomCallTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetCall1v1ConnectSuccessTrend(ctx context.Context, req *trackingdto.CMSCall1v1ConnectSuccessTrendReq) (*trackingdto.CMSCall1v1ConnectSuccessTrendRes, error) {
	return tracking.GetCMSCall1v1ConnectSuccessTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetOneToOneRoomListTrend(ctx context.Context, req *trackingdto.CMSOneToOneRoomListTrendReq) (*trackingdto.CMSOneToOneRoomListTrendRes, error) {
	return tracking.GetCMSOneToOneRoomListTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetMiniGameRoundStartTrend(ctx context.Context, req *trackingdto.CMSMiniGameRoundStartTrendReq) (*trackingdto.CMSMiniGameRoundStartTrendRes, error) {
	return tracking.GetCMSMiniGameRoundStartTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetMiniGameRoundResultTrend(ctx context.Context, req *trackingdto.CMSMiniGameRoundResultTrendReq) (*trackingdto.CMSMiniGameRoundResultTrendRes, error) {
	return tracking.GetCMSMiniGameRoundResultTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetMiniGameExposureTrend(ctx context.Context, req *trackingdto.CMSMiniGameExposureTrendReq) (*trackingdto.CMSMiniGameExposureTrendRes, error) {
	return tracking.GetCMSMiniGameExposureTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetMiniGameWebViewLoadTrend(ctx context.Context, req *trackingdto.CMSMiniGameWebViewLoadTrendReq) (*trackingdto.CMSMiniGameWebViewLoadTrendRes, error) {
	return tracking.GetCMSMiniGameWebViewLoadTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetMiniGameControllerClickTrend(ctx context.Context, req *trackingdto.CMSMiniGameControllerClickTrendReq) (*trackingdto.CMSMiniGameControllerClickTrendRes, error) {
	return tracking.GetCMSMiniGameControllerClickTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetMiniGameSpinClickTrend(ctx context.Context, req *trackingdto.CMSMiniGameSpinClickTrendReq) (*trackingdto.CMSMiniGameSpinClickTrendRes, error) {
	return tracking.GetCMSMiniGameSpinClickTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetMiniGamePlusClickTrend(ctx context.Context, req *trackingdto.CMSMiniGamePlusClickTrendReq) (*trackingdto.CMSMiniGamePlusClickTrendRes, error) {
	return tracking.GetCMSMiniGamePlusClickTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetOfflineGameRoomPlayTrend(ctx context.Context, req *trackingdto.CMSOfflineGameRoomPlayTrendReq) (*trackingdto.CMSOfflineGameRoomPlayTrendRes, error) {
	return tracking.GetCMSOfflineGameRoomPlayTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetShortVideoPlayStartTrend(ctx context.Context, req *trackingdto.CMSShortVideoPlayStartTrendReq) (*trackingdto.CMSShortVideoPlayStartTrendRes, error) {
	return tracking.GetCMSShortVideoPlayStartTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetShortVideoPaywallClickTrend(ctx context.Context, req *trackingdto.CMSShortVideoPaywallClickTrendReq) (*trackingdto.CMSShortVideoPaywallClickTrendRes, error) {
	return tracking.GetCMSShortVideoPaywallClickTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetShortVideoPaySuccessTrend(ctx context.Context, req *trackingdto.CMSShortVideoPaySuccessTrendReq) (*trackingdto.CMSShortVideoPaySuccessTrendRes, error) {
	return tracking.GetCMSShortVideoPaySuccessTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetShortVideoWatchEndTrend(ctx context.Context, req *trackingdto.CMSShortVideoWatchEndTrendReq) (*trackingdto.CMSShortVideoWatchEndTrendRes, error) {
	return tracking.GetCMSShortVideoWatchEndTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetVoiceRoomEnterTrend(ctx context.Context, req *trackingdto.CMSVoiceRoomEnterTrendReq) (*trackingdto.CMSVoiceRoomEnterTrendRes, error) {
	return tracking.GetCMSVoiceRoomEnterTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetVoiceRoomMicRequestTrend(ctx context.Context, req *trackingdto.CMSVoiceRoomMicRequestTrendReq) (*trackingdto.CMSVoiceRoomMicRequestTrendRes, error) {
	return tracking.GetCMSVoiceRoomMicRequestTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetVoiceRoomGiftClickTrend(ctx context.Context, req *trackingdto.CMSVoiceRoomGiftClickTrendReq) (*trackingdto.CMSVoiceRoomGiftClickTrendRes, error) {
	return tracking.GetCMSVoiceRoomGiftClickTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetChargePanelEnterTrend(ctx context.Context, req *trackingdto.CMSChargePanelEnterTrendReq) (*trackingdto.CMSChargePanelEnterTrendRes, error) {
	return tracking.GetCMSChargePanelEnterTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetChargeClickTrend(ctx context.Context, req *trackingdto.CMSChargeClickTrendReq) (*trackingdto.CMSChargeClickTrendRes, error) {
	return tracking.GetCMSChargeClickTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetChargeResultServerTrend(ctx context.Context, req *trackingdto.CMSChargeResultServerTrendReq) (*trackingdto.CMSChargeResultServerTrendRes, error) {
	return tracking.GetCMSChargeResultServerTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetH5DownloadTrend(ctx context.Context, req *trackingdto.CMSH5DownloadTrendReq) (*trackingdto.CMSH5DownloadTrendRes, error) {
	return tracking.GetCMSH5DownloadTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetOfficialSiteDownloadTrend(ctx context.Context, req *trackingdto.CMSOfficialSiteDownloadTrendReq) (*trackingdto.CMSOfficialSiteDownloadTrendRes, error) {
	return tracking.GetCMSOfficialSiteDownloadTrend(ctx, req)
}

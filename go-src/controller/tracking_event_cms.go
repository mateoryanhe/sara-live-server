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

func (c *TrackingEventCMSController) GetLiveFirstFrameRenderTrend(ctx context.Context, req *trackingdto.CMSLiveFirstFrameRenderTrendReq) (*trackingdto.CMSLiveFirstFrameRenderTrendRes, error) {
	return tracking.GetCMSLiveFirstFrameRenderTrend(ctx, req)
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

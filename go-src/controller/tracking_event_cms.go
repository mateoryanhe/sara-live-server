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

func (c *TrackingEventCMSController) GetLiveFirstFrameRenderTrend(ctx context.Context, req *trackingdto.CMSLiveFirstFrameRenderTrendReq) (*trackingdto.CMSLiveFirstFrameRenderTrendRes, error) {
	return tracking.GetCMSLiveFirstFrameRenderTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetHotLiveRoomLeaveTrend(ctx context.Context, req *trackingdto.CMSHotLiveRoomLeaveTrendReq) (*trackingdto.CMSHotLiveRoomLeaveTrendRes, error) {
	return tracking.GetCMSHotLiveRoomLeaveTrend(ctx, req)
}

func (c *TrackingEventCMSController) GetCall1v1InitiateTrend(ctx context.Context, req *trackingdto.CMSCall1v1InitiateTrendReq) (*trackingdto.CMSCall1v1InitiateTrendRes, error) {
	return tracking.GetCMSCall1v1InitiateTrend(ctx, req)
}

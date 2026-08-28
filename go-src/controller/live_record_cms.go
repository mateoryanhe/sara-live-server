package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liverecorddto"
	"xr-game-server/module/liverecord"
)

const (
	LiveRecordCMSUrl = "/liveRecord"
)

type LiveRecordCMSController struct{}

func initLiveRecordCMSController() {
	httpserver.RegCMS(LiveRecordCMSUrl, &LiveRecordCMSController{})
}

// CMSLiveRecordList CMS分页查询直播记录
func (c *LiveRecordCMSController) CMSLiveRecordList(ctx context.Context, req *liverecorddto.CMSLiveRecordListReq) (res *httpserver.CMSQueryResp, err error) {
	return liverecord.GetCMSList(ctx, req)
}

// CMSDailyEffectiveLiveList CMS分页查询每日流水
func (c *LiveRecordCMSController) CMSDailyEffectiveLiveList(ctx context.Context, req *liverecorddto.CMSDailyEffectiveLiveListReq) (res *httpserver.CMSQueryResp, err error) {
	return liverecord.GetCMSDailyEffectiveLiveList(ctx, req)
}

// CMSWeeklyUnsettledLiveList CMS分页查询本周未结算流水
func (c *LiveRecordCMSController) CMSWeeklyUnsettledLiveList(ctx context.Context, req *liverecorddto.CMSWeeklyUnsettledLiveListReq) (res *httpserver.CMSQueryResp, err error) {
	return liverecord.GetCMSWeeklyUnsettledLiveList(ctx, req)
}

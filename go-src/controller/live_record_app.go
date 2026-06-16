package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/liverecorddto"
	"xr-game-server/module/liverecord"
)

const (
	LiveRecordAppUrl = "/liveRecord"
)

type LiveRecordAppController struct{}

func initLiveRecordAppController() {
	httpserver.RegAPI(LiveRecordAppUrl, &LiveRecordAppController{})
}

// AppLiveRecordList App端分页查询主播直播数据
func (c *LiveRecordAppController) AppLiveRecordList(ctx context.Context, req *liverecorddto.AppLiveRecordListReq) (res *liverecorddto.AppLiveRecordListRes, err error) {
	return liverecord.GetAppList(ctx, req)
}

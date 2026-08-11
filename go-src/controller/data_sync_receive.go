package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/datasyncdto"
	"xr-game-server/module/datasync"
)

type DataSyncReceiveController struct{}

func initDataSyncReceiveController() {
	httpserver.RegDataSyncReceive(DataSyncCMSUrl, datasync.MiddlewareDataSyncAuth, &DataSyncReceiveController{})
}

func (c *DataSyncReceiveController) ReceiveVipCfg(ctx context.Context, req *datasyncdto.ReceiveVipCfgReq) (*datasyncdto.ReceiveVipCfgRes, error) {
	return datasync.ReceiveVipCfg(ctx, req)
}

package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/datasyncdto"
	"xr-game-server/module/datasync"
)

const DataSyncCMSUrl = "/dataSync"

type DataSyncCMSController struct{}

func initDataSyncCMSController() {
	httpserver.RegCMS(DataSyncCMSUrl, &DataSyncCMSController{})
}

func (c *DataSyncCMSController) GetDataSyncCfg(ctx context.Context, req *datasyncdto.GetDataSyncCfgReq) (*datasyncdto.GetDataSyncCfgRes, error) {
	return datasync.GetDataSyncCfg(ctx, req)
}

func (c *DataSyncCMSController) SaveDataSyncCfg(ctx context.Context, req *datasyncdto.SaveDataSyncCfgReq) (*datasyncdto.SaveDataSyncCfgRes, error) {
	return datasync.SaveDataSyncCfg(ctx, req)
}

func (c *DataSyncCMSController) SyncVipCfg(ctx context.Context, req *datasyncdto.SyncVipCfgReq) (*datasyncdto.SyncVipCfgRes, error) {
	return datasync.SyncVipCfg(ctx, req)
}

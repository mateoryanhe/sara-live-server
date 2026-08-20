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

func (c *DataSyncCMSController) SyncActivityMessage(ctx context.Context, req *datasyncdto.SyncActivityMessageReq) (*datasyncdto.SyncActivityMessageRes, error) {
	return datasync.SyncActivityMessage(ctx, req)
}

func (c *DataSyncCMSController) SyncBanner(ctx context.Context, req *datasyncdto.SyncBannerReq) (*datasyncdto.SyncBatchRes, error) {
	return datasync.SyncBanner(ctx, req)
}

func (c *DataSyncCMSController) SyncGift(ctx context.Context, req *datasyncdto.SyncGiftReq) (*datasyncdto.SyncBatchRes, error) {
	return datasync.SyncGift(ctx, req)
}

func (c *DataSyncCMSController) SyncRechargeCfg(ctx context.Context, req *datasyncdto.SyncRechargeCfgReq) (*datasyncdto.SyncBatchRes, error) {
	return datasync.SyncRechargeCfg(ctx, req)
}

func (c *DataSyncCMSController) SyncFirstRechargeActivityCfg(ctx context.Context, req *datasyncdto.SyncFirstRechargeActivityCfgReq) (*datasyncdto.SyncFirstRechargeActivityCfgRes, error) {
	return datasync.SyncFirstRechargeActivityCfg(ctx, req)
}

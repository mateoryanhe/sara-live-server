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

func (c *DataSyncReceiveController) ReceiveActivityMessage(ctx context.Context, req *datasyncdto.ReceiveActivityMessageReq) (*datasyncdto.ReceiveActivityMessageRes, error) {
	return datasync.ReceiveActivityMessage(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveBanner(ctx context.Context, req *datasyncdto.ReceiveBannerReq) (*datasyncdto.ReceiveBatchRes, error) {
	return datasync.ReceiveBanner(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveGift(ctx context.Context, req *datasyncdto.ReceiveGiftReq) (*datasyncdto.ReceiveBatchRes, error) {
	return datasync.ReceiveGift(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveRechargeCfg(ctx context.Context, req *datasyncdto.ReceiveRechargeCfgReq) (*datasyncdto.ReceiveBatchRes, error) {
	return datasync.ReceiveRechargeCfg(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveFirstRechargeActivityCfg(ctx context.Context, req *datasyncdto.ReceiveFirstRechargeActivityCfgReq) (*datasyncdto.ReceiveFirstRechargeActivityCfgRes, error) {
	return datasync.ReceiveFirstRechargeActivityCfg(ctx, req)
}

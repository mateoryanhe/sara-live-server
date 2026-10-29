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

func (c *DataSyncReceiveController) ReceiveAnchorSalarySocialShareCfg(ctx context.Context, req *datasyncdto.ReceiveAnchorSalarySocialShareCfgReq) (*datasyncdto.ReceiveBatchRes, error) {
	return datasync.ReceiveAnchorSalarySocialShareCfg(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveAnchorNoSalaryShareCfg(ctx context.Context, req *datasyncdto.ReceiveAnchorNoSalaryShareCfgReq) (*datasyncdto.ReceiveBatchRes, error) {
	return datasync.ReceiveAnchorNoSalaryShareCfg(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveAnchorGameShareCfg(ctx context.Context, req *datasyncdto.ReceiveAnchorGameShareCfgReq) (*datasyncdto.ReceiveBatchRes, error) {
	return datasync.ReceiveAnchorGameShareCfg(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveFiatCurrency(ctx context.Context, req *datasyncdto.ReceiveFiatCurrencyReq) (*datasyncdto.ReceiveBatchRes, error) {
	return datasync.ReceiveFiatCurrency(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveFirstRechargeActivityCfg(ctx context.Context, req *datasyncdto.ReceiveFirstRechargeActivityCfgReq) (*datasyncdto.ReceiveFirstRechargeActivityCfgRes, error) {
	return datasync.ReceiveFirstRechargeActivityCfg(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveCfEmailCfg(ctx context.Context, req *datasyncdto.ReceiveCfEmailCfgReq) (*datasyncdto.ReceiveBatchRes, error) {
	return datasync.ReceiveCfEmailCfg(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveSofiePrivacyPolicyCfg(ctx context.Context, req *datasyncdto.ReceiveSofiePrivacyPolicyCfgReq) (*datasyncdto.ReceiveBatchRes, error) {
	return datasync.ReceiveSofiePrivacyPolicyCfg(ctx, req)
}

func (c *DataSyncReceiveController) ReceiveVoiceChatMicSeatTier(ctx context.Context, req *datasyncdto.ReceiveVoiceChatMicSeatTierReq) (*datasyncdto.ReceiveBatchRes, error) {
	return datasync.ReceiveVoiceChatMicSeatTier(ctx, req)
}

package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/incomesettlementdto"
	"xr-game-server/module/incomesettlement"
)

const GuildIncomeSettlementLogCMSUrl = "/guildIncomeSettlementLog"

type GuildIncomeSettlementLogCMSController struct{}

func initGuildIncomeSettlementLogCMSController() {
	httpserver.RegCMS(GuildIncomeSettlementLogCMSUrl, &GuildIncomeSettlementLogCMSController{})
}

func (c *GuildIncomeSettlementLogCMSController) CMSGuildIncomeSettlementLogList(ctx context.Context, req *incomesettlementdto.CMSGuildIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	return incomesettlement.GetGuildCMSList(ctx, req)
}

func (c *GuildIncomeSettlementLogCMSController) CMSGuildIncomeSettlementLogDetail(ctx context.Context, req *incomesettlementdto.CMSGuildIncomeSettlementLogDetailReq) (*incomesettlementdto.CMSGuildIncomeSettlementLogDetailRes, error) {
	return incomesettlement.GetGuildCMSDetail(ctx, req)
}

func (c *GuildIncomeSettlementLogCMSController) CMSBatchApproveGuildSettlement(ctx context.Context, req *incomesettlementdto.CMSBatchApproveGuildSettlementReq) (*incomesettlementdto.CMSBatchApproveGuildSettlementRes, error) {
	return incomesettlement.BatchApproveGuildSettlement(ctx, req)
}

func (c *GuildIncomeSettlementLogCMSController) CMSReopenGuildSettlementApproval(ctx context.Context, req *incomesettlementdto.CMSReopenGuildSettlementApprovalReq) (*incomesettlementdto.CMSReopenGuildSettlementApprovalRes, error) {
	return incomesettlement.ReopenGuildSettlementApproval(ctx, req)
}

func (c *GuildIncomeSettlementLogCMSController) CMSCopyGuildSettlementPayout(ctx context.Context, req *incomesettlementdto.CMSCopyGuildSettlementPayoutReq) (*incomesettlementdto.CMSCopyGuildSettlementPayoutRes, error) {
	return incomesettlement.CopyGuildSettlementPayout(ctx, req)
}

func (c *GuildIncomeSettlementLogCMSController) CMSUpdateGuildSettlementReceivableUsd(ctx context.Context, req *incomesettlementdto.CMSUpdateGuildSettlementReceivableUsdReq) (*incomesettlementdto.CMSUpdateGuildSettlementReceivableUsdRes, error) {
	return incomesettlement.UpdateGuildSettlementReceivableUsd(ctx, req)
}

func (c *GuildIncomeSettlementLogCMSController) CMSBatchTransferGuildSettlement(ctx context.Context, req *incomesettlementdto.CMSBatchTransferGuildSettlementReq) (*incomesettlementdto.CMSBatchTransferGuildSettlementRes, error) {
	return incomesettlement.BatchTransferGuildSettlement(ctx, req)
}

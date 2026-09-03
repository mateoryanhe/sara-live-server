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

func (c *GuildIncomeSettlementLogCMSController) CMSBatchApproveGuildSettlement(ctx context.Context, req *incomesettlementdto.CMSBatchApproveGuildSettlementReq) (*incomesettlementdto.CMSBatchApproveGuildSettlementRes, error) {
	return incomesettlement.BatchApproveGuildSettlement(ctx, req)
}

func (c *GuildIncomeSettlementLogCMSController) CMSBatchTransferGuildSettlement(ctx context.Context, req *incomesettlementdto.CMSBatchTransferGuildSettlementReq) (*incomesettlementdto.CMSBatchTransferGuildSettlementRes, error) {
	return incomesettlement.BatchTransferGuildSettlement(ctx, req)
}

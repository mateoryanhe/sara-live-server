package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/incomesettlementdto"
	"xr-game-server/module/incomesettlement"
)

const AnchorIncomeSettlementLogCMSUrl = "/anchorIncomeSettlementLog"

type AnchorIncomeSettlementLogCMSController struct{}

func initAnchorIncomeSettlementLogCMSController() {
	httpserver.RegCMS(AnchorIncomeSettlementLogCMSUrl, &AnchorIncomeSettlementLogCMSController{})
}

func (c *AnchorIncomeSettlementLogCMSController) CMSAnchorIncomeSettlementLogList(ctx context.Context, req *incomesettlementdto.CMSAnchorIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	return incomesettlement.GetAnchorCMSList(ctx, req)
}

func (c *AnchorIncomeSettlementLogCMSController) CMSBatchApproveAnchorSettlement(ctx context.Context, req *incomesettlementdto.CMSBatchApproveAnchorSettlementReq) (*incomesettlementdto.CMSBatchApproveAnchorSettlementRes, error) {
	return incomesettlement.BatchApproveAnchorSettlement(ctx, req)
}

func (c *AnchorIncomeSettlementLogCMSController) CMSBatchTransferAnchorSettlement(ctx context.Context, req *incomesettlementdto.CMSBatchTransferAnchorSettlementReq) (*incomesettlementdto.CMSBatchTransferAnchorSettlementRes, error) {
	return incomesettlement.BatchTransferAnchorSettlement(ctx, req)
}

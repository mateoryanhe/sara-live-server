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

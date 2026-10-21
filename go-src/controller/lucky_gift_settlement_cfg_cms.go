package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/luckygiftsettlementcfgdto"
	"xr-game-server/module/luckygiftsettlementcfg"
)

const LuckyGiftSettlementCfgCMSUrl = "/luckyGiftSettlementCfg"

type LuckyGiftSettlementCfgCMSController struct{}

func initLuckyGiftSettlementCfgCMSController() {
	httpserver.RegCMS(LuckyGiftSettlementCfgCMSUrl, &LuckyGiftSettlementCfgCMSController{})
}

func (c *LuckyGiftSettlementCfgCMSController) GetLuckyGiftSettlementCfg(ctx context.Context, req *luckygiftsettlementcfgdto.GetLuckyGiftSettlementCfgReq) (*luckygiftsettlementcfgdto.GetLuckyGiftSettlementCfgRes, error) {
	return luckygiftsettlementcfg.GetLuckyGiftSettlementCfg(ctx, req)
}

func (c *LuckyGiftSettlementCfgCMSController) SaveLuckyGiftSettlementCfg(ctx context.Context, req *luckygiftsettlementcfgdto.SaveLuckyGiftSettlementCfgReq) (*luckygiftsettlementcfgdto.SaveLuckyGiftSettlementCfgRes, error) {
	return luckygiftsettlementcfg.SaveLuckyGiftSettlementCfg(ctx, req)
}

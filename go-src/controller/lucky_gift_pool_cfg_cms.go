package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/luckygiftpoolcfgdto"
	"xr-game-server/module/luckygiftpoolcfg"
)

const LuckyGiftPoolCfgCMSUrl = "/luckyGiftPoolCfg"

type LuckyGiftPoolCfgCMSController struct{}

func initLuckyGiftPoolCfgCMSController() {
	httpserver.RegCMS(LuckyGiftPoolCfgCMSUrl, &LuckyGiftPoolCfgCMSController{})
}

func (c *LuckyGiftPoolCfgCMSController) GetLuckyGiftPoolCfg(ctx context.Context, req *luckygiftpoolcfgdto.GetLuckyGiftPoolCfgReq) (*luckygiftpoolcfgdto.GetLuckyGiftPoolCfgRes, error) {
	return luckygiftpoolcfg.GetLuckyGiftPoolCfg(ctx, req)
}

func (c *LuckyGiftPoolCfgCMSController) SaveLuckyGiftPoolCfg(ctx context.Context, req *luckygiftpoolcfgdto.SaveLuckyGiftPoolCfgReq) (*luckygiftpoolcfgdto.SaveLuckyGiftPoolCfgRes, error) {
	return luckygiftpoolcfg.SaveLuckyGiftPoolCfg(ctx, req)
}

func (c *LuckyGiftPoolCfgCMSController) ListPendingBalance(ctx context.Context, req *luckygiftpoolcfgdto.ListLuckyGiftPoolPendingBalanceReq) (*httpserver.CMSQueryResp, error) {
	return luckygiftpoolcfg.ListLuckyGiftPoolPendingBalance(ctx, req)
}

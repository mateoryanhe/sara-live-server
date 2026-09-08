package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/golddto"
	"xr-game-server/module/wallet"
)

type GoldAppController struct{}

func initGoldAppController() {
	httpserver.RegAPI(GoldUrl, &GoldAppController{})
}

// ExchangeGoldToDiamond App端金币兑换钻石(1金币=100钻石)
func (c *GoldAppController) ExchangeGoldToDiamond(ctx context.Context, req *golddto.AppExchangeGoldToDiamondReq) (*golddto.AppExchangeGoldToDiamondRes, error) {
	return wallet.AppExchangeGoldToDiamond(ctx, req)
}

// TransferGold App端转赠金币给指定用户(币商扣款,目标用户到账)
func (c *GoldAppController) TransferGold(ctx context.Context, req *golddto.AppTransferGoldReq) (*golddto.AppTransferGoldRes, error) {
	return wallet.AppTransferGold(ctx, req)
}

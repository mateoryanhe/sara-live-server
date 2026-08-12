package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/walletdto"
	"xr-game-server/module/wallet"
)

const WalletCMSUrl = "/wallet"

type WalletCMSController struct{}

func initWalletCMSController() {
	httpserver.RegCMS(WalletCMSUrl, &WalletCMSController{})
}

func (c *WalletCMSController) GetWalletExchangeCfg(ctx context.Context, req *walletdto.GetWalletExchangeCfgReq) (*walletdto.GetWalletExchangeCfgRes, error) {
	return wallet.GetWalletExchangeCfg(ctx, req)
}

func (c *WalletCMSController) SaveWalletExchangeCfg(ctx context.Context, req *walletdto.SaveWalletExchangeCfgReq) (*walletdto.SaveWalletExchangeCfgRes, error) {
	return wallet.SaveWalletExchangeCfg(ctx, req)
}

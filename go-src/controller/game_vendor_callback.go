package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/module/game"
)

type GameVendorCallbackController struct{}

func initGameVendorCallbackController() {
	// 第三方固定访问: POST https://{domain}/verify, POST https://{domain}/balance, POST https://{domain}/transfer
	httpserver.RegRootNonAuthCustomizeRes(&GameVendorCallbackController{})
}

func (c *GameVendorCallbackController) Verify(ctx context.Context, req *gameplatformdto.VendorVerifyReq) (*gameplatformdto.VendorVerifyRes, error) {
	return game.HandleVendorVerify(ctx, req)
}

func (c *GameVendorCallbackController) Balance(ctx context.Context, req *gameplatformdto.VendorBalanceReq) (*gameplatformdto.VendorBalanceRes, error) {
	return game.HandleVendorBalance(ctx, req)
}

func (c *GameVendorCallbackController) Transfer(ctx context.Context, req *gameplatformdto.VendorTransferReq) (*gameplatformdto.VendorTransferRes, error) {
	return game.HandleVendorTransfer(ctx, req)
}

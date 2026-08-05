package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/module/game"
)

const GameVendorCallbackUrl = "/game"

type GameVendorCallbackController struct{}

func initGameVendorCallbackController() {
	httpserver.RegCMSNonAuthCustomizeRes(GameVendorCallbackUrl, &GameVendorCallbackController{})
}

func (c *GameVendorCallbackController) Verify(ctx context.Context, req *gameplatformdto.VendorVerifyReq) (*gameplatformdto.VendorVerifyRes, error) {
	return game.HandleVendorVerify(ctx, req)
}

func (c *GameVendorCallbackController) Balance(ctx context.Context, req *gameplatformdto.VendorBalanceReq) (*gameplatformdto.VendorBalanceRes, error) {
	return game.HandleVendorBalance(ctx, req)
}

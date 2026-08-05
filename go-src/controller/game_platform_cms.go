package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/module/game"
)

const GamePlatformCMSUrl = "/gamePlatform"

type GamePlatformCMSController struct{}

func initGamePlatformCMSController() {
	httpserver.RegCMS(GamePlatformCMSUrl, &GamePlatformCMSController{})
}

func (c *GamePlatformCMSController) GetGamePlatformCfg(ctx context.Context, req *gameplatformdto.GetGamePlatformCfgReq) (*gameplatformdto.GetGamePlatformCfgRes, error) {
	return game.GetGamePlatformCfg(ctx, req)
}

func (c *GamePlatformCMSController) SaveGamePlatformCfg(ctx context.Context, req *gameplatformdto.SaveGamePlatformCfgReq) (*gameplatformdto.SaveGamePlatformCfgRes, error) {
	return game.SaveGamePlatformCfg(ctx, req)
}

func (c *GamePlatformCMSController) VendorGameList(ctx context.Context, req *gameplatformdto.VendorGameListReq) (*httpserver.CMSQueryResp, error) {
	return game.GetVendorGameList(ctx, req)
}

func (c *GamePlatformCMSController) ReloadVendorGameCache(ctx context.Context, req *gameplatformdto.ReloadVendorGameCacheReq) (*gameplatformdto.ReloadVendorGameCacheRes, error) {
	return game.ReloadVendorGameCacheCMS(ctx, req)
}

func (c *GamePlatformCMSController) GameShelfList(ctx context.Context, req *gameplatformdto.GameShelfListReq) (*httpserver.CMSQueryResp, error) {
	return game.GetGameShelfList(ctx, req)
}

func (c *GamePlatformCMSController) AddGameShelf(ctx context.Context, req *gameplatformdto.AddGameShelfReq) (*gameplatformdto.AddGameShelfRes, error) {
	return game.AddGameShelf(ctx, req)
}

func (c *GamePlatformCMSController) DeleteGameShelf(ctx context.Context, req *gameplatformdto.DeleteGameShelfReq) (*gameplatformdto.DeleteGameShelfRes, error) {
	return game.DeleteGameShelf(ctx, req)
}

func (c *GamePlatformCMSController) BatchAddGameShelf(ctx context.Context, req *gameplatformdto.BatchAddGameShelfReq) (*gameplatformdto.BatchAddGameShelfRes, error) {
	return game.BatchAddGameShelf(ctx, req)
}

func (c *GamePlatformCMSController) BatchDeleteGameShelf(ctx context.Context, req *gameplatformdto.BatchDeleteGameShelfReq) (*gameplatformdto.BatchDeleteGameShelfRes, error) {
	return game.BatchDeleteGameShelf(ctx, req)
}

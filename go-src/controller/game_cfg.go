package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gamecfgdto"
	"xr-game-server/module/game"
)

const GameCfgUrl = "/gameCfg"

type GameCfgController struct{}

func initGameCfgController() {
	httpserver.RegCMS(GameCfgUrl, &GameCfgController{})
}

func (c *GameCfgController) GameCfgList(ctx context.Context, req *gamecfgdto.GameCfgListReq) (*httpserver.CMSQueryResp, error) {
	return game.GetList(ctx, req)
}

func (c *GameCfgController) CreateGameCfg(ctx context.Context, req *gamecfgdto.CreateGameCfgReq) (*gamecfgdto.CreateGameCfgRes, error) {
	return game.Create(ctx, req)
}

func (c *GameCfgController) UpdateGameCfg(ctx context.Context, req *gamecfgdto.UpdateGameCfgReq) (*gamecfgdto.UpdateGameCfgRes, error) {
	return game.Update(ctx, req)
}

func (c *GameCfgController) DeleteGameCfg(ctx context.Context, req *gamecfgdto.DeleteGameCfgReq) (*gamecfgdto.DeleteGameCfgRes, error) {
	return game.Delete(ctx, req)
}

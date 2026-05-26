package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gamecfgdto"
	"xr-game-server/module/game"
)

const GameCfgAppUrl = "/gameCfg"

type GameCfgAppController struct{}

func initGameCfgAppController() {
	httpserver.RegAPI(GameCfgAppUrl, &GameCfgAppController{})
}

// AppGameCfgList App端分页查询游戏列表(仅已上架,走缓存)
func (c *GameCfgAppController) AppGameCfgList(ctx context.Context, req *gamecfgdto.AppGameCfgListReq) (*gamecfgdto.AppGameCfgListRes, error) {
	return game.GetAppGameCfgList(ctx, req)
}

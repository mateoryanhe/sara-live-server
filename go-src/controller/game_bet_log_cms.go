package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gamebetdto"
	"xr-game-server/module/gamebet"
)

const GameBetLogCMSUrl = "/gameBetLog"

type GameBetLogCMSController struct{}

func initGameBetLogCMSController() {
	httpserver.RegCMS(GameBetLogCMSUrl, &GameBetLogCMSController{})
}

// CMSGameBetLogList CMS 分页查询游戏下注记录
func (c *GameBetLogCMSController) CMSGameBetLogList(ctx context.Context, req *gamebetdto.CMSGameBetLogListReq) (*httpserver.CMSQueryResp, error) {
	return gamebet.GetCMSList(ctx, req)
}

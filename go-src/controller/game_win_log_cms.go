package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gamewindto"
	"xr-game-server/module/gamewin"
)

const GameWinLogCMSUrl = "/gameWinLog"

type GameWinLogCMSController struct{}

func initGameWinLogCMSController() {
	httpserver.RegCMS(GameWinLogCMSUrl, &GameWinLogCMSController{})
}

// CMSGameWinLogList CMS 分页查询游戏派彩记录
func (c *GameWinLogCMSController) CMSGameWinLogList(ctx context.Context, req *gamewindto.CMSGameWinLogListReq) (*httpserver.CMSQueryResp, error) {
	return gamewin.GetCMSList(ctx, req)
}

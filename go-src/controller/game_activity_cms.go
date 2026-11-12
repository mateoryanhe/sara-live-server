package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gameactivitydto"
	"xr-game-server/module/gameactivity"
)

const GameActivityCMSUrl = "/gameActivity"

type GameActivityCMSController struct{}

func initGameActivityCMSController() {
	httpserver.RegCMS(GameActivityCMSUrl, &GameActivityCMSController{})
}

func (c *GameActivityCMSController) GetGameActivityCfg(ctx context.Context, req *gameactivitydto.GetGameActivityCfgReq) (*gameactivitydto.GetGameActivityCfgRes, error) {
	return gameactivity.GetGameActivityCfg(ctx, req)
}

func (c *GameActivityCMSController) SaveGameActivityCfg(ctx context.Context, req *gameactivitydto.SaveGameActivityCfgReq) (*gameactivitydto.SaveGameActivityCfgRes, error) {
	return gameactivity.SaveGameActivityCfg(ctx, req)
}

package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gameconsumrankdto"
	"xr-game-server/module/gameconsumrank"
)

const (
	GameConsumeRankAppUrl = "/gameConsumeRank"
)

type GameConsumeRankAppController struct{}

func initGameConsumeRankAppController() {
	httpserver.RegAPI(GameConsumeRankAppUrl, &GameConsumeRankAppController{})
}

// AppGameConsumeRankList App端查询游戏消费榜
func (c *GameConsumeRankAppController) AppGameConsumeRankList(ctx context.Context, req *gameconsumrankdto.AppGameConsumeRankListReq) (res *gameconsumrankdto.AppGameConsumeRankListRes, err error) {
	return gameconsumrank.GetAppGameConsumeRankList(ctx, req)
}

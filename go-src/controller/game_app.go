package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/gamebetdto"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/module/game"
)

const GameAppUrl = "/game"

type GameAppController struct{}

func initGameAppController() {
	httpserver.RegAPI(GameAppUrl, &GameAppController{})
}

// AppGameList App 分页查询已上架游戏列表
func (c *GameAppController) AppGameList(ctx context.Context, req *gameplatformdto.AppGameListReq) (*gameplatformdto.AppGameListRes, error) {
	return game.GetAppGameList(ctx, req)
}

// AppGameStart App 获取游戏启动链接
func (c *GameAppController) AppGameStart(ctx context.Context, req *gameplatformdto.AppGameStartReq) (*gameplatformdto.AppGameStartRes, error) {
	return game.GetAppGameStartLink(ctx, req)
}

// AppGameBetList App 分页查询游戏下注记录
func (c *GameAppController) AppGameBetList(ctx context.Context, req *gamebetdto.AppGameBetListReq) (*gamebetdto.AppGameBetListRes, error) {
	return game.GetAppBetList(ctx, req)
}

// AppGameConsumeRank App 分页查询单场直播游戏消费榜
func (c *GameAppController) AppGameConsumeRank(ctx context.Context, req *gamebetdto.AppGameConsumeRankReq) (*gamebetdto.AppGameConsumeRankRes, error) {
	return game.GetAppGameConsumeRank(ctx, req)
}

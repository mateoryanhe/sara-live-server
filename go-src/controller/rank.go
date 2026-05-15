package controller

import (
	"context"
	"time"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/rankdto"
	"xr-game-server/gameevent"
	"xr-game-server/module/rank"
)

type RankApp struct {
}

type RankCMS struct {
}

func initRank() {
	httpserver.RegAPI("/rank", &RankApp{})
	httpserver.RegCMS("/rank", &RankCMS{})
}

func (receiver *RankApp) GetRank(context context.Context, req *rankdto.GetRankReq) (res *rankdto.GetRankRes, err error) {
	return rank.GetRank(context, req)
}
func (receiver *RankApp) GetRankVersion(context context.Context, req *rankdto.GetRankVersionReq) (int64, error) {
	return rank.GetRankVersionReq(context, req)
}
func (receiver *RankCMS) AddRankDataReq(ctx context.Context, req *rankdto.AddRankDataReq) (int64, error) {
	event.Pub(gameevent.AddRankValEvent, gameevent.NewAddRankValEvent(httpserver.GetAuthId(ctx), req.Val, req.TypeId))
	return time.Now().Unix(), nil
}

func (receiver *RankCMS) ReduceRankDataReq(ctx context.Context, req *rankdto.ReduceRankReq) (bool, error) {
	event.Pub(gameevent.ReduceRankValEvent, gameevent.NewReduceRankEvent(httpserver.GetAuthId(ctx), req.Val, req.TypeId))
	return false, nil
}
func (receiver *RankCMS) UpRank(ctx context.Context, req *rankdto.UpRankReq) (bool, error) {
	return rank.UpRankReq(ctx, req)
}
func (receiver *RankCMS) DownRank(ctx context.Context, req *rankdto.DownRankReq) (bool, error) {
	return rank.DownRankReq(ctx, req)
}

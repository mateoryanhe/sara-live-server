package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/richrankdto"
	"xr-game-server/module/richrank"
)

const (
	RichRankAppUrl = "/richRank"
)

type RichRankAppController struct{}

func initRichRankAppController() {
	httpserver.RegAPI(RichRankAppUrl, &RichRankAppController{})
}

// AppRichRankList App端查询富豪榜
func (c *RichRankAppController) AppRichRankList(ctx context.Context, req *richrankdto.AppRichRankListReq) (res *richrankdto.AppRichRankListRes, err error) {
	return richrank.GetAppRichRankList(ctx, req)
}

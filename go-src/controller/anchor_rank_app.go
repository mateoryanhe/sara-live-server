package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/anchorrankdto"
	"xr-game-server/module/anchorrank"
)

const (
	AnchorRankAppUrl = "/anchorRank"
)

type AnchorRankAppController struct{}

func initAnchorRankAppController() {
	httpserver.RegAPI(AnchorRankAppUrl, &AnchorRankAppController{})
}

// AppAnchorRankList App端查询主播红人榜
func (c *AnchorRankAppController) AppAnchorRankList(ctx context.Context, req *anchorrankdto.AppAnchorRankListReq) (res *anchorrankdto.AppAnchorRankListRes, err error) {
	return anchorrank.GetAppAnchorRankList(ctx, req)
}

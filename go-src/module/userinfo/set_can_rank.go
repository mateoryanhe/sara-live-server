package userinfo

import (
	"context"
	"xr-game-server/core/event"
	"xr-game-server/dao/userextdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/gameevent"
)

// SetCanRank CMS 设置用户是否可上排行榜
func SetCanRank(_ context.Context, req *accountdto.SetCanRankReq) (*accountdto.SetCanRankRes, error) {
	ext := userextdao.GetByUserId(req.AccountId)
	if ext.CanRank == req.CanRank {
		return &accountdto.SetCanRankRes{Success: true}, nil
	}
	ext.SetCanRank(req.CanRank)
	AddIdToCache(req.AccountId)
	event.Pub(gameevent.RankListRefreshEvent, nil)
	return &accountdto.SetCanRankRes{Success: true}, nil
}

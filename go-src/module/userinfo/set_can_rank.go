package userinfo

import (
	"context"
	"xr-game-server/core/event"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/gameevent"
)

// SetCanRank CMS 设置用户是否可上排行榜
func SetCanRank(_ context.Context, req *accountdto.SetCanRankReq) (*accountdto.SetCanRankRes, error) {
	ext := userinfodao.GetUserExtByUserId(req.AccountId)
	if ext.CanRank == req.CanRank {
		return &accountdto.SetCanRankRes{Success: true}, nil
	}
	ext.SetCanRank(req.CanRank)
	userinfodao.PublishUserExt(ext)
	event.Pub(gameevent.RankListRefreshEvent, nil)
	return &accountdto.SetCanRankRes{Success: true}, nil
}

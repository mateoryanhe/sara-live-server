package account

import (
	"context"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/accountdto"
)

func CancelUser(ctx context.Context, req *accountdto.CancelReq) (bool, error) {
	accountdao.CancelAccount(req.AccountId)
	return true, nil
}

func UnCancelUser(ctx context.Context, req *accountdto.UnCancelReq) (bool, error) {
	accountdao.UnCancelAccount(req.AccountId)
	return true, nil
}

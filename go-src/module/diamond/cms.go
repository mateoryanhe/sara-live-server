package diamond

import (
	"context"
	"xr-game-server/constants/currency"
	"xr-game-server/dto/diamonddto"
)

// CMSAdd CMS 后台增加用户钻石
func CMSAdd(_ context.Context, req *diamonddto.CMSAddDiamondReq) (*diamonddto.CMSAddDiamondRes, error) {
	after, err := Add(req.UserId, req.Amount, currency.ReasonGmAdjust)
	if err != nil {
		return nil, err
	}
	return &diamonddto.CMSAddDiamondRes{Diamond: after}, nil
}

// CMSSub CMS 后台扣减用户钻石
func CMSSub(_ context.Context, req *diamonddto.CMSSubDiamondReq) (*diamonddto.CMSSubDiamondRes, error) {
	after, err := Sub(req.UserId, req.Amount, currency.ReasonGmAdjust)
	if err != nil {
		return nil, err
	}
	return &diamonddto.CMSSubDiamondRes{Diamond: after}, nil
}

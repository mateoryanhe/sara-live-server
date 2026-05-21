package gold

import (
	"context"
	"xr-game-server/constants/currency"
	"xr-game-server/dto/golddto"
)

// CMSAdd CMS 后台增加用户金币
func CMSAdd(_ context.Context, req *golddto.CMSAddGoldReq) (*golddto.CMSAddGoldRes, error) {
	after, err := Add(req.UserId, req.Amount, currency.ReasonGmAdjust)
	if err != nil {
		return nil, err
	}
	return &golddto.CMSAddGoldRes{Gold: after}, nil
}

// CMSSub CMS 后台扣减用户金币
func CMSSub(_ context.Context, req *golddto.CMSSubGoldReq) (*golddto.CMSSubGoldRes, error) {
	after, err := Sub(req.UserId, req.Amount, currency.ReasonGmAdjust)
	if err != nil {
		return nil, err
	}
	return &golddto.CMSSubGoldRes{Gold: after}, nil
}

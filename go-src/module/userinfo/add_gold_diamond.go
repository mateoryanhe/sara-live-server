package userinfo

import (
	"context"
	"xr-game-server/constants/currency"
	"xr-game-server/dto/diamonddto"
	"xr-game-server/dto/golddto"
	"xr-game-server/module/wallet"
)

// DiamondCMSAdd CMS 后台增加用户钻石
func DiamondCMSAdd(_ context.Context, req *diamonddto.CMSAddDiamondReq) (*diamonddto.CMSAddDiamondRes, error) {
	after, err := wallet.DiamondAdd(req.UserId, req.Amount, currency.ReasonGmAdjust)
	if err != nil {
		return nil, err
	}
	AddIdToCache(req.UserId)
	return &diamonddto.CMSAddDiamondRes{Diamond: after}, nil
}

// DiamondCMSSub CMS 后台扣减用户钻石
func DiamondCMSSub(_ context.Context, req *diamonddto.CMSSubDiamondReq) (*diamonddto.CMSSubDiamondRes, error) {
	after, err := wallet.DiamondSub(req.UserId, req.Amount, currency.ReasonGmAdjust)
	if err != nil {
		return nil, err
	}
	AddIdToCache(req.UserId)
	return &diamonddto.CMSSubDiamondRes{Diamond: after}, nil
}

// GoldCMSAdd CMS 后台增加用户金币
func GoldCMSAdd(_ context.Context, req *golddto.CMSAddGoldReq) (*golddto.CMSAddGoldRes, error) {
	after, err := wallet.GoldAdd(req.UserId, req.Amount, currency.ReasonGmAdjust)
	if err != nil {
		return nil, err
	}
	AddIdToCache(req.UserId)
	return &golddto.CMSAddGoldRes{Gold: after}, nil
}

// GoldCMSSub CMS 后台扣减用户金币
func GoldCMSSub(_ context.Context, req *golddto.CMSSubGoldReq) (*golddto.CMSSubGoldRes, error) {
	after, err := wallet.GoldSub(req.UserId, req.Amount, currency.ReasonGmAdjust)
	if err != nil {
		return nil, err
	}
	AddIdToCache(req.UserId)
	return &golddto.CMSSubGoldRes{Gold: after}, nil
}

package wallet

import (
	"context"

	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/golddto"
	"xr-game-server/errercode"
)

type appExchangeResult struct {
	goldAfter          float64
	diamondAfter       float64
	goldAmount         float64
	goldDeduct         float64
	feeDiamond         float64
	diamondGross       float64
	diamondAmount      float64
	goldToDiamondRate  int
	exchangeFeePercent float64
}

func exchangeGoldToDiamond(userId uint64, goldAmount float64, applyAppFee bool) (*appExchangeResult, error) {
	if goldAmount <= 0 {
		return nil, errercode.CreateCode(errercode.GoldAmountInvalid)
	}

	snap := GetExchangeCfgSnapshot()
	goldDeduct := goldAmount
	diamondGross, feeDiamond, diamondAmount := calcExchangeDiamond(goldAmount, snap, applyAppFee)

	goldAfter, err := GoldSub(userId, goldDeduct, currency.ReasonGoldExchange)
	if err != nil {
		return nil, err
	}
	diamondAfter, err := DiamondAdd(userId, diamondAmount, currency.ReasonGoldExchange)
	if err != nil {
		if _, refundErr := GoldAdd(userId, goldDeduct, currency.ReasonRefund); refundErr != nil {
			return nil, err
		}
		return nil, err
	}
	return &appExchangeResult{
		goldAfter:          goldAfter,
		diamondAfter:       diamondAfter,
		goldAmount:         goldAmount,
		goldDeduct:         goldDeduct,
		feeDiamond:         feeDiamond,
		diamondGross:       diamondGross,
		diamondAmount:      diamondAmount,
		goldToDiamondRate:  snap.GoldToDiamondRate,
		exchangeFeePercent: snap.ExchangeFeePercent,
	}, nil
}

// AppExchangeGoldToDiamond App端金币兑换钻石;手续费比例>0 时从兑换钻石中扣除
func AppExchangeGoldToDiamond(ctx context.Context, req *golddto.AppExchangeGoldToDiamondReq) (*golddto.AppExchangeGoldToDiamondRes, error) {
	userId := httpserver.GetAuthId(ctx)
	result, err := exchangeGoldToDiamond(userId, req.GoldAmount, true)
	if err != nil {
		return nil, err
	}
	return &golddto.AppExchangeGoldToDiamondRes{
		Gold:               result.goldAfter,
		Diamond:            result.diamondAfter,
		ExchangedGold:      result.goldAmount,
		ExchangedDiamond:   result.diamondAmount,
		GoldDeduct:         result.goldDeduct,
		FeeDiamond:         result.feeDiamond,
		DiamondGross:       result.diamondGross,
		GoldToDiamondRate:  result.goldToDiamondRate,
		ExchangeFeePercent: result.exchangeFeePercent,
	}, nil
}

// exchangeGoldToDiamondForAutoPay 业务内自动兑换(免手续费)
func exchangeGoldToDiamondForAutoPay(userId uint64, goldAmount float64) (goldAfter, diamondAfter float64, err error) {
	result, err := exchangeGoldToDiamond(userId, goldAmount, false)
	if err != nil {
		return 0, 0, err
	}
	return result.goldAfter, result.diamondAfter, nil
}

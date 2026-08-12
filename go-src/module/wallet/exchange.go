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
	feeGold            float64
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
	feeGold := 0.0
	if applyAppFee {
		goldDeduct = calcAppExchangeGoldDeduct(goldAmount, snap)
		feeGold = goldDeduct - goldAmount
		if feeGold < 0 {
			feeGold = 0
		}
	}
	diamondAmount := goldAmount * float64(snap.GoldToDiamondRate)

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
		feeGold:            feeGold,
		diamondAmount:      diamondAmount,
		goldToDiamondRate:  snap.GoldToDiamondRate,
		exchangeFeePercent: snap.ExchangeFeePercent,
	}, nil
}

// AppExchangeGoldToDiamond App端金币兑换钻石;手续费比例>0 时收取手续费
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
		FeeGold:            result.feeGold,
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

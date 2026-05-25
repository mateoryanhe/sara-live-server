package wallet

import (
	"context"
	"xr-game-server/constants/currency"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/golddto"
	"xr-game-server/errercode"
)

const goldToDiamondRate = 100

// AppExchangeGoldToDiamond App端金币兑换钻石,兑换比例 1 金币 = 100 钻石
func AppExchangeGoldToDiamond(ctx context.Context, req *golddto.AppExchangeGoldToDiamondReq) (*golddto.AppExchangeGoldToDiamondRes, error) {
	userId := httpserver.GetAuthId(ctx)
	goldAmount := req.GoldAmount
	if goldAmount <= 0 {
		return nil, errercode.CreateCode(errercode.GoldAmountInvalid)
	}
	diamondAmount := goldAmount * goldToDiamondRate

	goldAfter, err := GoldSub(userId, goldAmount, currency.ReasonGoldExchange)
	if err != nil {
		return nil, err
	}
	diamondAfter, err := DiamondAdd(userId, diamondAmount, currency.ReasonGoldExchange)
	if err != nil {
		if _, refundErr := GoldAdd(userId, goldAmount, currency.ReasonRefund); refundErr != nil {
			return nil, err
		}
		return nil, err
	}

	return &golddto.AppExchangeGoldToDiamondRes{
		Gold:             goldAfter,
		Diamond:          diamondAfter,
		ExchangedGold:    goldAmount,
		ExchangedDiamond: diamondAmount,
	}, nil
}

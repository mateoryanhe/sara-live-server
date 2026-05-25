package golddto

import "github.com/gogf/gf/v2/frame/g"

// AppExchangeGoldToDiamondReq App端金币兑换钻石(1金币=100钻石)
type AppExchangeGoldToDiamondReq struct {
	g.Meta     `path:"/exchangeGoldToDiamond" method:"post" summary:"金币兑换钻石" tags:"金币"`
	GoldAmount float64 `json:"goldAmount" v:"required|min:0.0001#金币数量不能为空|金币数量必须大于0" dc:"兑换金币数量(正数)"`
}

type AppExchangeGoldToDiamondRes struct {
	Gold             float64 `json:"gold" dc:"兑换后金币余额"`
	Diamond          float64 `json:"diamond" dc:"兑换后钻石余额"`
	ExchangedGold    float64 `json:"exchangedGold" dc:"本次兑换金币数量"`
	ExchangedDiamond float64 `json:"exchangedDiamond" dc:"本次获得钻石数量"`
}

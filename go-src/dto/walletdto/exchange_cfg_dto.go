package walletdto

import "github.com/gogf/gf/v2/frame/g"

type GetWalletExchangeCfgReq struct {
	g.Meta `path:"/getWalletExchangeCfg" method:"post" summary:"查询金币兑换钻石配置" tags:"钱包配置"`
}

type WalletExchangeCfgItem struct {
	ID                 string  `json:"id"`
	GoldToDiamondRate  int     `json:"goldToDiamondRate"`
	ExchangeFeePercent float64 `json:"exchangeFeePercent"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

type GetWalletExchangeCfgRes struct {
	Cfg *WalletExchangeCfgItem `json:"cfg"`
}

type SaveWalletExchangeCfgReq struct {
	g.Meta             `path:"/saveWalletExchangeCfg" method:"post" summary:"保存金币兑换钻石配置" tags:"钱包配置"`
	ID                 uint64  `json:"id,string" dc:"配置ID,更新时传"`
	GoldToDiamondRate  int     `json:"goldToDiamondRate" v:"required|min:1#兑换比例不能为空|兑换比例必须大于0" dc:"1金币兑换钻石数"`
	ExchangeFeePercent float64 `json:"exchangeFeePercent" v:"min:0#手续费不能为负数" dc:"App手动兑换手续费(%)，大于0时收取"`
}

type SaveWalletExchangeCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

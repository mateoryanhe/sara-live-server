package fiatcurrencydto

import "github.com/gogf/gf/v2/frame/g"

// AppFiatCurrencyListReq App 查询币种列表(无需鉴权,仅返回已启用)
type AppFiatCurrencyListReq struct {
	g.Meta     `path:"/fiatCurrencyListForApp" method:"post" summary:"App查询币种列表" tags:"币种配置"`
	TypeFilter int `json:"typeFilter" v:"in:0,1,2#币种类型无效" dc:"币种类型(0全部,1法币,2加密币)"`
}

// AppFiatCurrencyItem App 币种项
type AppFiatCurrencyItem struct {
	CurrencyCode string `json:"currencyCode"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Icon         string `json:"icon"`
	CurrencyType uint8  `json:"currencyType" dc:"1法币,2加密币"`
	Sort         int    `json:"sort"`
}

type AppFiatCurrencyListRes struct {
	List []*AppFiatCurrencyItem `json:"list"`
}

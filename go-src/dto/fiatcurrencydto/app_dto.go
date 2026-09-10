package fiatcurrencydto

import "github.com/gogf/gf/v2/frame/g"

// AppFiatCurrencyListReq App 查询区域列表(硬编码 HaiPay region,无需鉴权)
type AppFiatCurrencyListReq struct {
	g.Meta     `path:"/fiatCurrencyListForApp" method:"post" summary:"App查询支付区域列表" tags:"币种配置"`
	TypeFilter int `json:"typeFilter" v:"in:0,1,2#币种类型无效" dc:"兼容旧字段(2=加密返回空,其它返回区域列表)"`
}

// AppFiatCurrencyItem App 区域项(currencyCode=HaiPay region,如 ID/PH)
type AppFiatCurrencyItem struct {
	CurrencyCode string `json:"currencyCode" dc:"HaiPay region,建单原样上报"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Icon         string `json:"icon"`
	CurrencyType uint8  `json:"currencyType" dc:"固定1"`
	Sort         int    `json:"sort"`
}

type AppFiatCurrencyListRes struct {
	List []*AppFiatCurrencyItem `json:"list"`
}

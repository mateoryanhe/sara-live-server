package fiatcurrencydto

import "github.com/gogf/gf/v2/frame/g"

// AppFiatCurrencyListReq App 查询支付区域列表(无需鉴权)
type AppFiatCurrencyListReq struct {
	g.Meta     `path:"/fiatCurrencyListForApp" method:"post" summary:"App查询HaiPay支付区域列表" tags:"币种配置"`
	TypeFilter int `json:"typeFilter" v:"in:0,1,2#币种类型无效" dc:"兼容旧字段(2=加密返回空,其它返回区域列表)"`
}

// HaiPayRegionListReq CMS 查询 HaiPay 支付区域列表
type HaiPayRegionListReq struct {
	g.Meta `path:"/haiPayRegionList" method:"post" summary:"CMS查询HaiPay支付区域列表" tags:"币种配置"`
}

// AppFiatCurrencyItem 支付区域项(currencyCode=HaiPay region 简码,与 country.Code 一致)
type AppFiatCurrencyItem struct {
	CurrencyCode string `json:"currencyCode" dc:"HaiPay region 简码,建单原样上报"`
	Name         string `json:"name" dc:"英文名(兼容旧字段,同 nameEn)"`
	NameEn       string `json:"nameEn" dc:"英文名称"`
	NameZh       string `json:"nameZh" dc:"中文名称"`
	Symbol       string `json:"symbol" dc:"展示符号,同简码"`
	Icon         string `json:"icon" dc:"国旗URL(country-flags/{version}/{code}.png)"`
	CurrencyType uint8  `json:"currencyType" dc:"固定1"`
	Sort         int    `json:"sort"`
}

type AppFiatCurrencyListRes struct {
	List []*AppFiatCurrencyItem `json:"list"`
}

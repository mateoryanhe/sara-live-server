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

// CoinMerchantPaymentRegionListReq 币商 App 查询专用支付区域列表。
type CoinMerchantPaymentRegionListReq struct {
	g.Meta `path:"/paymentRegionList" method:"post" summary:"币商App查询HaiPay支付区域列表" tags:"币商充值"`
}

type AppHaiPayPaymentMethod struct {
	PayType     string `json:"payType" dc:"HaiPay本地代收支付类型"`
	InBankCode  string `json:"inBankCode" dc:"HaiPay本地代收支付编码"`
	MinAmount   string `json:"minAmount" dc:"该币种单笔最小金额"`
	MaxAmount   string `json:"maxAmount" dc:"该币种单笔最大金额"`
	Description string `json:"description" dc:"支付方式说明"`
}

// AppFiatCurrencyItem 支付区域项(currencyCode=HaiPay region 简码,与 country.Code 一致)。
// App 区域列表只返回地区展示字段；CMS 查询可额外返回币种与支付方式配置。
type AppFiatCurrencyItem struct {
	CurrencyCode      string                    `json:"currencyCode" dc:"HaiPay region 简码,建单原样上报"`
	FiatCurrencyCode  string                    `json:"fiatCurrencyCode,omitempty" dc:"CMS查询字段;服务端配置的实际下单币种"`
	FiatCurrencyCodes []string                  `json:"fiatCurrencyCodes,omitempty" dc:"CMS查询字段;国家/地区支持的全部ISO 4217法币代码"`
	Continent         string                    `json:"continent" dc:"HaiPay地区分组"`
	PayType           string                    `json:"payType,omitempty" dc:"CMS查询字段;配置的支付类型,逗号分隔"`
	InBankCode        string                    `json:"inBankCode,omitempty" dc:"CMS查询字段;配置的支付编码,逗号分隔"`
	PayTypes          []string                  `json:"payTypes,omitempty" dc:"CMS查询字段;配置的支付类型列表"`
	InBankCodes       []string                  `json:"inBankCodes,omitempty" dc:"CMS查询字段;配置的支付编码列表"`
	PaymentMethods    []*AppHaiPayPaymentMethod `json:"paymentMethods,omitempty" dc:"CMS查询字段;启用的本地代收支付类型与编码组合"`
	Name              string                    `json:"name" dc:"英文名(兼容旧字段,同 nameEn)"`
	NameEn            string                    `json:"nameEn" dc:"英文名称"`
	NameZh            string                    `json:"nameZh" dc:"中文名称"`
	Symbol            string                    `json:"symbol" dc:"展示符号,同简码"`
	Icon              string                    `json:"icon" dc:"国旗URL(country-flags/{version}/{code}.png)"`
	CurrencyType      uint8                     `json:"currencyType" dc:"固定1"`
	Sort              int                       `json:"sort"`
}

type AppFiatCurrencyListRes struct {
	List []*AppFiatCurrencyItem `json:"list"`
}

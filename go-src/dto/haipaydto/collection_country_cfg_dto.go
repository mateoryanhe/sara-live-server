package haipaydto

import "github.com/gogf/gf/v2/frame/g"

type GetCollectionCountryCfgReq struct {
	g.Meta `path:"/getCollectionCountryCfg" method:"post" summary:"查询HaiPay国家支付配置" tags:"HaiPay支付配置"`
}

type CollectionCountryCfgItem struct {
	ID                     string                              `json:"id"`
	CountryCode            string                              `json:"countryCode"`
	CountryNameEn          string                              `json:"countryNameEn"`
	CountryNameZh          string                              `json:"countryNameZh"`
	Continent              string                              `json:"continent"`
	SupportedCurrencies    []string                            `json:"supportedCurrencies"`
	CurrencyCode           string                              `json:"currencyCode"`
	AppId                  int64                               `json:"appId"`
	Enabled                bool                                `json:"enabled"`
	PaymentMethods         []*CollectionPaymentMethodOption    `json:"paymentMethods"`
	SelectedPaymentMethods []*CollectionPaymentMethodSelection `json:"selectedPaymentMethods"`
}

type CollectionPaymentMethodOption struct {
	CurrencyCode string `json:"currencyCode"`
	PayType      string `json:"payType"`
	InBankCode   string `json:"inBankCode"`
	MinAmount    string `json:"minAmount"`
	MaxAmount    string `json:"maxAmount"`
	Description  string `json:"description"`
	Available    bool   `json:"available"`
}

type CollectionPaymentMethodSelection struct {
	CurrencyCode string `json:"currencyCode"`
	PayType      string `json:"payType"`
	InBankCode   string `json:"inBankCode"`
}

type CollectionCountryCfgGroup struct {
	Continent string                      `json:"continent"`
	Countries []*CollectionCountryCfgItem `json:"countries"`
}

type GetCollectionCountryCfgRes struct {
	Continents    []*CollectionCountryCfgGroup `json:"continents"`
	PaymentTypes  []string                     `json:"paymentTypes"`
	DefaultAppIds map[string]int64             `json:"defaultAppIds"`
}

const SaveCollectionCountryCfgSectionBasic = "basic"

type SaveCollectionCountryCfgReq struct {
	g.Meta       `path:"/saveCollectionCountryCfg" method:"post" summary:"保存HaiPay国家支付配置" tags:"HaiPay支付配置"`
	SaveSection  string `json:"saveSection" v:"required|in:basic#保存类型不能为空|保存类型无效"`
	CountryCode  string `json:"countryCode" v:"required|length:2,8#国家地区不能为空|国家地区编码无效"`
	CurrencyCode string `json:"currencyCode" v:"required|length:3,8#支付币种不能为空|支付币种无效"`
	Enabled      bool   `json:"enabled"`
}

type SaveCollectionCountryCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type SaveCoinMerchantCollectionCountryCfgReq struct {
	g.Meta       `path:"/saveCollectionCountryCfg" method:"post" summary:"保存币商HaiPay国家代收配置" tags:"HaiPay支付配置"`
	CountryCode  string `json:"countryCode" v:"required|length:2,8#国家地区不能为空|国家地区编码无效"`
	CurrencyCode string `json:"currencyCode" v:"required|length:3,8#支付币种不能为空|支付币种无效"`
	AppId        int64  `json:"appId" v:"min:1#AppId无效"`
	PayType      string `json:"payType" v:"required#支付类型不能为空"`
	InBankCode   string `json:"inBankCode" v:"required#支付编码不能为空"`
	Enabled      bool   `json:"enabled"`
}

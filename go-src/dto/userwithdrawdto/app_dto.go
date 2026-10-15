package userwithdrawdto

import "github.com/gogf/gf/v2/frame/g"

type AppGetWithdrawCfgReq struct {
	g.Meta `path:"/getWithdrawCfg" method:"post" summary:"查询用户提现配置" tags:"用户提现"`
}

type AppGetWithdrawCfgRes struct {
	Enabled       bool `json:"enabled"`
	DailyMaxCount int  `json:"dailyMaxCount" dc:"每天最多提现次数(0=不限制)"`
}

type AppListWithdrawCurrenciesReq struct {
	g.Meta `path:"/listWithdrawCurrencies" method:"post" summary:"App查询支持的提现币种" tags:"用户提现"`
}

// AppWithdrawCurrencyItem 提现代付币种(对齐充值 fiatCurrencyListForApp 的区域展示字段)
type AppWithdrawCurrencyItem struct {
	CountryCode string `json:"countryCode" dc:"国家/地区简码"`
	Currency    string `json:"currency" dc:"代付 ISO 4217 币种,保存收款信息与代付使用"`
	Continent   string `json:"continent" dc:"HaiPay 代付地区分组"`
	NameEn      string `json:"nameEn"`
	NameZh      string `json:"nameZh"`
	Icon        string `json:"icon" dc:"国旗 URL"`
	Sort        int    `json:"sort"`
}

type AppListWithdrawCurrenciesRes struct {
	List []*AppWithdrawCurrencyItem `json:"list"`
}

type AppGetPayoutInfoReq struct {
	g.Meta `path:"/getPayoutInfo" method:"post" summary:"查询收款信息" tags:"用户提现"`
}

type AppPayoutInfoItem struct {
	CountryCode  string `json:"countryCode"`
	Currency     string `json:"currency"`
	AccountType  string `json:"accountType"`
	PayeeName    string `json:"payeeName"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	BankName     string `json:"bankName"`
	AccountNo    string `json:"accountNo"`
	BankCode     string `json:"bankCode"`
	IdentifyType string `json:"identifyType"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	Address3     string `json:"address3"`
	PostalCode   string `json:"postalCode"`
	Remark       string `json:"remark"`
	UpdatedAt    string `json:"updatedAt"`
}

type AppPayoutCountryOption struct {
	CountryCode  string                      `json:"countryCode"`
	NameEn       string                      `json:"nameEn"`
	NameZh       string                      `json:"nameZh"`
	Currency     string                      `json:"currency"`
	Region       string                      `json:"region"`
	AccountTypes []string                    `json:"accountTypes"`
	Wallets      []AppPayoutWalletOption     `json:"wallets"`
	Methods      []AppPayoutMethodOption     `json:"methods"`
	Icon         string                      `json:"icon"`
}

type AppPayoutWalletOption struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type AppPayoutMethodOption struct {
	AccountType          string   `json:"accountType"`
	BankCode             string   `json:"bankCode"`
	Limit                string   `json:"limit"`
	Description          string   `json:"description"`
	IdentifyTypeRequired bool     `json:"identifyTypeRequired"`
	IdentifyTypeOptions  []string `json:"identifyTypeOptions"`
	CountryRequired      bool     `json:"countryRequired"`
	AddressRequired      bool     `json:"addressRequired"`
}

type AppGetPayoutInfoRes struct {
	Info      *AppPayoutInfoItem        `json:"info"`
	Countries []*AppPayoutCountryOption `json:"countries"`
}

type AppSavePayoutInfoReq struct {
	g.Meta       `path:"/savePayoutInfo" method:"post" summary:"保存收款信息" tags:"用户提现"`
	CountryCode  string `json:"countryCode"`
	Currency     string `json:"currency"`
	AccountType  string `json:"accountType"`
	PayeeName    string `json:"payeeName"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	AccountNo    string `json:"accountNo"`
	BankCode     string `json:"bankCode"`
	IdentifyType string `json:"identifyType"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	Address3     string `json:"address3"`
	PostalCode   string `json:"postalCode"`
	Remark       string `json:"remark"`
}

type AppSavePayoutInfoRes struct {
	Success bool `json:"success"`
}

type AppApplyWithdrawReq struct {
	g.Meta     `path:"/applyWithdraw" method:"post" summary:"申请提现" tags:"用户提现"`
	TierId     uint64  `json:"tierId" v:"required|min:1#请选择提现档位"`
	Currency   string  `json:"currency" v:"required#请选择提现币种" dc:"代付 ISO 4217 币种,须与收款信息一致且在 CMS 已启用"`
	GoldAmount float64 `json:"goldAmount" dc:"已废弃,请传 tierId"`
}

type AppApplyWithdrawRes struct {
	Id            string  `json:"id"`
	DiamondAmount float64 `json:"diamondAmount"`
	WithdrawUsd   float64 `json:"withdrawUsd"`
	Status        uint8   `json:"status"`
}

type AppWithdrawLogListReq struct {
	g.Meta   `path:"/listWithdrawLogs" method:"post" summary:"我的提现记录" tags:"用户提现"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type AppWithdrawLogItem struct {
	Id                  string  `json:"id"`
	DiamondAmount       float64 `json:"diamondAmount"`
	WithdrawUsd         float64 `json:"withdrawUsd"`
	Status              uint8   `json:"status"`
	TransferCurrency    string  `json:"transferCurrency"`
	TransferLocalAmount float64 `json:"transferLocalAmount"`
	TransferFailMsg     string  `json:"transferFailMsg"`
	RejectReason        string  `json:"rejectReason"`
	CreatedAt           string  `json:"createdAt"`
	TransferAt          string  `json:"transferAt"`
}

type AppWithdrawLogListRes struct {
	Total int                   `json:"total"`
	List  []*AppWithdrawLogItem `json:"list"`
}

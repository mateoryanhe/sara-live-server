package userwithdrawdto

import "github.com/gogf/gf/v2/frame/g"

type AppGetWithdrawCfgReq struct {
	g.Meta `path:"/getWithdrawCfg" method:"post" summary:"查询用户提现配置" tags:"用户提现"`
}

type AppGetWithdrawCfgRes struct {
	Enabled  bool    `json:"enabled"`
	MinUsd   float64 `json:"minUsd"`
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
	GoldAmount float64 `json:"goldAmount" v:"required|min:0.01#金币数量无效"`
}

type AppApplyWithdrawRes struct {
	Id         string  `json:"id"`
	GoldAmount float64 `json:"goldAmount"`
	WithdrawUsd float64 `json:"withdrawUsd"`
	Status     uint8   `json:"status"`
}

type AppWithdrawLogListReq struct {
	g.Meta   `path:"/listWithdrawLogs" method:"post" summary:"我的提现记录" tags:"用户提现"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type AppWithdrawLogItem struct {
	Id                  string  `json:"id"`
	GoldAmount          float64 `json:"goldAmount"`
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

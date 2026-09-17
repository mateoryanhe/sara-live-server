package haipaydto

import "github.com/gogf/gf/v2/frame/g"

type GetHaiPayCfgReq struct {
	g.Meta `path:"/getHaiPayCfg" method:"post" summary:"查询HaiPay支付配置" tags:"HaiPay支付配置"`
}

type HaiPayCfgItem struct {
	ID                 string `json:"id"`
	ApiHost            string `json:"apiHost"`
	GlobalCashierAppId int64  `json:"globalCashierAppId"`
	TVisable           bool   `json:"tVisable"`
	MerchantSecretKey  string `json:"merchantSecretKey"`
	MerchantPrivateKey string `json:"merchantPrivateKey"`
	CallbackBaseUrl    string `json:"callbackBaseUrl"`
	ReturnUrl          string `json:"returnUrl"`
	FailReturnUrl      string `json:"failReturnUrl"`
	CancelUrl          string `json:"cancelUrl"`
	PaymentMethods     string `json:"paymentMethods"`
	Subject            string `json:"subject"`
	PayoutEnabled      bool   `json:"payoutEnabled"`
	PayoutAppIds       string `json:"payoutAppIds"`
	PayoutSubject      string `json:"payoutSubject"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type GetHaiPayCfgRes struct {
	Cfg *HaiPayCfgItem `json:"cfg"`
}

type SaveHaiPayCfgReq struct {
	g.Meta             `path:"/saveHaiPayCfg" method:"post" summary:"保存HaiPay支付配置" tags:"HaiPay支付配置"`
	ID                 uint64 `json:"id"`
	ApiHost            string `json:"apiHost" v:"required#apiHost不能为空"`
	GlobalCashierAppId int64  `json:"globalCashierAppId" v:"min:1#全球收银台AppId必须大于0"`
	TVisable           bool   `json:"tVisable"`
	MerchantSecretKey  string `json:"merchantSecretKey" v:"required#商户密钥不能为空"`
	MerchantPrivateKey string `json:"merchantPrivateKey" v:"required#商户私钥不能为空"`
	CallbackBaseUrl    string `json:"callbackBaseUrl"`
	ReturnUrl          string `json:"returnUrl"`
	FailReturnUrl      string `json:"failReturnUrl"`
	CancelUrl          string `json:"cancelUrl"`
	PaymentMethods     string `json:"paymentMethods"`
	Subject            string `json:"subject"`
	PayoutEnabled      bool   `json:"payoutEnabled"`
	PayoutAppIds       string `json:"payoutAppIds"`
	PayoutSubject      string `json:"payoutSubject"`
}

type SaveHaiPayCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

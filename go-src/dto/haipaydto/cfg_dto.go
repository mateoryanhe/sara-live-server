package haipaydto

import "github.com/gogf/gf/v2/frame/g"

type GetHaiPayCfgReq struct {
	g.Meta `path:"/getHaiPayCfg" method:"post" summary:"查询HaiPay支付配置" tags:"HaiPay支付配置"`
}

type HaiPayCfgItem struct {
	ID                 string `json:"id"`
	Enabled            bool   `json:"enabled"`
	AppId              int64  `json:"appId"`
	ApiHost            string `json:"apiHost"`
	MerchantSecretKey  string `json:"merchantSecretKey"`
	MerchantPrivateKey string `json:"merchantPrivateKey"`
	HaiPayPublicKey    string `json:"haiPayPublicKey"`
	CallbackBaseUrl    string `json:"callbackBaseUrl"`
	ReturnUrl          string `json:"returnUrl"`
	FailReturnUrl      string `json:"failReturnUrl"`
	CancelUrl          string `json:"cancelUrl"`
	PaymentMethods     string `json:"paymentMethods"`
	Subject            string `json:"subject"`
	DefaultRegion      string `json:"defaultRegion"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type GetHaiPayCfgRes struct {
	Cfg *HaiPayCfgItem `json:"cfg"`
}

type SaveHaiPayCfgReq struct {
	g.Meta             `path:"/saveHaiPayCfg" method:"post" summary:"保存HaiPay支付配置" tags:"HaiPay支付配置"`
	ID                 uint64 `json:"id"`
	Enabled            bool   `json:"enabled"`
	AppId              int64  `json:"appId" v:"required#appId不能为空"`
	ApiHost            string `json:"apiHost" v:"required#apiHost不能为空"`
	MerchantSecretKey  string `json:"merchantSecretKey" v:"required#商户密钥不能为空"`
	MerchantPrivateKey string `json:"merchantPrivateKey" v:"required#商户私钥不能为空"`
	HaiPayPublicKey    string `json:"haiPayPublicKey" v:"required#HaiPay公钥不能为空"`
	CallbackBaseUrl    string `json:"callbackBaseUrl"`
	ReturnUrl          string `json:"returnUrl"`
	FailReturnUrl      string `json:"failReturnUrl"`
	CancelUrl          string `json:"cancelUrl"`
	PaymentMethods     string `json:"paymentMethods"`
	Subject            string `json:"subject"`
	DefaultRegion      string `json:"defaultRegion"`
}

type SaveHaiPayCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

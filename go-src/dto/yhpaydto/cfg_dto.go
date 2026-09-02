package yhpaydto

import "github.com/gogf/gf/v2/frame/g"

type GetYhPayCfgReq struct {
	g.Meta `path:"/getYhPayCfg" method:"post" summary:"查询yhpay支付配置" tags:"yhpay支付配置"`
}

type YhPayCfgItem struct {
	ID              string `json:"id"`
	Enabled         bool   `json:"enabled"`
	MerchantCode    string `json:"merchantCode"`
	ApiKey          string `json:"apiKey"`
	ApiHost         string `json:"apiHost"`
	CryptoApiHost   string `json:"cryptoApiHost"`
	CallbackBaseUrl string `json:"callbackBaseUrl"`
	ReturnUrl       string `json:"returnUrl"`
	FailedReturnUrl string `json:"failedReturnUrl"`
	CryptoNetwork   string `json:"cryptoNetwork"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type GetYhPayCfgRes struct {
	Cfg *YhPayCfgItem `json:"cfg"`
}

type SaveYhPayCfgReq struct {
	g.Meta          `path:"/saveYhPayCfg" method:"post" summary:"保存yhpay支付配置" tags:"yhpay支付配置"`
	ID              uint64 `json:"id" dc:"配置ID,新建传0"`
	Enabled         bool   `json:"enabled" dc:"是否启用"`
	MerchantCode    string `json:"merchantCode" v:"required|length:1,64#Merchant ID不能为空|Merchant ID长度需在1到64之间" dc:"Merchant ID"`
	ApiKey          string `json:"apiKey" v:"required|length:1,256#API Key不能为空|API Key长度需在1到256之间" dc:"API Key"`
	ApiHost         string `json:"apiHost" v:"required|length:1,256#法币API Host不能为空|法币API Host长度需在1到256之间" dc:"法币API Host(DepositV2)"`
	CryptoApiHost   string `json:"cryptoApiHost" v:"required|length:1,256#USDT API Host不能为空|USDT API Host长度需在1到256之间" dc:"USDT加密入款API Host"`
	CallbackBaseUrl string `json:"callbackBaseUrl" v:"max-length:256#回调根地址最长256" dc:"本服务对外回调根地址"`
	ReturnUrl       string `json:"returnUrl" v:"max-length:512#成功跳转最长512" dc:"支付成功浏览器跳转"`
	FailedReturnUrl string `json:"failedReturnUrl" v:"max-length:512#失败跳转最长512" dc:"支付失败浏览器跳转"`
	CryptoNetwork   string `json:"cryptoNetwork" v:"max-length:32#网络最长32" dc:"USDT网络(固定TRC20)"`
}

type SaveYhPayCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

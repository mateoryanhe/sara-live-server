package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbYhPayCfg db.TbName = "yhpay_cfgs"
)

// YhPayCfg yhpay 第三方支付配置(CMS 管理,通常仅一条)
type YhPayCfg struct {
	migrate.OneModel
	Enabled         bool   `gorm:"default:0;comment:是否启用" json:"enabled"`
	MerchantCode    string `gorm:"size:64;default:'';comment:Merchant ID/商户号" json:"merchantCode"`
	ApiKey          string `gorm:"size:256;default:'';comment:API Key/密钥" json:"apiKey"`
	ApiHost         string `gorm:"size:256;default:'';comment:API Host(IDR手动入款)" json:"apiHost"`
	CryptoApiHost   string `gorm:"size:256;default:'';comment:USDT加密入款API Host(预留)" json:"cryptoApiHost"`
	CallbackBaseUrl string `gorm:"size:256;default:'';comment:本服务对外回调根地址(如https://www.example.com)" json:"callbackBaseUrl"`
	ReturnUrl       string `gorm:"size:512;default:'';comment:支付成功浏览器跳转" json:"returnUrl"`
	FailedReturnUrl string `gorm:"size:512;default:'';comment:支付失败浏览器跳转" json:"failedReturnUrl"`
	CryptoNetwork   string `gorm:"size:32;default:'TRC20';comment:USDT网络(固定波场TRC20)" json:"cryptoNetwork"`
}

func (YhPayCfg) TableName() string {
	return string(TbYhPayCfg)
}

func initYhPayCfg() {
	migrate.AutoMigrate(&YhPayCfg{})
}

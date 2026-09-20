package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const TbHaiPayCfg db.TbName = "haipay_cfgs"

// HaiPayCfg HaiPay 代收 + 代付配置，CMS 管理，通常一条
type HaiPayCfg struct {
	migrate.OneModel
	ApiHost            string `gorm:"size:256;default:'';comment:API Host如https://interface.haipay.asia" json:"apiHost"`
	TVisable           bool   `gorm:"column:t_visable;default:0;comment:App是否显示第三方支付" json:"tVisable"`
	MerchantSecretKey  string `gorm:"size:256;default:'';comment:签名串末尾key=商户密钥" json:"merchantSecretKey"`
	MerchantPrivateKey string `gorm:"type:text;comment:商户RSA私钥(PKCS8,可无PEM头)" json:"merchantPrivateKey"`
	CallbackBaseUrl    string `gorm:"size:512;default:'';comment:回调根地址" json:"callbackBaseUrl"`
	ReturnUrl          string `gorm:"size:512;default:'';comment:支付成功跳转" json:"returnUrl"`
	FailReturnUrl      string `gorm:"size:512;default:'';comment:支付失败跳转" json:"failReturnUrl"`
	CancelUrl          string `gorm:"size:512;default:'';comment:取消支付跳转(可选)" json:"cancelUrl"`
	PaymentMethods     string `gorm:"size:256;default:'';comment:历史全局收银台支付方式,本地代收不使用" json:"paymentMethods"`
	Subject            string `gorm:"size:128;default:'Recharge';comment:支付标题" json:"subject"`
}

func (HaiPayCfg) TableName() string {
	return string(TbHaiPayCfg)
}

func initHaiPayCfg() {
	migrate.AutoMigrate(&HaiPayCfg{})
}

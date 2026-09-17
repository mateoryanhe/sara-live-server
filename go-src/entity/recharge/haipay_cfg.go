package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbHaiPayCfg                     db.TbName = "haipay_cfgs"
	DefaultHaiPayGlobalCashierAppID int64     = 25272
)

// HaiPayCfg HaiPay 代收 + 代付配置，CMS 管理，通常一条
type HaiPayCfg struct {
	migrate.OneModel
	ApiHost            string `gorm:"size:256;default:'';comment:API Host如https://interface.haipay.asia" json:"apiHost"`
	GlobalCashierAppId int64  `gorm:"not null;default:25272;comment:普通用户全球收银台AppId" json:"globalCashierAppId"`
	TVisable           bool   `gorm:"column:t_visable;default:0;comment:App是否显示第三方支付" json:"tVisable"`
	MerchantSecretKey  string `gorm:"size:256;default:'';comment:签名串末尾key=商户密钥" json:"merchantSecretKey"`
	MerchantPrivateKey string `gorm:"type:text;comment:商户RSA私钥(PKCS8,可无PEM头)" json:"merchantPrivateKey"`
	CallbackBaseUrl    string `gorm:"size:512;default:'';comment:回调根地址" json:"callbackBaseUrl"`
	ReturnUrl          string `gorm:"size:512;default:'';comment:支付成功跳转" json:"returnUrl"`
	FailReturnUrl      string `gorm:"size:512;default:'';comment:支付失败跳转" json:"failReturnUrl"`
	CancelUrl          string `gorm:"size:512;default:'';comment:取消支付跳转(可选)" json:"cancelUrl"`
	PaymentMethods     string `gorm:"size:256;default:'';comment:历史全局收银台支付方式,本地代收不使用" json:"paymentMethods"`
	Subject            string `gorm:"size:128;default:'Recharge';comment:支付标题" json:"subject"`
	PayoutEnabled      bool   `gorm:"default:0;comment:是否启用工会代付" json:"payoutEnabled"`
	// PayoutAppIds 代付 appId 映射,如 IDR:25280,PHP:25281 (按币种,与收银台 appId 通常不同)
	PayoutAppIds  string `gorm:"size:512;default:'';comment:代付appId映射 CUR:appId 逗号分隔" json:"payoutAppIds"`
	PayoutSubject string `gorm:"size:128;default:'GuildSettlement';comment:代付标题" json:"payoutSubject"`
}

func (HaiPayCfg) TableName() string {
	return string(TbHaiPayCfg)
}

func initHaiPayCfg() {
	migrate.AutoMigrate(&HaiPayCfg{})
}

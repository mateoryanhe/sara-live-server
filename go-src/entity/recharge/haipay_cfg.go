package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbHaiPayCfg db.TbName = "haipay_cfgs"
)

// HaiPayCfg HaiPay 全球收银台 + 代付配置，CMS 管理，通常一条
type HaiPayCfg struct {
	migrate.OneModel
	Enabled            bool   `gorm:"default:0;comment:是否启用收银台代收" json:"enabled"`
	AppId              int64  `gorm:"default:0;comment:收银台专属appId" json:"appId"`
	ApiHost            string `gorm:"size:256;default:'';comment:API Host如https://interface.haipay.asia" json:"apiHost"`
	MerchantSecretKey  string `gorm:"size:256;default:'';comment:签名串末尾key=商户密钥" json:"merchantSecretKey"`
	MerchantPrivateKey string `gorm:"type:text;comment:商户RSA私钥(PKCS8,可无PEM头)" json:"merchantPrivateKey"`
	HaiPayPublicKey    string `gorm:"type:text;comment:HaiPay平台公钥(可无PEM头)" json:"haiPayPublicKey"`
	CallbackBaseUrl    string `gorm:"size:512;default:'';comment:回调根地址" json:"callbackBaseUrl"`
	ReturnUrl          string `gorm:"size:512;default:'';comment:支付成功跳转" json:"returnUrl"`
	FailReturnUrl      string `gorm:"size:512;default:'';comment:支付失败跳转" json:"failReturnUrl"`
	CancelUrl          string `gorm:"size:512;default:'';comment:取消支付跳转(可选)" json:"cancelUrl"`
	PaymentMethods     string `gorm:"size:256;default:'';comment:可选支付方式逗号分隔" json:"paymentMethods"`
	Subject            string `gorm:"size:128;default:'Recharge';comment:支付标题" json:"subject"`
	DefaultRegion      string `gorm:"size:8;default:'ID';comment:全球收银台默认region(如ID/PH/US)" json:"defaultRegion"`
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

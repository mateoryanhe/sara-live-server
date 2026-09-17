package entity

import (
	"strings"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const TbHaiPayCoinMerchantCollectionCfg db.TbName = "haipay_coin_merchant_collection_cfgs"

// HaiPayCoinMerchantCollectionCfg 币商本地代收配置。
// 每个国家/地区只保存一行，币种、appId、支付方式和 App 可见性一起生效。
type HaiPayCoinMerchantCollectionCfg struct {
	CountryCode  string    `gorm:"primaryKey;size:8;comment:HaiPay地区编码" json:"countryCode"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	CurrencyCode string    `gorm:"size:8;not null;comment:本地代收接口下单币种" json:"currencyCode"`
	AppId        int64     `gorm:"not null;comment:所选币种对应的HaiPay appId" json:"appId"`
	PayType      string    `gorm:"size:64;not null;comment:支付类型" json:"payType"`
	InBankCode   string    `gorm:"size:256;not null;comment:支付编码" json:"inBankCode"`
	Enabled      bool      `gorm:"not null;comment:币商App是否展示该国家地区" json:"enabled"`
}

func (HaiPayCoinMerchantCollectionCfg) TableName() string {
	return string(TbHaiPayCoinMerchantCollectionCfg)
}

func (r *HaiPayCoinMerchantCollectionCfg) Normalize() {
	if r == nil {
		return
	}
	r.CountryCode = strings.ToUpper(strings.TrimSpace(r.CountryCode))
	r.CurrencyCode = strings.ToUpper(strings.TrimSpace(r.CurrencyCode))
	r.PayType = strings.ToUpper(strings.TrimSpace(r.PayType))
	r.InBankCode = strings.TrimSpace(r.InBankCode)
}

func initHaiPayCoinMerchantCollectionCfg() {
	migrate.AutoMigrate(&HaiPayCoinMerchantCollectionCfg{})
}

package entity

import (
	"strings"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const TbHaiPayUserWithdrawPayoutCountryCfg db.TbName = "haipay_user_withdraw_payout_country_cfgs"

// HaiPayUserWithdrawPayoutCountryCfg 用户提现代付国家/地区配置(一国一币种,含默认代付方式).
type HaiPayUserWithdrawPayoutCountryCfg struct {
	CountryCode  string    `gorm:"primaryKey;size:8;comment:HaiPay地区编码" json:"countryCode"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	CurrencyCode string    `gorm:"size:8;not null;comment:代付 ISO 4217 币种" json:"currencyCode"`
	AppId        int64     `gorm:"not null;comment:服务端按币种自动写入的HaiPay appId" json:"-"`
	AccountType  string    `gorm:"size:64;not null;comment:账户类型(BANK_ACCOUNT/EWALLET)" json:"accountType"`
	BankCode     string    `gorm:"size:256;not null;comment:代付支付编码" json:"bankCode"`
	Enabled      bool      `gorm:"not null;comment:App是否展示该国家/地区" json:"enabled"`
}

func (HaiPayUserWithdrawPayoutCountryCfg) TableName() string {
	return string(TbHaiPayUserWithdrawPayoutCountryCfg)
}

func (r *HaiPayUserWithdrawPayoutCountryCfg) Normalize() {
	if r == nil {
		return
	}
	r.CountryCode = strings.ToUpper(strings.TrimSpace(r.CountryCode))
	r.CurrencyCode = strings.ToUpper(strings.TrimSpace(r.CurrencyCode))
	r.AccountType = strings.ToUpper(strings.TrimSpace(r.AccountType))
	r.BankCode = strings.TrimSpace(r.BankCode)
}

func initHaiPayUserWithdrawPayoutCountryCfg() {
	migrate.AutoMigrate(&HaiPayUserWithdrawPayoutCountryCfg{})
}

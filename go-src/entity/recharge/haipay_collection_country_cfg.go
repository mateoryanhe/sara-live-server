package entity

import (
	"strconv"
	"strings"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const TbHaiPayCollectionCountryCfg db.TbName = "haipay_collection_country_cfgs"

// HaiPayCollectionCountryCfg 普通用户全球收银台国家/地区配置。
// 币商本地代收使用独立的 HaiPayCoinMerchantCollectionCfg 表，不写入本表。
type HaiPayCollectionCountryCfg struct {
	ID           string        `gorm:"primaryKey;size:32;comment:业务类型_国家编码" json:"id"`
	CreatedAt    time.Time     `json:"-"`
	UpdatedAt    time.Time     `json:"-"`
	BizType      HaiPayBizType `gorm:"uniqueIndex:uniq_haipay_collection_country,priority:1;not null;comment:业务类型,1普通用户代收,2币商代收,3普通用户代付,4工会代付" json:"bizType"`
	CountryCode  string        `gorm:"size:8;index;uniqueIndex:uniq_haipay_collection_country,priority:2;not null;comment:HaiPay地区编码" json:"countryCode"`
	CurrencyCode string        `gorm:"size:8;not null;comment:本地代收接口下单币种" json:"currencyCode"`
	Enabled      bool          `gorm:"not null;comment:App是否展示该国家地区" json:"enabled"`
}

func (HaiPayCollectionCountryCfg) TableName() string {
	return string(TbHaiPayCollectionCountryCfg)
}

// HaiPayCollectionCountryCfgID returns the stable primary key shared by the
// country master row and its payment-method rows, for example 1_ID or 2_HK.
func HaiPayCollectionCountryCfgID(bizType HaiPayBizType, countryCode string) string {
	return strconv.FormatUint(uint64(bizType), 10) + "_" + strings.ToUpper(strings.TrimSpace(countryCode))
}

func initHaiPayCollectionCountryCfg() {
	migrate.AutoMigrate(&HaiPayCollectionCountryCfg{})
}

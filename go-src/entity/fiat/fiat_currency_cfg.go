package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbFiatCurrencyCfg db.TbName = "fiat_currency_cfgs"
)

const (
	FiatCurrencyCfgCurrencyCode   db.TbCol = "currency_code"
	FiatCurrencyCfgName           db.TbCol = "name"
	FiatCurrencyCfgSymbol         db.TbCol = "symbol"
	FiatCurrencyCfgIcon           db.TbCol = "icon"
	FiatCurrencyCfgAdjustPercent  db.TbCol = "adjust_percent"
	FiatCurrencyCfgCurrencyType   db.TbCol = "currency_type"
	FiatCurrencyCfgSort           db.TbCol = "sort"
	FiatCurrencyCfgStatus         db.TbCol = "status"
)

const (
	FiatCurrencyStatusDisabled uint8 = 0
	FiatCurrencyStatusEnabled  uint8 = 1
)

const (
	FiatCurrencyTypeFiat   uint8 = 1 // 法币
	FiatCurrencyTypeCrypto uint8 = 2 // 加密币
)

// FiatCurrencyCfg 法币/加密币配置(CMS 管理,供服务端汇率换算与白名单)
type FiatCurrencyCfg struct {
	migrate.OneModel
	CurrencyCode  string  `gorm:"size:8;uniqueIndex;default:'';comment:币种代码(ISO4217,如IDR)" json:"currencyCode"`
	Name          string  `gorm:"size:64;default:'';comment:币种名称" json:"name"`
	Symbol        string  `gorm:"size:16;default:'';comment:币种符号" json:"symbol"`
	Icon          string  `gorm:"size:255;default:'';comment:图标资源文件名" json:"icon"`
	AdjustPercent float64 `gorm:"type:decimal(8,4);default:0;comment:汇率加点比例(%)" json:"adjustPercent"`
	CurrencyType  uint8   `gorm:"default:1;comment:币种类型(1法币,2加密币)" json:"currencyType"`
	Sort          int     `gorm:"default:0;comment:排序(越大越靠前)" json:"sort"`
	Status        uint8   `gorm:"default:0;comment:状态(0禁用,1启用)" json:"status"`
}

func initFiatCurrencyCfg() {
	migrate.AutoMigrate(&FiatCurrencyCfg{})
}

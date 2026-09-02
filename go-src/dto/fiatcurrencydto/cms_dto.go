package fiatcurrencydto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// FiatCurrencyListReq CMS 分页查询法币配置
type FiatCurrencyListReq struct {
	g.Meta `path:"/fiatCurrencyList" method:"post" summary:"查询法币配置列表" tags:"法币配置"`
	httpserver.CMSQueryReq
	CurrencyCode string `json:"currencyCode" dc:"币种代码(模糊)"`
	Name         string `json:"name" dc:"名称(模糊)"`
	TypeFilter   int    `json:"typeFilter" dc:"币种类型(0全部,1法币,2加密币)"`
	StatusFilter int    `json:"statusFilter" dc:"状态(0全部,1禁用,2启用)"`
}

// FiatCurrencyItem CMS 法币配置项
type FiatCurrencyItem struct {
	ID            string  `json:"id"`
	CurrencyCode  string  `json:"currencyCode"`
	Name          string  `json:"name"`
	Symbol        string  `json:"symbol"`
	Icon          string  `json:"icon" dc:"图标完整URL"`
	IconName      string  `json:"iconName" dc:"图标资源文件名"`
	AdjustPercent float64 `json:"adjustPercent"`
	CurrencyType  uint8   `json:"currencyType"`
	Sort          int     `json:"sort"`
	Status        uint8   `json:"status"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

type CreateFiatCurrencyReq struct {
	g.Meta        `path:"/createFiatCurrency" method:"post" summary:"创建法币配置" tags:"法币配置"`
	CurrencyCode  string  `json:"currencyCode" v:"required|length:3,8#币种代码不能为空|币种代码长度需在3到8之间" dc:"币种代码(如IDR)"`
	Name          string  `json:"name" v:"required|length:1,64#名称不能为空|名称长度需在1到64之间" dc:"币种名称"`
	Symbol        string  `json:"symbol" v:"required|length:1,16#符号不能为空|符号长度需在1到16之间" dc:"币种符号"`
	Icon          string  `json:"icon" v:"max-length:255#图标最长255字符" dc:"图标资源文件名"`
	AdjustPercent float64 `json:"adjustPercent" dc:"汇率加点比例(%)"`
	CurrencyType  uint8   `json:"currencyType" v:"required|in:1,2#币种类型不能为空|币种类型无效" dc:"币种类型(1法币,2加密币)"`
	Sort          int     `json:"sort" dc:"排序"`
	Status        uint8   `json:"status" v:"in:0,1#状态无效" dc:"状态(0禁用,1启用)"`
}

type CreateFiatCurrencyRes struct {
	ID string `json:"id"`
}

type UpdateFiatCurrencyReq struct {
	g.Meta        `path:"/updateFiatCurrency" method:"post" summary:"修改法币配置" tags:"法币配置"`
	ID            uint64  `json:"id" v:"required#ID不能为空" dc:"配置ID"`
	CurrencyCode  string  `json:"currencyCode" v:"required|length:3,8#币种代码不能为空|币种代码长度需在3到8之间" dc:"币种代码(如IDR)"`
	Name          string  `json:"name" v:"required|length:1,64#名称不能为空|名称长度需在1到64之间" dc:"币种名称"`
	Symbol        string  `json:"symbol" v:"required|length:1,16#符号不能为空|符号长度需在1到16之间" dc:"币种符号"`
	Icon          string  `json:"icon" v:"max-length:255#图标最长255字符" dc:"图标资源文件名"`
	AdjustPercent float64 `json:"adjustPercent" dc:"汇率加点比例(%)"`
	CurrencyType  uint8   `json:"currencyType" v:"required|in:1,2#币种类型不能为空|币种类型无效" dc:"币种类型(1法币,2加密币)"`
	Sort          int     `json:"sort" dc:"排序"`
	Status        uint8   `json:"status" v:"in:0,1#状态无效" dc:"状态(0禁用,1启用)"`
}

type UpdateFiatCurrencyRes struct {
	Success bool `json:"success"`
}

type DeleteFiatCurrencyReq struct {
	g.Meta `path:"/deleteFiatCurrency" method:"post" summary:"删除法币配置" tags:"法币配置"`
	ID     uint64 `json:"id" v:"required#ID不能为空" dc:"配置ID"`
}

type DeleteFiatCurrencyRes struct {
	Success bool `json:"success"`
}

type ReloadFiatCurrencyCacheReq struct {
	g.Meta `path:"/reloadFiatCurrencyCache" method:"post" summary:"刷新法币配置缓存" tags:"法币配置"`
}

type ReloadFiatCurrencyCacheRes struct {
	Success bool `json:"success"`
}

type ReloadFiatExchangeRateCacheReq struct {
	g.Meta       `path:"/reloadFiatExchangeRateCache" method:"post" summary:"刷新法币汇率缓存" tags:"法币配置"`
	CurrencyCode string `json:"currencyCode" dc:"币种代码(空=全部)"`
}

type ReloadFiatExchangeRateCacheRes struct {
	Success bool `json:"success"`
}

type GetFiatExchangeRateReq struct {
	g.Meta       `path:"/getFiatExchangeRate" method:"post" summary:"查询USD兑法币汇率" tags:"法币配置"`
	CurrencyCode string `json:"currencyCode" v:"required|length:3,8#币种代码不能为空|币种代码长度需在3到8之间" dc:"币种代码(如IDR)"`
}

type GetFiatExchangeRateRes struct {
	Base           string  `json:"base"`
	Quote          string  `json:"quote"`
	MarketRate     float64 `json:"marketRate" dc:"市场参考汇率(1 USD = X Quote)"`
	AdjustPercent  float64 `json:"adjustPercent" dc:"加点比例(%)"`
	Rate           float64 `json:"rate" dc:"最终汇率(1 USD = X Quote)"`
	InverseRate    float64 `json:"inverseRate" dc:"1 Quote = X USD"`
	Source         string  `json:"source"`
	RateDate       string  `json:"rateDate"`
	Cached         bool    `json:"cached"`
	CacheExpiresAt int64   `json:"cacheExpiresAt" dc:"缓存过期Unix秒,0=未知"`
}

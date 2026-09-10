package fiatcurrency

import (
	"context"

	"xr-game-server/dto/fiatcurrencydto"
	fiatentity "xr-game-server/entity/fiat"
)

// 硬编码 HaiPay 全球收银台 region 列表(App 选用后原样上报 currencyCode)。
// 与 module/recharge 内 haiPayValidRegions 保持一致即可。
var appHaiPayRegions = []struct {
	Code string
	Name string
}{
	{Code: "ID", Name: "Indonesia"},
	{Code: "PH", Name: "Philippines"},
	{Code: "MY", Name: "Malaysia"},
	{Code: "IN", Name: "India"},
	{Code: "TH", Name: "Thailand"},
	{Code: "VN", Name: "Vietnam"},
	{Code: "SG", Name: "Singapore"},
	{Code: "HK", Name: "Hong Kong"},
	{Code: "TW", Name: "Taiwan"},
	{Code: "JP", Name: "Japan"},
	{Code: "KR", Name: "South Korea"},
	{Code: "PK", Name: "Pakistan"},
	{Code: "BR", Name: "Brazil"},
	{Code: "US", Name: "United States"},
	{Code: "GB", Name: "United Kingdom"},
	{Code: "EU", Name: "European Union"},
	{Code: "IT", Name: "Italy"},
	{Code: "AT", Name: "Austria"},
	{Code: "BE", Name: "Belgium"},
	{Code: "NL", Name: "Netherlands"},
	{Code: "PL", Name: "Poland"},
	{Code: "TR", Name: "Türkiye"},
	{Code: "AE", Name: "United Arab Emirates"},
	{Code: "SA", Name: "Saudi Arabia"},
	{Code: "QA", Name: "Qatar"},
	{Code: "KW", Name: "Kuwait"},
	{Code: "BH", Name: "Bahrain"},
	{Code: "OM", Name: "Oman"},
	{Code: "EG", Name: "Egypt"},
}

// GetAppList App 查询区域列表(硬编码 HaiPay region; currencyCode 即 region)
func GetAppList(_ context.Context, req *fiatcurrencydto.AppFiatCurrencyListReq) (*fiatcurrencydto.AppFiatCurrencyListRes, error) {
	if req != nil && req.TypeFilter == int(fiatentity.FiatCurrencyTypeCrypto) {
		return &fiatcurrencydto.AppFiatCurrencyListRes{List: []*fiatcurrencydto.AppFiatCurrencyItem{}}, nil
	}
	list := make([]*fiatcurrencydto.AppFiatCurrencyItem, 0, len(appHaiPayRegions))
	for i, row := range appHaiPayRegions {
		list = append(list, &fiatcurrencydto.AppFiatCurrencyItem{
			CurrencyCode: row.Code,
			Name:         row.Name,
			Symbol:       row.Code,
			CurrencyType: fiatentity.FiatCurrencyTypeFiat,
			Sort:         len(appHaiPayRegions) - i,
		})
	}
	return &fiatcurrencydto.AppFiatCurrencyListRes{List: list}, nil
}

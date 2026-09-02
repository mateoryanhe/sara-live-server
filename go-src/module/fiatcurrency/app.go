package fiatcurrency

import (
	"context"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/fiatcurrencydto"
	fiatentity "xr-game-server/entity/fiat"
	"xr-game-server/module/upload"
)

// GetAppList App 查询已启用币种列表(走内存缓存)
func GetAppList(_ context.Context, req *fiatcurrencydto.AppFiatCurrencyListReq) (*fiatcurrencydto.AppFiatCurrencyListRes, error) {
	if req == nil {
		req = &fiatcurrencydto.AppFiatCurrencyListReq{}
	}
	rows := cfgdao.GetFiatCurrencyCfgListCached()
	list := make([]*fiatcurrencydto.AppFiatCurrencyItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 || row.Status != fiatentity.FiatCurrencyStatusEnabled {
			continue
		}
		switch req.TypeFilter {
		case int(fiatentity.FiatCurrencyTypeFiat):
			if row.CurrencyType != fiatentity.FiatCurrencyTypeFiat {
				continue
			}
		case int(fiatentity.FiatCurrencyTypeCrypto):
			if row.CurrencyType != fiatentity.FiatCurrencyTypeCrypto {
				continue
			}
		}
		list = append(list, &fiatcurrencydto.AppFiatCurrencyItem{
			CurrencyCode: row.CurrencyCode,
			Name:         row.Name,
			Symbol:       row.Symbol,
			Icon:         upload.GetUrlByName(row.Icon),
			CurrencyType: row.CurrencyType,
			Sort:         row.Sort,
		})
	}
	return &fiatcurrencydto.AppFiatCurrencyListRes{List: list}, nil
}

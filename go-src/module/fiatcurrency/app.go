package fiatcurrency

import (
	"context"

	"xr-game-server/constants/country"
	"xr-game-server/dto/fiatcurrencydto"
	fiatentity "xr-game-server/entity/fiat"
	"xr-game-server/module/countryflagdeploy"
	"xr-game-server/module/upload"
)

// GetAppList App 查询 HaiPay 支付区域(统一 constants/country; currencyCode=region 简码)
func GetAppList(_ context.Context, req *fiatcurrencydto.AppFiatCurrencyListReq) (*fiatcurrencydto.AppFiatCurrencyListRes, error) {
	if req != nil && req.TypeFilter == int(fiatentity.FiatCurrencyTypeCrypto) {
		return &fiatcurrencydto.AppFiatCurrencyListRes{List: []*fiatcurrencydto.AppFiatCurrencyItem{}}, nil
	}
	return &fiatcurrencydto.AppFiatCurrencyListRes{List: buildHaiPayRegionItems()}, nil
}

// GetHaiPayRegionList CMS/App 共用的 HaiPay region 选型列表
func GetHaiPayRegionList(_ context.Context, _ *fiatcurrencydto.HaiPayRegionListReq) (*fiatcurrencydto.AppFiatCurrencyListRes, error) {
	return &fiatcurrencydto.AppFiatCurrencyListRes{List: buildHaiPayRegionItems()}, nil
}

func buildHaiPayRegionItems() []*fiatcurrencydto.AppFiatCurrencyItem {
	regions := country.ListHaiPayRegions()
	version := countryflagdeploy.CurrentVersion()
	list := make([]*fiatcurrencydto.AppFiatCurrencyItem, 0, len(regions))
	for i, c := range regions {
		icon := ""
		if rel := country.RelPath(c.Code, version); rel != "" {
			icon = upload.GetUrlByName(rel)
		}
		list = append(list, &fiatcurrencydto.AppFiatCurrencyItem{
			CurrencyCode: c.Code,
			Name:         c.NameEn,
			NameEn:       c.NameEn,
			NameZh:       c.NameZh,
			Symbol:       c.Code,
			Icon:         icon,
			CurrencyType: fiatentity.FiatCurrencyTypeFiat,
			Sort:         len(regions) - i,
		})
	}
	return list
}

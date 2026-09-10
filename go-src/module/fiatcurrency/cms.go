package fiatcurrency

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/fiatcurrencydto"
	fiatentity "xr-game-server/entity/fiat"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func normalizeCurrencyCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func normalizeCurrencyType(currencyType uint8) uint8 {
	if currencyType == fiatentity.FiatCurrencyTypeCrypto {
		return fiatentity.FiatCurrencyTypeCrypto
	}
	return fiatentity.FiatCurrencyTypeFiat
}

func validateCurrencyCode(code string) error {
	code = normalizeCurrencyCode(code)
	if len(code) < 3 || len(code) > 8 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if code == "USD" {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}

func reloadCaches() {
	cfgdao.ReloadFiatCurrencyCfgCache()
}

// GetList CMS 分页查询法币配置
func GetList(_ context.Context, req *fiatcurrencydto.FiatCurrencyListReq) (*httpserver.CMSQueryResp, error) {
	total, rows := cfgdao.GetFiatCurrencyCfgList(req)
	list := make([]*fiatcurrencydto.FiatCurrencyItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		item := &fiatcurrencydto.FiatCurrencyItem{
			ID:           row.ID,
			CurrencyCode: row.CurrencyCode,
			Name:         row.Name,
			Symbol:       row.Symbol,
			IconName:     row.Icon,
			CurrencyType: row.CurrencyType,
			Sort:         row.Sort,
			Status:       row.Status,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
		}
		item.Icon = upload.GetUrlByName(item.IconName)
		list = append(list, item)
	}
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func Create(_ context.Context, req *fiatcurrencydto.CreateFiatCurrencyReq) (*fiatcurrencydto.CreateFiatCurrencyRes, error) {
	code := normalizeCurrencyCode(req.CurrencyCode)
	if err := validateCurrencyCode(code); err != nil {
		return nil, err
	}
	if existing := cfgdao.GetFiatCurrencyCfgByCode(code); existing != nil {
		return nil, errercode.CreateCode(errercode.FiatCurrencyExist)
	}
	row := &fiatentity.FiatCurrencyCfg{
		CurrencyCode: code,
		Name:         strings.TrimSpace(req.Name),
		Symbol:       strings.TrimSpace(req.Symbol),
		Icon:         strings.TrimSpace(req.Icon),
		CurrencyType: normalizeCurrencyType(req.CurrencyType),
		Sort:         req.Sort,
		Status:       req.Status,
	}
	if err := cfgdao.CreateFiatCurrencyCfg(row); err != nil {
		return nil, err
	}
	reloadCaches()
	return &fiatcurrencydto.CreateFiatCurrencyRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *fiatcurrencydto.UpdateFiatCurrencyReq) (*fiatcurrencydto.UpdateFiatCurrencyRes, error) {
	code := normalizeCurrencyCode(req.CurrencyCode)
	if err := validateCurrencyCode(code); err != nil {
		return nil, err
	}
	row := cfgdao.GetFiatCurrencyCfgById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.FiatCurrencyNonExist)
	}
	if existing := cfgdao.GetFiatCurrencyCfgByCode(code); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.FiatCurrencyExist)
	}
	row.CurrencyCode = code
	row.Name = strings.TrimSpace(req.Name)
	row.Symbol = strings.TrimSpace(req.Symbol)
	row.Icon = strings.TrimSpace(req.Icon)
	row.CurrencyType = normalizeCurrencyType(req.CurrencyType)
	row.Sort = req.Sort
	row.Status = req.Status
	if err := cfgdao.UpdateFiatCurrencyCfg(row); err != nil {
		return nil, err
	}
	reloadCaches()
	return &fiatcurrencydto.UpdateFiatCurrencyRes{Success: true}, nil
}

func Delete(_ context.Context, req *fiatcurrencydto.DeleteFiatCurrencyReq) (*fiatcurrencydto.DeleteFiatCurrencyRes, error) {
	row := cfgdao.GetFiatCurrencyCfgById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.FiatCurrencyNonExist)
	}
	if err := cfgdao.DeleteFiatCurrencyCfg(req.ID); err != nil {
		return nil, err
	}
	reloadCaches()
	return &fiatcurrencydto.DeleteFiatCurrencyRes{Success: true}, nil
}

func ReloadCfgCache(_ context.Context, _ *fiatcurrencydto.ReloadFiatCurrencyCacheReq) (*fiatcurrencydto.ReloadFiatCurrencyCacheRes, error) {
	cfgdao.ReloadFiatCurrencyCfgCache()
	return &fiatcurrencydto.ReloadFiatCurrencyCacheRes{Success: true}, nil
}

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
	"xr-game-server/module/fxrate"
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

func reloadCaches(currencyCode string) {
	cfgdao.ReloadFiatCurrencyCfgCache()
	fxrate.ReloadExchangeRateCache(currencyCode)
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
			ID:            row.ID,
			CurrencyCode:  row.CurrencyCode,
			Name:          row.Name,
			Symbol:        row.Symbol,
			IconName:      row.Icon,
			AdjustPercent: row.AdjustPercent,
			CurrencyType:  row.CurrencyType,
			Sort:          row.Sort,
			Status:        row.Status,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
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
		CurrencyCode:  code,
		Name:          strings.TrimSpace(req.Name),
		Symbol:        strings.TrimSpace(req.Symbol),
		Icon:          strings.TrimSpace(req.Icon),
		AdjustPercent: req.AdjustPercent,
		CurrencyType:  normalizeCurrencyType(req.CurrencyType),
		Sort:          req.Sort,
		Status:        req.Status,
	}
	if err := cfgdao.CreateFiatCurrencyCfg(row); err != nil {
		return nil, err
	}
	reloadCaches(code)
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
	oldCode := row.CurrencyCode
	row.CurrencyCode = code
	row.Name = strings.TrimSpace(req.Name)
	row.Symbol = strings.TrimSpace(req.Symbol)
	row.Icon = strings.TrimSpace(req.Icon)
	row.AdjustPercent = req.AdjustPercent
	row.CurrencyType = normalizeCurrencyType(req.CurrencyType)
	row.Sort = req.Sort
	row.Status = req.Status
	if err := cfgdao.UpdateFiatCurrencyCfg(row); err != nil {
		return nil, err
	}
	reloadCaches(oldCode)
	if oldCode != code {
		reloadCaches(code)
	}
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
	reloadCaches(row.CurrencyCode)
	return &fiatcurrencydto.DeleteFiatCurrencyRes{Success: true}, nil
}

func ReloadCfgCache(_ context.Context, _ *fiatcurrencydto.ReloadFiatCurrencyCacheReq) (*fiatcurrencydto.ReloadFiatCurrencyCacheRes, error) {
	cfgdao.ReloadFiatCurrencyCfgCache()
	return &fiatcurrencydto.ReloadFiatCurrencyCacheRes{Success: true}, nil
}

func ReloadExchangeRateCache(_ context.Context, req *fiatcurrencydto.ReloadFiatExchangeRateCacheReq) (*fiatcurrencydto.ReloadFiatExchangeRateCacheRes, error) {
	code := normalizeCurrencyCode(req.CurrencyCode)
	if code != "" {
		if err := validateCurrencyCode(code); err != nil {
			return nil, err
		}
	}
	fxrate.ReloadExchangeRateCache(code)
	return &fiatcurrencydto.ReloadFiatExchangeRateCacheRes{Success: true}, nil
}

func GetExchangeRate(ctx context.Context, req *fiatcurrencydto.GetFiatExchangeRateReq) (*fiatcurrencydto.GetFiatExchangeRateRes, error) {
	code := normalizeCurrencyCode(req.CurrencyCode)
	if err := validateCurrencyCode(code); err != nil {
		return nil, err
	}
	cfg := cfgdao.GetFiatCurrencyCfgByCode(code)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.FiatCurrencyNonExist)
	}
	if cfg.Status != fiatentity.FiatCurrencyStatusEnabled {
		return nil, errercode.CreateCode(errercode.FiatCurrencyDisabled)
	}
	rate, err := fxrate.GetUsdToQuoteRate(ctx, code, cfg.AdjustPercent)
	if err != nil {
		return nil, errercode.CreateCode(errercode.FiatExchangeRateUnavailable)
	}
	return toExchangeRateRes(rate), nil
}

// ResolveEnabledUsdToQuoteRate 服务端内部:按 CMS 启用的法币配置查询最终汇率
func ResolveEnabledUsdToQuoteRate(ctx context.Context, currencyCode string) (*fxrate.QuoteRate, error) {
	code := normalizeCurrencyCode(currencyCode)
	if err := validateCurrencyCode(code); err != nil {
		return nil, err
	}
	cfg := cfgdao.GetFiatCurrencyCfgByCode(code)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.FiatCurrencyNonExist)
	}
	if cfg.Status != fiatentity.FiatCurrencyStatusEnabled {
		return nil, errercode.CreateCode(errercode.FiatCurrencyDisabled)
	}
	rate, err := fxrate.GetUsdToQuoteRate(ctx, code, cfg.AdjustPercent)
	if err != nil {
		return nil, errercode.CreateCode(errercode.FiatExchangeRateUnavailable)
	}
	return rate, nil
}

func toExchangeRateRes(rate *fxrate.QuoteRate) *fiatcurrencydto.GetFiatExchangeRateRes {
	if rate == nil {
		return &fiatcurrencydto.GetFiatExchangeRateRes{}
	}
	return &fiatcurrencydto.GetFiatExchangeRateRes{
		Base:           rate.Base,
		Quote:          rate.Quote,
		MarketRate:     rate.MarketRate,
		AdjustPercent:  rate.AdjustPercent,
		Rate:           rate.Rate,
		InverseRate:    rate.InverseRate,
		Source:         rate.Source,
		RateDate:       rate.RateDate,
		Cached:         rate.Cached,
		CacheExpiresAt: rate.CacheExpiresAt,
	}
}

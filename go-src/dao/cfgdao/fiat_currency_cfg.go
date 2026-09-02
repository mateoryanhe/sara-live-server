package cfgdao

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/dto/fiatcurrencydto"
	fiatentity "xr-game-server/entity/fiat"
)

const fiatCurrencyCfgListCacheKey = "fiat_currency_cfg_list"

var fiatCurrencyCfgListCache *cache.ListCache[*fiatentity.FiatCurrencyCfg]

func InitFiatCurrencyCfgDao() {
	fiatCurrencyCfgListCache = cache.NewPermanentListCache[*fiatentity.FiatCurrencyCfg]()
}

func loadFiatCurrencyCfgListFromDB() []*fiatentity.FiatCurrencyCfg {
	rows := make([]*fiatentity.FiatCurrencyCfg, 0)
	_ = g.DB().Model(string(fiatentity.TbFiatCurrencyCfg)).
		Order("sort desc, currency_code asc").
		Scan(&rows)
	return rows
}

func ReloadFiatCurrencyCfgCache() {
	if fiatCurrencyCfgListCache == nil {
		return
	}
	fiatCurrencyCfgListCache.PublishList(gctx.New(), fiatCurrencyCfgListCacheKey, loadFiatCurrencyCfgListFromDB())
}

func GetFiatCurrencyCfgListCached() []*fiatentity.FiatCurrencyCfg {
	if fiatCurrencyCfgListCache == nil {
		return loadFiatCurrencyCfgListFromDB()
	}
	return fiatCurrencyCfgListCache.MustGetList(gctx.New(), fiatCurrencyCfgListCacheKey, func(ctx context.Context) ([]*fiatentity.FiatCurrencyCfg, error) {
		return loadFiatCurrencyCfgListFromDB(), nil
	})
}

func GetEnabledFiatCurrencyCfgMapCached() map[string]*fiatentity.FiatCurrencyCfg {
	rows := GetFiatCurrencyCfgListCached()
	ret := make(map[string]*fiatentity.FiatCurrencyCfg, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 || row.Status != fiatentity.FiatCurrencyStatusEnabled {
			continue
		}
		code := normalizeFiatCurrencyCode(row.CurrencyCode)
		if code == "" || code == "USD" {
			continue
		}
		ret[code] = row
	}
	return ret
}

func GetFiatCurrencyCfgById(id uint64) *fiatentity.FiatCurrencyCfg {
	if id == 0 {
		return nil
	}
	for _, row := range GetFiatCurrencyCfgListCached() {
		if row != nil && row.ID == id {
			return row
		}
	}
	return nil
}

func GetFiatCurrencyCfgByCode(currencyCode string) *fiatentity.FiatCurrencyCfg {
	code := normalizeFiatCurrencyCode(currencyCode)
	if code == "" {
		return nil
	}
	for _, row := range GetFiatCurrencyCfgListCached() {
		if row != nil && normalizeFiatCurrencyCode(row.CurrencyCode) == code {
			return row
		}
	}
	return nil
}

func CreateFiatCurrencyCfg(row *fiatentity.FiatCurrencyCfg) error {
	_, err := g.DB().Model(string(fiatentity.TbFiatCurrencyCfg)).Save(row)
	return err
}

func UpdateFiatCurrencyCfg(row *fiatentity.FiatCurrencyCfg) error {
	return CreateFiatCurrencyCfg(row)
}

func DeleteFiatCurrencyCfg(id uint64) error {
	_, err := g.DB().Model(string(fiatentity.TbFiatCurrencyCfg)).WherePri(id).Delete()
	return err
}

func GetFiatCurrencyCfgList(req *fiatcurrencydto.FiatCurrencyListReq) (int, []*fiatcurrencydto.FiatCurrencyItem) {
	if req == nil {
		return 0, nil
	}
	filtered := make([]*fiatentity.FiatCurrencyCfg, 0)
	for _, row := range GetFiatCurrencyCfgListCached() {
		if matchFiatCurrencyListFilter(row, req) {
			filtered = append(filtered, row)
		}
	}
	total := len(filtered)
	offset := req.PageOffset()
	if offset >= total {
		return total, []*fiatcurrencydto.FiatCurrencyItem{}
	}
	end := offset + req.PageSize
	if end > total {
		end = total
	}
	pageRows := filtered[offset:end]
	ret := make([]*fiatcurrencydto.FiatCurrencyItem, 0, len(pageRows))
	for _, row := range pageRows {
		if item := toFiatCurrencyListItem(row); item != nil {
			ret = append(ret, item)
		}
	}
	return total, ret
}

func matchFiatCurrencyListFilter(row *fiatentity.FiatCurrencyCfg, req *fiatcurrencydto.FiatCurrencyListReq) bool {
	if row == nil || row.ID == 0 {
		return false
	}
	if code := strings.TrimSpace(req.CurrencyCode); code != "" {
		if !strings.Contains(strings.ToUpper(row.CurrencyCode), strings.ToUpper(code)) {
			return false
		}
	}
	if name := strings.TrimSpace(req.Name); name != "" {
		if !strings.Contains(row.Name, name) {
			return false
		}
	}
	switch req.TypeFilter {
	case int(fiatentity.FiatCurrencyTypeFiat):
		if row.CurrencyType != fiatentity.FiatCurrencyTypeFiat {
			return false
		}
	case int(fiatentity.FiatCurrencyTypeCrypto):
		if row.CurrencyType != fiatentity.FiatCurrencyTypeCrypto {
			return false
		}
	}
	switch req.StatusFilter {
	case 1:
		if row.Status != fiatentity.FiatCurrencyStatusDisabled {
			return false
		}
	case 2:
		if row.Status != fiatentity.FiatCurrencyStatusEnabled {
			return false
		}
	}
	return true
}

func toFiatCurrencyListItem(row *fiatentity.FiatCurrencyCfg) *fiatcurrencydto.FiatCurrencyItem {
	if row == nil || row.ID == 0 {
		return nil
	}
	return &fiatcurrencydto.FiatCurrencyItem{
		ID:            strconv.FormatUint(row.ID, 10),
		CurrencyCode:  row.CurrencyCode,
		Name:          row.Name,
		Symbol:        row.Symbol,
		Icon:          row.Icon,
		AdjustPercent: row.AdjustPercent,
		CurrencyType:  row.CurrencyType,
		Sort:          row.Sort,
		Status:        row.Status,
		CreatedAt:     FormatFiatCurrencyTime(row.CreatedAt),
		UpdatedAt:     FormatFiatCurrencyTime(row.UpdatedAt),
	}
}

func normalizeFiatCurrencyCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func FormatFiatCurrencyTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

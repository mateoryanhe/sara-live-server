package cfgdao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/country"
	"xr-game-server/core/cache"
	rechargeentity "xr-game-server/entity/recharge"
)

const haiPayUserWithdrawPayoutCountryCfgCacheKey = "haipay_user_withdraw_payout_country_cfg_list"

var haiPayUserWithdrawPayoutCountryCfgCache *cache.ListCache[*rechargeentity.HaiPayUserWithdrawPayoutCountryCfg]

func InitHaiPayUserWithdrawPayoutCountryCfgDao() {
	haiPayUserWithdrawPayoutCountryCfgCache = cache.NewPermanentListCache[*rechargeentity.HaiPayUserWithdrawPayoutCountryCfg]()
}

func loadHaiPayUserWithdrawPayoutCountryCfgList() []*rechargeentity.HaiPayUserWithdrawPayoutCountryCfg {
	rows := make([]*rechargeentity.HaiPayUserWithdrawPayoutCountryCfg, 0)
	_ = g.DB().Model(string(rechargeentity.TbHaiPayUserWithdrawPayoutCountryCfg)).
		Order("country_code asc").Scan(&rows)
	return rows
}

func ReloadHaiPayUserWithdrawPayoutCountryCfgCache() {
	if haiPayUserWithdrawPayoutCountryCfgCache == nil {
		return
	}
	haiPayUserWithdrawPayoutCountryCfgCache.PublishList(
		gctx.New(), haiPayUserWithdrawPayoutCountryCfgCacheKey, loadHaiPayUserWithdrawPayoutCountryCfgList(),
	)
}

func getHaiPayUserWithdrawPayoutCountryCfgListCached() []*rechargeentity.HaiPayUserWithdrawPayoutCountryCfg {
	if haiPayUserWithdrawPayoutCountryCfgCache == nil {
		return loadHaiPayUserWithdrawPayoutCountryCfgList()
	}
	return haiPayUserWithdrawPayoutCountryCfgCache.MustGetList(
		gctx.New(), haiPayUserWithdrawPayoutCountryCfgCacheKey,
		func(context.Context) ([]*rechargeentity.HaiPayUserWithdrawPayoutCountryCfg, error) {
			return loadHaiPayUserWithdrawPayoutCountryCfgList(), nil
		},
	)
}

func GetHaiPayUserWithdrawPayoutCountryCfgMapCached() map[string]*rechargeentity.HaiPayUserWithdrawPayoutCountryCfg {
	result := make(map[string]*rechargeentity.HaiPayUserWithdrawPayoutCountryCfg)
	for _, row := range getHaiPayUserWithdrawPayoutCountryCfgListCached() {
		if row == nil {
			continue
		}
		countryCode := strings.ToUpper(strings.TrimSpace(row.CountryCode))
		if countryCode != "" {
			result[countryCode] = row
		}
	}
	return result
}

func GetHaiPayUserWithdrawPayoutCountryCfgCached(countryCode string) *rechargeentity.HaiPayUserWithdrawPayoutCountryCfg {
	return GetHaiPayUserWithdrawPayoutCountryCfgMapCached()[strings.ToUpper(strings.TrimSpace(countryCode))]
}

func HaiPayUserWithdrawPayoutCountryCfgComplete(row *rechargeentity.HaiPayUserWithdrawPayoutCountryCfg) bool {
	return row != nil && strings.TrimSpace(row.CountryCode) != "" &&
		strings.TrimSpace(row.CurrencyCode) != "" &&
		strings.TrimSpace(row.AccountType) != "" && strings.TrimSpace(row.BankCode) != ""
}

func fillHaiPayUserWithdrawPayoutCountryCfgAppID(row *rechargeentity.HaiPayUserWithdrawPayoutCountryCfg) error {
	if row == nil {
		return nil
	}
	if row.AppId > 0 {
		return nil
	}
	appID, ok := country.LookupHaiPayAppID(row.CurrencyCode, HaiPayUseProdAppID())
	if !ok {
		return fmt.Errorf("haipay app_id not found for currency %s", strings.ToUpper(strings.TrimSpace(row.CurrencyCode)))
	}
	row.AppId = appID
	return nil
}

func SaveHaiPayUserWithdrawPayoutCountryCfg(row *rechargeentity.HaiPayUserWithdrawPayoutCountryCfg) (string, error) {
	if row == nil {
		return "", nil
	}
	row.Normalize()
	if err := fillHaiPayUserWithdrawPayoutCountryCfgAppID(row); err != nil {
		return "", err
	}
	now := time.Now()
	if existing := GetHaiPayUserWithdrawPayoutCountryCfgCached(row.CountryCode); existing != nil && !existing.CreatedAt.IsZero() {
		row.CreatedAt = existing.CreatedAt
	} else if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	row.UpdatedAt = now
	// 显式列写入，避免 Save(struct) 省略 app_id 导致 NOT NULL 报错
	_, err := g.DB().Model(string(rechargeentity.TbHaiPayUserWithdrawPayoutCountryCfg)).Data(g.Map{
		"country_code":  row.CountryCode,
		"created_at":    row.CreatedAt,
		"updated_at":    row.UpdatedAt,
		"currency_code": row.CurrencyCode,
		"app_id":        row.AppId,
		"account_type":  row.AccountType,
		"bank_code":     row.BankCode,
		"enabled":       row.Enabled,
	}).Save()
	if err != nil {
		return "", err
	}
	return row.CountryCode, nil
}

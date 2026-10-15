package cfgdao

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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

func SaveHaiPayUserWithdrawPayoutCountryCfg(row *rechargeentity.HaiPayUserWithdrawPayoutCountryCfg) (string, error) {
	if row == nil {
		return "", nil
	}
	row.Normalize()
	if _, err := g.DB().Model(string(rechargeentity.TbHaiPayUserWithdrawPayoutCountryCfg)).Save(row); err != nil {
		return "", err
	}
	return row.CountryCode, nil
}

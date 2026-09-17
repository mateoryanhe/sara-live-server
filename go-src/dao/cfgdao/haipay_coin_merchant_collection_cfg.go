package cfgdao

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	rechargeentity "xr-game-server/entity/recharge"
)

const haiPayCoinMerchantCollectionCfgCacheKey = "haipay_coin_merchant_collection_cfg_list"

var haiPayCoinMerchantCollectionCfgCache *cache.ListCache[*rechargeentity.HaiPayCoinMerchantCollectionCfg]

func InitHaiPayCoinMerchantCollectionCfgDao() {
	haiPayCoinMerchantCollectionCfgCache = cache.NewPermanentListCache[*rechargeentity.HaiPayCoinMerchantCollectionCfg]()
}

func loadHaiPayCoinMerchantCollectionCfgList() []*rechargeentity.HaiPayCoinMerchantCollectionCfg {
	rows := make([]*rechargeentity.HaiPayCoinMerchantCollectionCfg, 0)
	_ = g.DB().Model(string(rechargeentity.TbHaiPayCoinMerchantCollectionCfg)).
		Order("country_code asc").Scan(&rows)
	return rows
}

func ReloadHaiPayCoinMerchantCollectionCfgCache() {
	if haiPayCoinMerchantCollectionCfgCache == nil {
		return
	}
	haiPayCoinMerchantCollectionCfgCache.PublishList(
		gctx.New(), haiPayCoinMerchantCollectionCfgCacheKey, loadHaiPayCoinMerchantCollectionCfgList(),
	)
}

func getHaiPayCoinMerchantCollectionCfgListCached() []*rechargeentity.HaiPayCoinMerchantCollectionCfg {
	if haiPayCoinMerchantCollectionCfgCache == nil {
		return loadHaiPayCoinMerchantCollectionCfgList()
	}
	return haiPayCoinMerchantCollectionCfgCache.MustGetList(
		gctx.New(), haiPayCoinMerchantCollectionCfgCacheKey,
		func(context.Context) ([]*rechargeentity.HaiPayCoinMerchantCollectionCfg, error) {
			return loadHaiPayCoinMerchantCollectionCfgList(), nil
		},
	)
}

func GetHaiPayCoinMerchantCollectionCfgMapCached() map[string]*rechargeentity.HaiPayCoinMerchantCollectionCfg {
	result := make(map[string]*rechargeentity.HaiPayCoinMerchantCollectionCfg)
	for _, row := range getHaiPayCoinMerchantCollectionCfgListCached() {
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

func GetHaiPayCoinMerchantCollectionCfgCached(countryCode string) *rechargeentity.HaiPayCoinMerchantCollectionCfg {
	return GetHaiPayCoinMerchantCollectionCfgMapCached()[strings.ToUpper(strings.TrimSpace(countryCode))]
}

func HaiPayCoinMerchantCollectionCfgComplete(row *rechargeentity.HaiPayCoinMerchantCollectionCfg) bool {
	return row != nil && strings.TrimSpace(row.CountryCode) != "" &&
		strings.TrimSpace(row.CurrencyCode) != "" && row.AppId > 0 &&
		strings.TrimSpace(row.PayType) != "" && strings.TrimSpace(row.InBankCode) != ""
}

func SaveHaiPayCoinMerchantCollectionCfg(row *rechargeentity.HaiPayCoinMerchantCollectionCfg) (string, error) {
	if row == nil {
		return "", nil
	}
	row.Normalize()
	if _, err := g.DB().Model(string(rechargeentity.TbHaiPayCoinMerchantCollectionCfg)).Save(row); err != nil {
		return "", err
	}
	return row.CountryCode, nil
}

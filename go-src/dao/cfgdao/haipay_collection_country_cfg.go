package cfgdao

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	rechargeentity "xr-game-server/entity/recharge"
)

const haiPayCollectionCountryCfgCacheKey = "haipay_collection_country_cfg_list"

type HaiPayCollectionCountryCfgAggregate struct {
	Config *rechargeentity.HaiPayCollectionCountryCfg
}

var haiPayCollectionCountryCfgCache *cache.ListCache[*rechargeentity.HaiPayCollectionCountryCfg]

func InitHaiPayCollectionCountryCfgDao() {
	haiPayCollectionCountryCfgCache = cache.NewPermanentListCache[*rechargeentity.HaiPayCollectionCountryCfg]()
}

func loadHaiPayCollectionCountryCfgList() []*rechargeentity.HaiPayCollectionCountryCfg {
	rows := make([]*rechargeentity.HaiPayCollectionCountryCfg, 0)
	_ = g.DB().Model(string(rechargeentity.TbHaiPayCollectionCountryCfg)).
		Order("biz_type asc, country_code asc").Scan(&rows)
	return rows
}

func ReloadHaiPayCollectionCountryCfgCache() {
	if haiPayCollectionCountryCfgCache == nil {
		return
	}
	ctx := gctx.New()
	haiPayCollectionCountryCfgCache.PublishList(ctx, haiPayCollectionCountryCfgCacheKey, loadHaiPayCollectionCountryCfgList())
}

func getHaiPayCollectionCountryCfgListCached() []*rechargeentity.HaiPayCollectionCountryCfg {
	if haiPayCollectionCountryCfgCache == nil {
		return loadHaiPayCollectionCountryCfgList()
	}
	return haiPayCollectionCountryCfgCache.MustGetList(gctx.New(), haiPayCollectionCountryCfgCacheKey, func(context.Context) ([]*rechargeentity.HaiPayCollectionCountryCfg, error) {
		return loadHaiPayCollectionCountryCfgList(), nil
	})
}

func GetHaiPayCollectionCountryCfgMapCached(bizType rechargeentity.HaiPayBizType) map[string]*HaiPayCollectionCountryCfgAggregate {
	result := make(map[string]*HaiPayCollectionCountryCfgAggregate)
	for _, row := range getHaiPayCollectionCountryCfgListCached() {
		if row == nil || strings.TrimSpace(row.ID) == "" || row.BizType != bizType {
			continue
		}
		aggregate := &HaiPayCollectionCountryCfgAggregate{Config: row}
		result[strings.ToUpper(strings.TrimSpace(row.CountryCode))] = aggregate
	}
	return result
}

func GetHaiPayCollectionCountryCfgCached(bizType rechargeentity.HaiPayBizType, countryCode string) *HaiPayCollectionCountryCfgAggregate {
	return GetHaiPayCollectionCountryCfgMapCached(bizType)[strings.ToUpper(strings.TrimSpace(countryCode))]
}

// SaveHaiPayCollectionCountryBasicCfg 只保存普通用户全球收银台国家配置。
func SaveHaiPayCollectionCountryBasicCfg(config *rechargeentity.HaiPayCollectionCountryCfg) (string, error) {
	if config == nil {
		return "", nil
	}
	config.CountryCode = strings.ToUpper(strings.TrimSpace(config.CountryCode))
	config.CurrencyCode = strings.ToUpper(strings.TrimSpace(config.CurrencyCode))
	config.ID = rechargeentity.HaiPayCollectionCountryCfgID(config.BizType, config.CountryCode)
	_, err := g.DB().Model(string(rechargeentity.TbHaiPayCollectionCountryCfg)).Save(config)
	if err != nil {
		return "", err
	}
	return config.ID, nil
}

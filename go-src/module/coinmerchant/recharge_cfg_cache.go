package coinmerchant

import (
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/coinmerchantrechargecfgdto"
	"xr-game-server/entity/recharge"
)

type coinMerchantRechargeCfgSnapshot struct {
	byId map[uint64]*coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem
	list []*coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem
}

var (
	coinMerchantRechargeCfgCache     atomic.Value // *coinMerchantRechargeCfgSnapshot
	emptyCoinMerchantRechargeCfgList = make([]*coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem, 0)
)

func toAppCoinMerchantRechargeCfgItem(r *entity.CoinMerchantRechargeCfg) *coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem {
	return &coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem{
		ID:    r.ID,
		Name:  r.Name,
		Price: r.Price,
		Gold:  r.Gold,
	}
}

func loadCoinMerchantRechargeCfgCache() []*coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem {
	rows := cfgdao.GetOnShelfCoinMerchantRechargeCfg()
	byId := make(map[uint64]*coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem, len(rows))
	list := make([]*coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem, 0, len(rows))
	for _, r := range rows {
		item := toAppCoinMerchantRechargeCfgItem(r)
		byId[r.ID] = item
		list = append(list, item)
	}
	coinMerchantRechargeCfgCache.Store(&coinMerchantRechargeCfgSnapshot{byId: byId, list: list})
	return list
}

func reloadCoinMerchantRechargeCfgCache() {
	loadCoinMerchantRechargeCfgCache()
}

// ReloadCoinMerchantRechargeCfgCache 供同步等模块调用
func ReloadCoinMerchantRechargeCfgCache() {
	reloadCoinMerchantRechargeCfgCache()
}

func getCoinMerchantRechargeCfgSnapshot() *coinMerchantRechargeCfgSnapshot {
	v := coinMerchantRechargeCfgCache.Load()
	if v == nil {
		return &coinMerchantRechargeCfgSnapshot{
			byId: make(map[uint64]*coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem),
			list: emptyCoinMerchantRechargeCfgList,
		}
	}
	return v.(*coinMerchantRechargeCfgSnapshot)
}

// GetOnShelfCoinMerchantRechargeCfgs 获取已上架币商充值档位(缓存)
func GetOnShelfCoinMerchantRechargeCfgs() []*coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem {
	if coinMerchantRechargeCfgCache.Load() == nil {
		return loadCoinMerchantRechargeCfgCache()
	}
	return getCoinMerchantRechargeCfgSnapshot().list
}

// GetCoinMerchantRechargeCfgFromCacheById 按 ID 取已上架档位
func GetCoinMerchantRechargeCfgFromCacheById(id uint64) *coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgItem {
	if coinMerchantRechargeCfgCache.Load() == nil {
		loadCoinMerchantRechargeCfgCache()
	}
	return getCoinMerchantRechargeCfgSnapshot().byId[id]
}

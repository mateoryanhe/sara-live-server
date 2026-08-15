package cfgdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/user"
)

const walletExchangeCfgCacheKey = "wallet_exchange_cfg"

var walletExchangeCfgCacheMgr *cache.CacheMgr

func InitWalletExchangeCfgDao() {
	walletExchangeCfgCacheMgr = cache.NewCacheMgr()
}

func loadWalletExchangeCfgFromDB() *entity.WalletExchangeCfg {
	var row entity.WalletExchangeCfg
	if err := g.DB().Model(string(entity.TbWalletExchangeCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveWalletExchangeCfg(row *entity.WalletExchangeCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbWalletExchangeCfg)).Save(row)
	return err
}

func ReloadWalletExchangeCfgCache() {
	if walletExchangeCfgCacheMgr == nil {
		return
	}
	walletExchangeCfgCacheMgr.FlushCache(walletExchangeCfgCacheKey, loadWalletExchangeCfgFromDB())
	walletExchangeCfgCacheMgr.Cache.UpdateExpire(gctx.New(), walletExchangeCfgCacheKey, time.Hour*24*365*100)
}

func GetWalletExchangeCfgCached() *entity.WalletExchangeCfg {
	if walletExchangeCfgCacheMgr == nil {
		return loadWalletExchangeCfgFromDB()
	}
	v := walletExchangeCfgCacheMgr.GetData(walletExchangeCfgCacheKey, func(ctx context.Context) (value interface{}, err error) {
		return loadWalletExchangeCfgFromDB(), nil
	})
	walletExchangeCfgCacheMgr.Cache.UpdateExpire(gctx.New(), walletExchangeCfgCacheKey, time.Hour*24*365*100)
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.WalletExchangeCfg)
	return row
}

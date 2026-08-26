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

var walletExchangeCfgCacheMgr *cache.RowCache[*entity.WalletExchangeCfg]

func InitWalletExchangeCfgDao() {
	walletExchangeCfgCacheMgr = cache.NewRowCache[*entity.WalletExchangeCfg]()
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
	walletExchangeCfgCacheMgr.PublishRow(gctx.New(), walletExchangeCfgCacheKey, loadWalletExchangeCfgFromDB())
	_ = walletExchangeCfgCacheMgr.SetRow(gctx.New(), walletExchangeCfgCacheKey, loadWalletExchangeCfgFromDB(), time.Hour*24*365*100)
}

func GetWalletExchangeCfgCached() *entity.WalletExchangeCfg {
	if walletExchangeCfgCacheMgr == nil {
		return loadWalletExchangeCfgFromDB()
	}
	v := walletExchangeCfgCacheMgr.MustGetRow(gctx.New(), walletExchangeCfgCacheKey, func(ctx context.Context) (*entity.WalletExchangeCfg, error) {
		return loadWalletExchangeCfgFromDB(), nil
	})
	_ = walletExchangeCfgCacheMgr.SetRow(gctx.New(), walletExchangeCfgCacheKey, v, time.Hour*24*365*100)
	return v
}

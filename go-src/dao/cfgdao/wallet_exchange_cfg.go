package cfgdao

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/user"
)

const (
	walletExchangeCfgCacheKey = "wallet_exchange_cfg"

	DefaultWalletGoldToDiamondRate            = 100
	DefaultWalletExchangeFeePercent           = 3.0
	DefaultWalletUsdToGoldRate                = 100
	DefaultEffectiveLiveMinSessionMinutes     = 30
	MaxEffectiveLiveMinSessionMinutes         = 24 * 60
	DefaultPlatformAnchorMinimumSettlementUsd = 5.0
	MaxPlatformAnchorMinimumSettlementUsd     = 1000000.0
)

var walletExchangeCfgCacheMgr *cache.RowCache[*entity.WalletExchangeCfg]
var walletExchangeCfgWriteMu sync.Mutex

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

// UpdateWalletExchangeCfg 串行更新统一配置行，避免三个独立保存入口互相覆盖字段。
// expectedID 大于0时必须与数据库现有配置行一致；matched=false 表示前端提交了失效ID。
func UpdateWalletExchangeCfg(expectedID uint64, update func(*entity.WalletExchangeCfg)) (row *entity.WalletExchangeCfg, matched bool, err error) {
	walletExchangeCfgWriteMu.Lock()
	defer walletExchangeCfgWriteMu.Unlock()

	existing := loadWalletExchangeCfgFromDB()
	if expectedID > 0 && (existing == nil || existing.ID != expectedID) {
		return nil, false, nil
	}
	if existing == nil {
		row = &entity.WalletExchangeCfg{
			GoldToDiamondRate:                  DefaultWalletGoldToDiamondRate,
			ExchangeFeePercent:                 DefaultWalletExchangeFeePercent,
			UsdToGoldRate:                      DefaultWalletUsdToGoldRate,
			EffectiveLiveMinSessionMinutes:     DefaultEffectiveLiveMinSessionMinutes,
			PlatformAnchorMinimumSettlementUsd: DefaultPlatformAnchorMinimumSettlementUsd,
		}
	} else {
		copyRow := *existing
		row = &copyRow
	}
	if update != nil {
		update(row)
	}
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err = SaveWalletExchangeCfg(row); err != nil {
		return nil, true, err
	}
	ReloadWalletExchangeCfgCache()
	return row, true, nil
}

func ReloadWalletExchangeCfgCache() {
	if walletExchangeCfgCacheMgr == nil {
		return
	}
	cfg := loadWalletExchangeCfgFromDB()
	walletExchangeCfgCacheMgr.PublishRow(gctx.New(), walletExchangeCfgCacheKey, cfg)
	_ = walletExchangeCfgCacheMgr.SetRow(gctx.New(), walletExchangeCfgCacheKey, cfg, time.Hour*24*365*100)
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

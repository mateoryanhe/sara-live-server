package cfgdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

const (
	effectiveLiveCfgCacheKey              = "effective_live_cfg"
	DefaultEffectiveLiveMinSessionMinutes = 30
	MaxEffectiveLiveMinSessionMinutes     = 24 * 60
)

var effectiveLiveCfgCacheMgr *cache.RowCache[*entity.EffectiveLiveCfg]

func InitEffectiveLiveCfgDao() {
	effectiveLiveCfgCacheMgr = cache.NewRowCache[*entity.EffectiveLiveCfg]()
}

func loadEffectiveLiveCfgFromDB() *entity.EffectiveLiveCfg {
	var row entity.EffectiveLiveCfg
	if err := g.DB().Model(string(entity.TbEffectiveLiveCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveEffectiveLiveCfg(row *entity.EffectiveLiveCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbEffectiveLiveCfg)).Save(row)
	return err
}

func ReloadEffectiveLiveCfgCache() {
	if effectiveLiveCfgCacheMgr == nil {
		return
	}
	cfg := loadEffectiveLiveCfgFromDB()
	effectiveLiveCfgCacheMgr.PublishRow(gctx.New(), effectiveLiveCfgCacheKey, cfg)
	_ = effectiveLiveCfgCacheMgr.SetRow(gctx.New(), effectiveLiveCfgCacheKey, cfg, time.Hour*24*365*100)
}

func GetEffectiveLiveCfgCached() *entity.EffectiveLiveCfg {
	if effectiveLiveCfgCacheMgr == nil {
		return loadEffectiveLiveCfgFromDB()
	}
	cfg := effectiveLiveCfgCacheMgr.MustGetRow(gctx.New(), effectiveLiveCfgCacheKey, func(ctx context.Context) (*entity.EffectiveLiveCfg, error) {
		return loadEffectiveLiveCfgFromDB(), nil
	})
	_ = effectiveLiveCfgCacheMgr.SetRow(gctx.New(), effectiveLiveCfgCacheKey, cfg, time.Hour*24*365*100)
	return cfg
}

func resolveEffectiveLiveMinSessionMinutes(cfg *entity.EffectiveLiveCfg) int {
	if cfg == nil || cfg.MinSessionMinutes <= 0 || cfg.MinSessionMinutes > MaxEffectiveLiveMinSessionMinutes {
		return DefaultEffectiveLiveMinSessionMinutes
	}
	return cfg.MinSessionMinutes
}

// EffectiveLiveMinSessionMinutes 返回进程缓存中的有效直播门槛;数据库未配置或值无效时返回30分钟.
func EffectiveLiveMinSessionMinutes() int {
	return resolveEffectiveLiveMinSessionMinutes(GetEffectiveLiveCfgCached())
}

package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	sysentity "xr-game-server/entity/sys"
)

const metaPixelCfgCacheKey = "meta_pixel_cfg"

var metaPixelCfgCacheMgr *cache.RowCache[*sysentity.MetaPixelCfg]

func InitMetaPixelCfgDao() {
	if metaPixelCfgCacheMgr != nil {
		return
	}
	metaPixelCfgCacheMgr = cache.NewPermanentRowCache[*sysentity.MetaPixelCfg]()
}

// LoadMetaPixelCfgRow 读库首条配置,不依赖缓存是否已初始化.
func LoadMetaPixelCfgRow() *sysentity.MetaPixelCfg {
	return loadMetaPixelCfgFromDB()
}

// ResolveMetaPixelCfg 优先读缓存;未初始化缓存时直读库.
func ResolveMetaPixelCfg() *sysentity.MetaPixelCfg {
	if metaPixelCfgCacheMgr != nil {
		return GetMetaPixelCfgCached()
	}
	return loadMetaPixelCfgFromDB()
}

func EnsureMetaPixelCfgDao() {
	InitMetaPixelCfgDao()
}

func loadMetaPixelCfgFromDB() *sysentity.MetaPixelCfg {
	var row sysentity.MetaPixelCfg
	if err := g.DB().Model(string(sysentity.TbMetaPixelCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveMetaPixelCfg(row *sysentity.MetaPixelCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(sysentity.TbMetaPixelCfg)).Save(row)
	return err
}

func ReloadMetaPixelCfgCache() *sysentity.MetaPixelCfg {
	if metaPixelCfgCacheMgr == nil {
		return loadMetaPixelCfgFromDB()
	}
	row := loadMetaPixelCfgFromDB()
	metaPixelCfgCacheMgr.PublishRow(gctx.New(), metaPixelCfgCacheKey, row)
	return row
}

func GetMetaPixelCfgCached() *sysentity.MetaPixelCfg {
	if metaPixelCfgCacheMgr == nil {
		return nil
	}
	row, _ := metaPixelCfgCacheMgr.GetRowCached(gctx.New(), metaPixelCfgCacheKey)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}

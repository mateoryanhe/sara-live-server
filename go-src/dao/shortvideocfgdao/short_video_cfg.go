package shortvideocfgdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const shortVideoCfgCacheKey uint64 = 1

var shortVideoCfgCacheMgr *cache.CacheMgr

func InitShortVideoCfgDao() {
	shortVideoCfgCacheMgr = cache.NewCacheMgr()
}

func Get() *entity.ShortVideoCfg {
	if shortVideoCfgCacheMgr == nil {
		return loadFromDB()
	}
	v := shortVideoCfgCacheMgr.GetData(shortVideoCfgCacheKey, func(ctx context.Context) (value interface{}, err error) {
		return loadFromDB(), nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.ShortVideoCfg)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}

func Save(row *entity.ShortVideoCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbShortVideoCfg)).Save(row)
	if err != nil {
		return err
	}
	if shortVideoCfgCacheMgr != nil {
		shortVideoCfgCacheMgr.FlushCache(shortVideoCfgCacheKey, row)
	}
	return nil
}

func loadFromDB() *entity.ShortVideoCfg {
	var row entity.ShortVideoCfg
	if err := g.DB().Model(string(entity.TbShortVideoCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

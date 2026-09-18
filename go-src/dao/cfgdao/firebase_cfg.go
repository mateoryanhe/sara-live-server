package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	sysentity "xr-game-server/entity/sys"
)

const firebaseCfgCacheKey = "firebase_cfg"

var firebaseCfgCacheMgr *cache.RowCache[*sysentity.FirebaseCfg]

func InitFirebaseCfgDao() {
	firebaseCfgCacheMgr = cache.NewPermanentRowCache[*sysentity.FirebaseCfg]()
}

func loadFirebaseCfgFromDB() *sysentity.FirebaseCfg {
	var row sysentity.FirebaseCfg
	if err := g.DB().Model(string(sysentity.TbFirebaseCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveFirebaseCfg(row *sysentity.FirebaseCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(sysentity.TbFirebaseCfg)).Save(row)
	return err
}

func ReloadFirebaseCfgCache() *sysentity.FirebaseCfg {
	if firebaseCfgCacheMgr == nil {
		return loadFirebaseCfgFromDB()
	}
	row := loadFirebaseCfgFromDB()
	firebaseCfgCacheMgr.PublishRow(gctx.New(), firebaseCfgCacheKey, row)
	return row
}

// GetFirebaseCfgCached 仅从内存读取配置.
func GetFirebaseCfgCached() *sysentity.FirebaseCfg {
	if firebaseCfgCacheMgr == nil {
		return nil
	}
	row, _ := firebaseCfgCacheMgr.GetRowCached(gctx.New(), firebaseCfgCacheKey)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}

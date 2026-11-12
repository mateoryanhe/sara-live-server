package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	sysentity "xr-game-server/entity/sys"
)

const firebasePushCfgCacheKey = "firebase_push_cfg"

var firebasePushCfgCacheMgr *cache.RowCache[*sysentity.FirebasePushCfg]

func InitFirebasePushCfgDao() {
	if firebasePushCfgCacheMgr != nil {
		return
	}
	firebasePushCfgCacheMgr = cache.NewPermanentRowCache[*sysentity.FirebasePushCfg]()
}

func EnsureFirebasePushCfgDao() {
	InitFirebasePushCfgDao()
}

func loadFirebasePushCfgFromDB() *sysentity.FirebasePushCfg {
	var row sysentity.FirebasePushCfg
	if err := g.DB().Model(string(sysentity.TbFirebasePushCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveFirebasePushCfg(row *sysentity.FirebasePushCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(sysentity.TbFirebasePushCfg)).Save(row)
	return err
}

func ReloadFirebasePushCfgCache() *sysentity.FirebasePushCfg {
	if firebasePushCfgCacheMgr == nil {
		return loadFirebasePushCfgFromDB()
	}
	row := loadFirebasePushCfgFromDB()
	firebasePushCfgCacheMgr.PublishRow(gctx.New(), firebasePushCfgCacheKey, row)
	return row
}

func GetFirebasePushCfgFromMemory() *sysentity.FirebasePushCfg {
	if firebasePushCfgCacheMgr == nil {
		return nil
	}
	row, _ := firebasePushCfgCacheMgr.GetRowCached(gctx.New(), firebasePushCfgCacheKey)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}

// ResolveFirebasePushCfg 优先读缓存;无记录返回带默认文案的空壳(ID=0,不入库).
func ResolveFirebasePushCfg() *sysentity.FirebasePushCfg {
	var row *sysentity.FirebasePushCfg
	if firebasePushCfgCacheMgr != nil {
		row = GetFirebasePushCfgFromMemory()
	}
	if row == nil {
		row = loadFirebasePushCfgFromDB()
	}
	if row == nil {
		return sysentity.NewDefaultFirebasePushCfg()
	}
	cloned := *row
	cloned.ApplyFirebasePushDefaults()
	return &cloned
}

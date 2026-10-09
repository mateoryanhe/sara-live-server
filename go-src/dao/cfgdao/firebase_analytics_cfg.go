package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	sysentity "xr-game-server/entity/sys"
)

const firebaseAnalyticsCfgCacheKey = "firebase_analytics_cfg"

var firebaseAnalyticsCfgCacheMgr *cache.RowCache[*sysentity.FirebaseAnalyticsCfg]

func InitFirebaseAnalyticsCfgDao() {
	if firebaseAnalyticsCfgCacheMgr != nil {
		return
	}
	firebaseAnalyticsCfgCacheMgr = cache.NewPermanentRowCache[*sysentity.FirebaseAnalyticsCfg]()
}

func EnsureFirebaseAnalyticsCfgDao() {
	InitFirebaseAnalyticsCfgDao()
}

func LoadFirebaseAnalyticsCfgRow() *sysentity.FirebaseAnalyticsCfg {
	return loadFirebaseAnalyticsCfgFromDB()
}

func ResolveFirebaseAnalyticsCfg() *sysentity.FirebaseAnalyticsCfg {
	if firebaseAnalyticsCfgCacheMgr != nil {
		return GetFirebaseAnalyticsCfgCached()
	}
	return loadFirebaseAnalyticsCfgFromDB()
}

func loadFirebaseAnalyticsCfgFromDB() *sysentity.FirebaseAnalyticsCfg {
	var row sysentity.FirebaseAnalyticsCfg
	if err := g.DB().Model(string(sysentity.TbFirebaseAnalyticsCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveFirebaseAnalyticsCfg(row *sysentity.FirebaseAnalyticsCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(sysentity.TbFirebaseAnalyticsCfg)).Save(row)
	return err
}

func ReloadFirebaseAnalyticsCfgCache() *sysentity.FirebaseAnalyticsCfg {
	if firebaseAnalyticsCfgCacheMgr == nil {
		return loadFirebaseAnalyticsCfgFromDB()
	}
	row := loadFirebaseAnalyticsCfgFromDB()
	firebaseAnalyticsCfgCacheMgr.PublishRow(gctx.New(), firebaseAnalyticsCfgCacheKey, row)
	return row
}

func GetFirebaseAnalyticsCfgCached() *sysentity.FirebaseAnalyticsCfg {
	if firebaseAnalyticsCfgCacheMgr == nil {
		return nil
	}
	row, _ := firebaseAnalyticsCfgCacheMgr.GetRowCached(gctx.New(), firebaseAnalyticsCfgCacheKey)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}

package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	sysentity "xr-game-server/entity/sys"
)

const appsFlyerCfgCacheKey = "apps_flyer_cfg"

var appsFlyerCfgCacheMgr *cache.RowCache[*sysentity.AppsFlyerCfg]

func InitAppsFlyerCfgDao() {
	if appsFlyerCfgCacheMgr != nil {
		return
	}
	appsFlyerCfgCacheMgr = cache.NewPermanentRowCache[*sysentity.AppsFlyerCfg]()
}

func loadAppsFlyerCfgFromDB() *sysentity.AppsFlyerCfg {
	var row sysentity.AppsFlyerCfg
	if err := g.DB().Model(string(sysentity.TbAppsFlyerCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveAppsFlyerCfg(row *sysentity.AppsFlyerCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(sysentity.TbAppsFlyerCfg)).Save(row)
	return err
}

func ReloadAppsFlyerCfgCache() *sysentity.AppsFlyerCfg {
	if appsFlyerCfgCacheMgr == nil {
		return loadAppsFlyerCfgFromDB()
	}
	row := loadAppsFlyerCfgFromDB()
	appsFlyerCfgCacheMgr.PublishRow(gctx.New(), appsFlyerCfgCacheKey, row)
	return row
}

// GetAppsFlyerCfgCached App/CMS 仅从内存读取配置.
func GetAppsFlyerCfgCached() *sysentity.AppsFlyerCfg {
	if appsFlyerCfgCacheMgr == nil {
		return nil
	}
	row, _ := appsFlyerCfgCacheMgr.GetRowCached(gctx.New(), appsFlyerCfgCacheKey)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}
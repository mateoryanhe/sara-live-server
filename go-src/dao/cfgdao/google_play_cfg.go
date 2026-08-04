package cfgdao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const googlePlayCfgCacheKey = "google_play_cfg"

var googlePlayCfgCacheMgr *cache.CacheMgr

func InitGooglePlayCfgDao() {
	googlePlayCfgCacheMgr = cache.NewPermanentCacheMgr()
}

func loadGooglePlayCfgFromDB() *entity.GooglePlayCfg {
	var row entity.GooglePlayCfg
	if err := g.DB().Model(string(entity.TbGooglePlayCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

// SaveGooglePlayCfg 保存 Google Play 配置到数据库
func SaveGooglePlayCfg(row *entity.GooglePlayCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbGooglePlayCfg)).Save(row)
	return err
}

// ReloadGooglePlayCfgCache 从数据库重新加载并刷新内存缓存(CMS 保存后调用)
func ReloadGooglePlayCfgCache() *entity.GooglePlayCfg {
	if googlePlayCfgCacheMgr == nil {
		return loadGooglePlayCfgFromDB()
	}
	row := loadGooglePlayCfgFromDB()
	googlePlayCfgCacheMgr.FlushCache(googlePlayCfgCacheKey, row)
	return row
}

// GetGooglePlayCfgFromMemory 仅从内存读取配置,不访问数据库
func GetGooglePlayCfgFromMemory() *entity.GooglePlayCfg {
	if googlePlayCfgCacheMgr == nil {
		return nil
	}
	v := googlePlayCfgCacheMgr.GetFromCache(googlePlayCfgCacheKey)
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.GooglePlayCfg)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}

// GetGooglePlayCfgCached 获取 Google Play 配置(仅读内存,等价于 GetGooglePlayCfgFromMemory)
func GetGooglePlayCfgCached() *entity.GooglePlayCfg {
	return GetGooglePlayCfgFromMemory()
}

func GooglePlayEnabled() bool {
	row := GetGooglePlayCfgFromMemory()
	if row == nil || !row.Enabled {
		return false
	}
	return strings.TrimSpace(row.PackageName) != "" && strings.TrimSpace(row.ServiceAccountJson) != ""
}

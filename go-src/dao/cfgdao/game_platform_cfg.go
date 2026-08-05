package cfgdao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const gamePlatformCfgCacheKey = "game_platform_cfg"

var gamePlatformCfgCacheMgr *cache.CacheMgr

func InitGamePlatformCfgDao() {
	gamePlatformCfgCacheMgr = cache.NewPermanentCacheMgr()
}

func loadGamePlatformCfgFromDB() *entity.GamePlatformCfg {
	var row entity.GamePlatformCfg
	if err := g.DB().Model(string(entity.TbGamePlatformCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveGamePlatformCfg(row *entity.GamePlatformCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbGamePlatformCfg)).Save(row)
	return err
}

func ReloadGamePlatformCfgCache() *entity.GamePlatformCfg {
	if gamePlatformCfgCacheMgr == nil {
		return loadGamePlatformCfgFromDB()
	}
	row := loadGamePlatformCfgFromDB()
	gamePlatformCfgCacheMgr.FlushCache(gamePlatformCfgCacheKey, row)
	return row
}

func GetGamePlatformCfgFromMemory() *entity.GamePlatformCfg {
	if gamePlatformCfgCacheMgr == nil {
		return nil
	}
	v := gamePlatformCfgCacheMgr.GetFromCache(gamePlatformCfgCacheKey)
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.GamePlatformCfg)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}

// GetVendorUrlFromMemory 从内存读取厂家 API 根地址.
func GetVendorUrlFromMemory() string {
	row := GetGamePlatformCfgFromMemory()
	if row == nil {
		return entity.GamePlatformDefaultVendorUrl
	}
	vendorUrl := strings.TrimRight(strings.TrimSpace(row.VendorUrl), "/")
	if vendorUrl == "" {
		return entity.GamePlatformDefaultVendorUrl
	}
	return vendorUrl
}

func GamePlatformCfgReady() bool {
	row := GetGamePlatformCfgFromMemory()
	if row == nil {
		return false
	}
	return strings.TrimSpace(row.Token) != "" &&
		strings.TrimSpace(row.SecretKey) != ""
}

// GetGameIconBaseUrlFromMemory 游戏封面 CDN 根地址(用于拼接第三方返回的相对路径).
func GetGameIconBaseUrlFromMemory() string {
	row := GetGamePlatformCfgFromMemory()
	if row == nil {
		return ""
	}
	return strings.TrimRight(strings.TrimSpace(row.IconUrl), "/")
}

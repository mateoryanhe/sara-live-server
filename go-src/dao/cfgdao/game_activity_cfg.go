package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	sysentity "xr-game-server/entity/sys"
)

const gameActivityCfgCacheKey = "game_activity_cfg"

var gameActivityCfgCacheMgr *cache.RowCache[*sysentity.GameActivityCfg]

func InitGameActivityCfgDao() {
	if gameActivityCfgCacheMgr != nil {
		return
	}
	gameActivityCfgCacheMgr = cache.NewPermanentRowCache[*sysentity.GameActivityCfg]()
}

func EnsureGameActivityCfgDao() {
	InitGameActivityCfgDao()
}

func loadGameActivityCfgFromDB() *sysentity.GameActivityCfg {
	var row sysentity.GameActivityCfg
	if err := g.DB().Model(string(sysentity.TbGameActivityCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveGameActivityCfg(row *sysentity.GameActivityCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(sysentity.TbGameActivityCfg)).Save(row)
	return err
}

func ReloadGameActivityCfgCache() *sysentity.GameActivityCfg {
	if gameActivityCfgCacheMgr == nil {
		return loadGameActivityCfgFromDB()
	}
	row := loadGameActivityCfgFromDB()
	gameActivityCfgCacheMgr.PublishRow(gctx.New(), gameActivityCfgCacheKey, row)
	return row
}

func GetGameActivityCfgFromMemory() *sysentity.GameActivityCfg {
	if gameActivityCfgCacheMgr == nil {
		return nil
	}
	row, _ := gameActivityCfgCacheMgr.GetRowCached(gctx.New(), gameActivityCfgCacheKey)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}

// ResolveGameActivityCfg 优先读缓存;无记录返回 nil.
func ResolveGameActivityCfg() *sysentity.GameActivityCfg {
	if gameActivityCfgCacheMgr != nil {
		if row := GetGameActivityCfgFromMemory(); row != nil {
			return row
		}
	}
	return loadGameActivityCfgFromDB()
}

package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	sysentity "xr-game-server/entity/sys"
)

const dbBackupCfgCacheKey = "db_backup_cfg"

var dbBackupCfgCacheMgr *cache.RowCache[*sysentity.DbBackupCfg]

func InitDbBackupCfgDao() {
	dbBackupCfgCacheMgr = cache.NewPermanentRowCache[*sysentity.DbBackupCfg]()
}

func loadDbBackupCfgFromDB() *sysentity.DbBackupCfg {
	var row sysentity.DbBackupCfg
	if err := g.DB().Model(string(sysentity.TbDbBackupCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveDbBackupCfg(row *sysentity.DbBackupCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(sysentity.TbDbBackupCfg)).Save(row)
	return err
}

func ReloadDbBackupCfgCache() *sysentity.DbBackupCfg {
	if dbBackupCfgCacheMgr == nil {
		return loadDbBackupCfgFromDB()
	}
	row := loadDbBackupCfgFromDB()
	dbBackupCfgCacheMgr.PublishRow(gctx.New(), dbBackupCfgCacheKey, row)
	return row
}

func GetDbBackupCfgCached() *sysentity.DbBackupCfg {
	if dbBackupCfgCacheMgr == nil {
		return nil
	}
	v, _ := dbBackupCfgCacheMgr.GetRowCached(gctx.New(), dbBackupCfgCacheKey)
	if v == nil || v.ID == 0 {
		return nil
	}
	return v
}

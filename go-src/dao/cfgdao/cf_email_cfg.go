package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	sysentity "xr-game-server/entity/sys"
)

const cfEmailCfgCacheKey = "cf_email_cfg"

var cfEmailCfgCacheMgr *cache.RowCache[*sysentity.CfEmailCfg]

func InitCfEmailCfgDao() {
	cfEmailCfgCacheMgr = cache.NewPermanentRowCache[*sysentity.CfEmailCfg]()
}

func loadCfEmailCfgFromDB() *sysentity.CfEmailCfg {
	var row sysentity.CfEmailCfg
	if err := g.DB().Model(string(sysentity.TbCfEmailCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveCfEmailCfg(row *sysentity.CfEmailCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(sysentity.TbCfEmailCfg)).Save(row)
	return err
}

func ReloadCfEmailCfgCache() *sysentity.CfEmailCfg {
	if cfEmailCfgCacheMgr == nil {
		return loadCfEmailCfgFromDB()
	}
	row := loadCfEmailCfgFromDB()
	cfEmailCfgCacheMgr.PublishRow(gctx.New(), cfEmailCfgCacheKey, row)
	return row
}

// GetCfEmailCfgCached 仅从内存读取配置
func GetCfEmailCfgCached() *sysentity.CfEmailCfg {
	if cfEmailCfgCacheMgr == nil {
		return nil
	}
	v, _ := cfEmailCfgCacheMgr.GetRowCached(gctx.New(), cfEmailCfgCacheKey)
	if v == nil || v.ID == 0 {
		return nil
	}
	return v
}

func CfEmailEnabled() bool {
	row := GetCfEmailCfgCached()
	return row != nil &&
		row.Enabled &&
		row.Region != "" &&
		row.AccessKeyId != "" &&
		row.SecretAccessKey != "" &&
		row.FromEmail != ""
}

package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	sysentity "xr-game-server/entity/sys"
)

const attrReportCfgCacheKey = "attr_report_cfg"

var attrReportCfgCacheMgr *cache.RowCache[*sysentity.AttrReportCfg]

func InitAttrReportCfgDao() {
	if attrReportCfgCacheMgr != nil {
		return
	}
	attrReportCfgCacheMgr = cache.NewPermanentRowCache[*sysentity.AttrReportCfg]()
}

func EnsureAttrReportCfgDao() {
	InitAttrReportCfgDao()
}

func loadAttrReportCfgFromDB() *sysentity.AttrReportCfg {
	var row sysentity.AttrReportCfg
	if err := g.DB().Model(string(sysentity.TbAttrReportCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveAttrReportCfg(row *sysentity.AttrReportCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(sysentity.TbAttrReportCfg)).Save(row)
	return err
}

func ReloadAttrReportCfgCache() *sysentity.AttrReportCfg {
	if attrReportCfgCacheMgr == nil {
		return loadAttrReportCfgFromDB()
	}
	row := loadAttrReportCfgFromDB()
	attrReportCfgCacheMgr.PublishRow(gctx.New(), attrReportCfgCacheKey, row)
	return row
}

func GetAttrReportCfgFromMemory() *sysentity.AttrReportCfg {
	if attrReportCfgCacheMgr == nil {
		return nil
	}
	row, _ := attrReportCfgCacheMgr.GetRowCached(gctx.New(), attrReportCfgCacheKey)
	if row == nil || row.ID == 0 {
		return nil
	}
	return row
}

// ResolveAttrReportCfg 优先读缓存;未初始化缓存时直读库.
func ResolveAttrReportCfg() *sysentity.AttrReportCfg {
	if attrReportCfgCacheMgr != nil {
		if row := GetAttrReportCfgFromMemory(); row != nil {
			return row
		}
	}
	return loadAttrReportCfgFromDB()
}

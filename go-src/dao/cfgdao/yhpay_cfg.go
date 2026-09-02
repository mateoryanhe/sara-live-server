package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/recharge"
)

const yhPayCfgCacheKey = "yhpay_cfg"

var yhPayCfgCacheMgr *cache.RowCache[*entity.YhPayCfg]

func InitYhPayCfgDao() {
	yhPayCfgCacheMgr = cache.NewPermanentRowCache[*entity.YhPayCfg]()
}

func loadYhPayCfgFromDB() *entity.YhPayCfg {
	var row entity.YhPayCfg
	if err := g.DB().Model(string(entity.TbYhPayCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveYhPayCfg(row *entity.YhPayCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbYhPayCfg)).Save(row)
	return err
}

func ReloadYhPayCfgCache() *entity.YhPayCfg {
	if yhPayCfgCacheMgr == nil {
		return loadYhPayCfgFromDB()
	}
	row := loadYhPayCfgFromDB()
	yhPayCfgCacheMgr.PublishRow(gctx.New(), yhPayCfgCacheKey, row)
	return row
}

// GetYhPayCfgCached 仅从内存读取配置
func GetYhPayCfgCached() *entity.YhPayCfg {
	if yhPayCfgCacheMgr == nil {
		return nil
	}
	v, _ := yhPayCfgCacheMgr.GetRowCached(gctx.New(), yhPayCfgCacheKey)
	if v == nil || v.ID == 0 {
		return nil
	}
	return v
}

func YhPayEnabled() bool {
	row := GetYhPayCfgCached()
	return row != nil && row.Enabled && row.MerchantCode != "" && row.ApiKey != "" && row.ApiHost != ""
}

package cfgdao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/recharge"
)

const haiPayCfgCacheKey = "haipay_cfg"

var haiPayCfgCacheMgr *cache.RowCache[*entity.HaiPayCfg]

func InitHaiPayCfgDao() {
	haiPayCfgCacheMgr = cache.NewPermanentRowCache[*entity.HaiPayCfg]()
}

func loadHaiPayCfgFromDB() *entity.HaiPayCfg {
	var row entity.HaiPayCfg
	if err := g.DB().Model(string(entity.TbHaiPayCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveHaiPayCfg(row *entity.HaiPayCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbHaiPayCfg)).Save(row)
	return err
}

func ReloadHaiPayCfgCache() *entity.HaiPayCfg {
	if haiPayCfgCacheMgr == nil {
		return loadHaiPayCfgFromDB()
	}
	row := loadHaiPayCfgFromDB()
	haiPayCfgCacheMgr.PublishRow(gctx.New(), haiPayCfgCacheKey, row)
	return row
}

func GetHaiPayCfgCached() *entity.HaiPayCfg {
	if haiPayCfgCacheMgr == nil {
		return nil
	}
	v, _ := haiPayCfgCacheMgr.GetRowCached(gctx.New(), haiPayCfgCacheKey)
	if v == nil || v.ID == 0 {
		return nil
	}
	return v
}

func HaiPayEnabled() bool {
	row := GetHaiPayCfgCached()
	if row == nil || !row.Enabled {
		return false
	}
	if row.AppId <= 0 {
		return false
	}
	if strings.TrimSpace(row.ApiHost) == "" {
		return false
	}
	if strings.TrimSpace(row.MerchantSecretKey) == "" {
		return false
	}
	if strings.TrimSpace(row.MerchantPrivateKey) == "" {
		return false
	}
	if strings.TrimSpace(row.HaiPayPublicKey) == "" {
		return false
	}
	return true
}

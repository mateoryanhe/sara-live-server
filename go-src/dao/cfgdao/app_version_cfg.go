package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/cms"
)

func LoadAppVersionCfg() *entity.AppVersionCfg {
	var row entity.AppVersionCfg
	if err := g.DB().Model(string(entity.TbAppVersionCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveAppVersionCfg(row *entity.AppVersionCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbAppVersionCfg)).Save(row)
	return err
}

package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/activity"
)

func LoadFirstRechargeActivityCfg() *activity.FirstRechargeActivityCfg {
	var row activity.FirstRechargeActivityCfg
	if err := g.DB().Model(string(activity.TbFirstRechargeActivityCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveFirstRechargeActivityCfg(row *activity.FirstRechargeActivityCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(activity.TbFirstRechargeActivityCfg)).Save(row)
	return err
}

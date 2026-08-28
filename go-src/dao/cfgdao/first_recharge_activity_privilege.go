package cfgdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/snowflake"
	"xr-game-server/entity/activity"
)

func LoadFirstRechargeActivityPrivileges() []*activity.FirstRechargeActivityPrivilege {
	var rows []*activity.FirstRechargeActivityPrivilege
	if err := g.DB().Model(string(activity.TbFirstRechargeActivityPrivilege)).
		Order("sort asc, id asc").
		Scan(&rows); err != nil {
		return nil
	}
	return rows
}

// ReplaceAllFirstRechargeActivityPrivileges 物理清空后全量写入特权列表
func ReplaceAllFirstRechargeActivityPrivileges(rows []*activity.FirstRechargeActivityPrivilege) error {
	if _, err := g.DB().Model(string(activity.TbFirstRechargeActivityPrivilege)).Where("id > ?", 0).Delete(); err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	now := time.Now()
	for _, row := range rows {
		if row == nil {
			continue
		}
		if row.ID == 0 {
			row.ID = snowflake.GetId()
		}
		if row.CreatedAt.IsZero() {
			row.CreatedAt = now
		}
		row.UpdatedAt = now
	}
	_, err := g.DB().Model(string(activity.TbFirstRechargeActivityPrivilege)).Data(rows).Batch(len(rows)).Insert()
	return err
}

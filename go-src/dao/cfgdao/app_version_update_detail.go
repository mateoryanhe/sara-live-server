package cfgdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/snowflake"
	"xr-game-server/entity/cms"
)

func LoadAppVersionUpdateDetails() []*entity.AppVersionUpdateDetail {
	var rows []*entity.AppVersionUpdateDetail
	if err := g.DB().Model(string(entity.TbAppVersionUpdateDetail)).
		Order("sort asc, id asc").
		Scan(&rows); err != nil {
		return nil
	}
	return rows
}

// ReplaceAllAppVersionUpdateDetails 物理清空后全量写入更新明细
func ReplaceAllAppVersionUpdateDetails(rows []*entity.AppVersionUpdateDetail) error {
	if _, err := g.DB().Model(string(entity.TbAppVersionUpdateDetail)).Where("id > ?", 0).Delete(); err != nil {
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
	_, err := g.DB().Model(string(entity.TbAppVersionUpdateDetail)).Data(rows).Batch(len(rows)).Insert()
	return err
}

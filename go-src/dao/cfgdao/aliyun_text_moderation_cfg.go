package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/live"
)

func LoadAliyunTextModerationCfg() *entity.AliyunTextModerationCfg {
	var row entity.AliyunTextModerationCfg
	if err := g.DB().Model(string(entity.TbAliyunTextModerationCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveAliyunTextModerationCfg(row *entity.AliyunTextModerationCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbAliyunTextModerationCfg)).Save(row)
	return err
}

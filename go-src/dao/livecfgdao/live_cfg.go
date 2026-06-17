package livecfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

// Load 从数据库加载直播配置(通常仅一条)
func Load() *entity.LiveCfg {
	var row entity.LiveCfg
	if err := g.DB().Model(string(entity.TbLiveCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

// Save 保存直播配置
func Save(row *entity.LiveCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbLiveCfg)).Save(row)
	return err
}

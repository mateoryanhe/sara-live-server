package agoracfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

// Load 从数据库加载声网配置(通常仅一条)
func Load() *entity.AgoraCfg {
	var row entity.AgoraCfg
	if err := g.DB().Model(string(entity.TbAgoraCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

// Save 保存声网配置
func Save(row *entity.AgoraCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbAgoraCfg)).Save(row)
	return err
}

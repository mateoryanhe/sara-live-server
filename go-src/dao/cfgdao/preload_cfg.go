package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

const DefaultRecentLoginPreloadLimit = 100

func LoadPreloadCfg() *entity.PreloadCfg {
	var row entity.PreloadCfg
	if err := g.DB().Model(string(entity.TbPreloadCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SavePreloadCfg(row *entity.PreloadCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbPreloadCfg)).Save(row)
	return err
}

// GetRecentLoginPreloadLimit 读取最近登录用户预热数量,未配置或无效时返回默认值
func GetRecentLoginPreloadLimit() int {
	cfg := LoadPreloadCfg()
	if cfg == nil || cfg.RecentLoginLimit <= 0 {
		return DefaultRecentLoginPreloadLimit
	}
	return cfg.RecentLoginLimit
}

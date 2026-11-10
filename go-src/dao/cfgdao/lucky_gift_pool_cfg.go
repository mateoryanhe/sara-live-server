package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/live"
)

func LoadLuckyGiftPoolCfg() *entity.LuckyGiftPoolCfg {
	var row entity.LuckyGiftPoolCfg
	if err := g.DB().Model(string(entity.TbLuckyGiftPoolCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

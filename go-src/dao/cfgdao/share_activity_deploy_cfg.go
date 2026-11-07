package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/sys"
)

func LoadShareActivityDeployCfg() *entity.ShareActivityDeployCfg {
	var row entity.ShareActivityDeployCfg
	if err := g.DB().Model(string(entity.TbShareActivityDeployCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveShareActivityDeployCfg(row *entity.ShareActivityDeployCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbShareActivityDeployCfg)).Save(row)
	return err
}

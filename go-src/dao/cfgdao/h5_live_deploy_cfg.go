package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/sys"
)

func LoadH5LiveDeployCfg() *entity.H5LiveDeployCfg {
	var row entity.H5LiveDeployCfg
	if err := g.DB().Model(string(entity.TbH5LiveDeployCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveH5LiveDeployCfg(row *entity.H5LiveDeployCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbH5LiveDeployCfg)).Save(row)
	return err
}

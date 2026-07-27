package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

func LoadAccountCfg() *entity.AccountCfg {
	var row entity.AccountCfg
	if err := g.DB().Model(string(entity.TbAccountCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveAccountCfg(row *entity.AccountCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbAccountCfg)).Save(row)
	return err
}

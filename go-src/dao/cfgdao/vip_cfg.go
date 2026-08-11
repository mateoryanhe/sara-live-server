package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

func CreateVipCfg(row *entity.VipCfg) error {
	_, err := g.DB().Model(string(entity.TbVipCfg)).Save(row)
	return err
}

func UpdateVipCfg(row *entity.VipCfg) error {
	return CreateVipCfg(row)
}

func DeleteVipCfg(id uint64) error {
	_, err := g.DB().Model(string(entity.TbVipCfg)).WherePri(id).Delete()
	return err
}

func GetAllVipCfg() []*entity.VipCfg {
	var rows []*entity.VipCfg
	_ = g.DB().Model(string(entity.TbVipCfg)).Order("level asc").Scan(&rows)
	return rows
}

func GetVipCfgsByIDs(ids []uint64) []*entity.VipCfg {
	if len(ids) == 0 {
		return nil
	}
	var rows []*entity.VipCfg
	_ = g.DB().Model(string(entity.TbVipCfg)).WhereIn("id", ids).Scan(&rows)
	return rows
}

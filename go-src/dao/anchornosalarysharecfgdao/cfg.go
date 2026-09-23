package anchornosalarysharecfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/live"
)

func GetFirst() *entity.AnchorNoSalaryShareCfg {
	var row entity.AnchorNoSalaryShareCfg
	if err := g.DB().Model(string(entity.TbAnchorNoSalaryShareCfg)).Order("id asc").Limit(1).Scan(&row); err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func Save(row *entity.AnchorNoSalaryShareCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbAnchorNoSalaryShareCfg)).Save(row)
	return err
}

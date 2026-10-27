package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/sys"
)

func LoadAuditCodeCfg() *entity.AuditCodeCfg {
	var row entity.AuditCodeCfg
	if err := g.DB().Model(string(entity.TbAuditCodeCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveAuditCodeCfg(row *entity.AuditCodeCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbAuditCodeCfg)).Save(row)
	return err
}

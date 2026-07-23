package privacypolicycfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

func Load() *entity.PrivacyPolicyCfg {
	var row entity.PrivacyPolicyCfg
	if err := g.DB().Model(string(entity.TbPrivacyPolicyCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func Save(row *entity.PrivacyPolicyCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbPrivacyPolicyCfg)).Save(row)
	return err
}

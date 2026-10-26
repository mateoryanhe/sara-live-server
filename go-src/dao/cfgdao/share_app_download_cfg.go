package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/user"
)

func LoadShareAppDownloadCfg() *entity.ShareAppDownloadCfg {
	var row entity.ShareAppDownloadCfg
	if err := g.DB().Model(string(entity.TbShareAppDownloadCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveShareAppDownloadCfg(row *entity.ShareAppDownloadCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbShareAppDownloadCfg)).Save(row)
	return err
}

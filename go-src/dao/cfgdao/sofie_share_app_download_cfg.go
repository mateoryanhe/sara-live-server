package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/user"
)

// LoadSofieShareAppDownloadCfg 单行配置(通常仅一条).
func LoadSofieShareAppDownloadCfg() *entity.SofieShareAppDownloadCfg {
	var row entity.SofieShareAppDownloadCfg
	if err := g.DB().Model(string(entity.TbSofieShareAppDownloadCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveSofieShareAppDownloadCfg(row *entity.SofieShareAppDownloadCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbSofieShareAppDownloadCfg)).Save(row)
	return err
}
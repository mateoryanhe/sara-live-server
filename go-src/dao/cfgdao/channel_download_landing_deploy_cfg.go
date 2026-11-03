package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/sys"
)

func LoadChannelDownloadLandingDeployCfg() *entity.ChannelDownloadLandingDeployCfg {
	var row entity.ChannelDownloadLandingDeployCfg
	if err := g.DB().Model(string(entity.TbChannelDownloadLandingDeployCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveChannelDownloadLandingDeployCfg(row *entity.ChannelDownloadLandingDeployCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbChannelDownloadLandingDeployCfg)).Save(row)
	return err
}

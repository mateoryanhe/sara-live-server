package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/entity/sys"
)

// GetDataSyncCfg 从数据库读取数据同步配置(不缓存)
func GetDataSyncCfg() *entity.DataSyncCfg {
	var row entity.DataSyncCfg
	if err := g.DB().Model(string(entity.TbDataSyncCfg)).Order(string(db.IdName) + " asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

// SaveDataSyncCfg 保存数据同步配置到数据库
func SaveDataSyncCfg(row *entity.DataSyncCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbDataSyncCfg)).Save(row)
	return err
}

package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/user"
)

func LoadLiveRoomShareCfg() *entity.LiveRoomShareCfg {
	var row entity.LiveRoomShareCfg
	if err := g.DB().Model(string(entity.TbLiveRoomShareCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveLiveRoomShareCfg(row *entity.LiveRoomShareCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbLiveRoomShareCfg)).Save(row)
	return err
}
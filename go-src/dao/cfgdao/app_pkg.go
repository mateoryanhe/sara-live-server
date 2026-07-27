package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

func CreateAppPkg(row *entity.AppPkg) error {
	_, err := g.DB().Model(string(entity.TbAppPkg)).Save(row)
	return err
}

func UpdateAppPkg(row *entity.AppPkg) error {
	return CreateAppPkg(row)
}

func DeleteAppPkg(id uint64) error {
	_, err := g.DB().Model(string(entity.TbAppPkg)).WherePri(id).Delete()
	return err
}

func GetAllAppPkg() []*entity.AppPkg {
	var rows []*entity.AppPkg
	_ = g.DB().Model(string(entity.TbAppPkg)).Order("id desc").Scan(&rows)
	return rows
}

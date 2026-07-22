package apppkgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

func Create(row *entity.AppPkg) error {
	_, err := g.DB().Model(string(entity.TbAppPkg)).Save(row)
	return err
}

func Update(row *entity.AppPkg) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbAppPkg)).WherePri(id).Delete()
	return err
}

func GetAll() []*entity.AppPkg {
	var rows []*entity.AppPkg
	_ = g.DB().Model(string(entity.TbAppPkg)).Order("id desc").Scan(&rows)
	return rows
}

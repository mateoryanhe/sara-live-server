package cmsuserdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

func SavePermissions(data []*entity.Permission) {
	DeleteRolePermissions(data[0].RoleId)
	g.DB().Model(string(entity.TbPermission)).Save(data)
}

// DeleteRolePermissions 删除角色关联的权限
func DeleteRolePermissions(roleId uint64) error {
	_, err := g.DB().Model(string(entity.TbPermission)).Delete("role_id = ?", roleId)
	if err == nil {
		return nil
	}
	return err
}

func GetGetPermissionList(roleId uint64) (list []*entity.Permission) {
	list = []*entity.Permission{}
	g.DB().Model(string(entity.TbPermission)).Where("role_id = ?", roleId).Scan(&list)
	return list
}

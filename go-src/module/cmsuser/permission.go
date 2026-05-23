package cmsuser

import (
	"context"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dto/permdto"
	"xr-game-server/entity"
)

// GetPermissionList 获取权限列表
func GetPermissionList(ctx context.Context, req *permdto.PermissionListReq) (list []*entity.Permission, err error) {
	permissions := cmsuserdao.GetGetPermissionList(req.RoleId)
	return permissions, nil
}

// CreatePermission 创建权限
func CreatePermission(ctx context.Context, req *permdto.CreatePermissionReq) (res bool, err error) {
	// 检查权限名称是否已存在
	cmsuserdao.SavePermissions(req.Data)
	return true, nil
}

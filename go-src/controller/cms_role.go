package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/cmsroledto"
	"xr-game-server/dto/permdto"
	"xr-game-server/entity"
	"xr-game-server/module/cmsuser"
)

const (
	RoleUrl = "/role"
)

type RoleController struct {
}

func initRoleController() {
	httpserver.RegCMS(RoleUrl, &RoleController{})
}

// RoleList 获取角色列表
func (r *RoleController) RoleList(ctx context.Context, req *cmsroledto.RoleListReq) (res *httpserver.CMSQueryResp, err error) {
	return cmsuser.GetRoleList(ctx, req)
}

// CreateRole 创建角色
func (r *RoleController) CreateRole(ctx context.Context, req *cmsroledto.CreateRoleReq) (res *cmsroledto.CreateRoleRes, err error) {
	return cmsuser.CreateRole(ctx, req)
}

// UpdateRole 更新角色
func (r *RoleController) UpdateRole(ctx context.Context, req *cmsroledto.UpdateRoleReq) (res *cmsroledto.UpdateRoleRes, err error) {
	return cmsuser.UpdateRole(ctx, req)
}

// DeleteRole 删除角色
func (r *RoleController) DeleteRole(ctx context.Context, req *cmsroledto.DeleteRoleReq) (res *cmsroledto.DeleteRoleRes, err error) {
	return cmsuser.DeleteRole(ctx, req)
}

// GetRolePermissions 获取角色权限
func (r *RoleController) GetRolePermissions(ctx context.Context, req *permdto.PermissionListReq) (res []*entity.Permission, err error) {
	return cmsuser.GetPermissionList(ctx, req)
}

// SetRolePermissions 设置角色权限
func (r *RoleController) CreatePermission(ctx context.Context, req *permdto.CreatePermissionReq) (res bool, err error) {
	return cmsuser.CreatePermission(ctx, req)
}

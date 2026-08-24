package cmsuser

import (
	"context"
	"errors"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dto/cmsroledto"
	"xr-game-server/entity/cms"
)

// GetRoleList 获取角色列表
func GetRoleList(ctx context.Context, req *cmsroledto.RoleListReq) (res *httpserver.CMSQueryResp, err error) {
	total, roles := cmsuserdao.GetRoleList(req)

	return &httpserver.CMSQueryResp{
		Total: total,
		Data:  roles,
	}, nil
}

// CreateRole 创建角色
func CreateRole(ctx context.Context, req *cmsroledto.CreateRoleReq) (res *cmsroledto.CreateRoleRes, err error) {
	// 检查角色名称是否已存在
	existingRole := cmsuserdao.GetRoleByName(req.Name)
	if existingRole != nil {
		return nil, errors.New("角色名称已存在")
	}

	roleType, err := normalizeRoleType(req.RoleType)
	if err != nil {
		return nil, err
	}

	role := entity.CMSRole{
		Name:        req.Name,
		Description: req.Description,
		RoleType:    roleType,
		Status:      req.Status,
	}

	err = cmsuserdao.CreateRole(&role)
	if err != nil {
		return nil, err
	}

	return &cmsroledto.CreateRoleRes{
		ID: role.ID,
	}, nil
}

// UpdateRole 更新角色
func UpdateRole(ctx context.Context, req *cmsroledto.UpdateRoleReq) (res *cmsroledto.UpdateRoleRes, err error) {
	role := cmsuserdao.GetRoleById(req.ID)
	if role == nil {
		return nil, errors.New("角色不存在")
	}

	// 检查角色名称是否与其他角色重复
	existingRole := cmsuserdao.GetRoleByName(req.Name)
	if existingRole != nil && existingRole.ID != req.ID {
		return nil, errors.New("角色名称已存在")
	}

	roleType, err := normalizeRoleType(req.RoleType)
	if err != nil {
		return nil, err
	}

	role.Name = req.Name
	role.Description = req.Description
	role.RoleType = roleType
	role.Status = req.Status

	err = cmsuserdao.UpdateRole(role)
	if err != nil {
		return nil, err
	}
	if req.Status == 0 {
		invalidateCmsTokensByRoleId(req.ID)
	}

	return &cmsroledto.UpdateRoleRes{
		Success: true,
	}, nil
}

// DeleteRole 删除角色
func DeleteRole(ctx context.Context, req *cmsroledto.DeleteRoleReq) (res *cmsroledto.DeleteRoleRes, err error) {
	invalidateCmsTokensByRoleId(req.ID)

	// 删除角色关联的权限
	err = cmsuserdao.DeleteRolePermissions(req.ID)
	if err != nil {
		return nil, err
	}

	// 删除角色
	err = cmsuserdao.DeleteRole(req.ID)
	if err != nil {
		return nil, err
	}

	return &cmsroledto.DeleteRoleRes{
		Success: true,
	}, nil
}

func normalizeRoleType(roleType uint8) (uint8, error) {
	if roleType == 0 {
		return entity.CMSRoleTypeInternal, nil
	}
	if roleType != entity.CMSRoleTypeInternal && roleType != entity.CMSRoleTypeExternal {
		return 0, errors.New("角色类型无效")
	}
	return roleType, nil
}

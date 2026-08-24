package cmsuser

import (
	"context"
	"errors"
	"xr-game-server/constants/common"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dto/cmsuserdto"
	"xr-game-server/entity/cms"
)

// GetCMSUserList 获取CMS用户列表
func GetCMSUserList(ctx context.Context, req *cmsuserdto.CMSUserListReq) (res *httpserver.CMSQueryResp, err error) {
	total, users := cmsuserdao.GetCMSUserList(req)

	return &httpserver.CMSQueryResp{
		Total: total,
		Data:  users,
	}, nil
}

// CreateCMSUser 创建CMS用户
func CreateCMSUser(ctx context.Context, req *cmsuserdto.CreateCMSUserReq) (res *cmsuserdto.CreateCMSUserRes, err error) {
	// 检查CMS用户名称是否已存在
	existingUser := cmsuserdao.GetCMSUserByName(req.Name)
	if existingUser != nil {
		return nil, errors.New("CMS用户名称已存在")
	}

	user := entity.CMSUser{
		Name:   req.Name,
		Pwd:    req.Pwd,
		Status: req.Status,
		Admin:  req.Admin,
		RoleId: req.RoleId,
		Remark: req.Remark,
	}

	err = cmsuserdao.CreateCMSUser(&user)
	if err != nil {
		return nil, err
	}

	return &cmsuserdto.CreateCMSUserRes{
		ID: user.ID,
	}, nil
}

// CreateGuildCMSUser 工会模块创建外部 CMS 用户(固定非管理员)
func CreateGuildCMSUser(_ context.Context, req *cmsuserdto.CreateGuildCMSUserReq) (*cmsuserdto.CreateCMSUserRes, error) {
	existingUser := cmsuserdao.GetCMSUserByName(req.Name)
	if existingUser != nil {
		return nil, errors.New("CMS用户名称已存在")
	}

	role := cmsuserdao.GetRoleById(req.RoleId)
	if role == nil {
		return nil, errors.New("角色不存在")
	}
	if role.Status != 1 {
		return nil, errors.New("角色已禁用")
	}
	if role.RoleType != entity.CMSRoleTypeExternal {
		return nil, errors.New("只能选择外部类型角色")
	}

	status := req.Status
	if status == 0 {
		status = 1
	}

	user := entity.CMSUser{
		Name:   req.Name,
		Pwd:    req.Pwd,
		Status: status,
		Admin:  false,
		RoleId: req.RoleId,
		Remark: req.Remark,
	}

	if err := cmsuserdao.CreateCMSUser(&user); err != nil {
		return nil, err
	}

	return &cmsuserdto.CreateCMSUserRes{
		ID: user.ID,
	}, nil
}

// UpdateCMSUser 更新CMS用户
func UpdateCMSUser(ctx context.Context, req *cmsuserdto.UpdateCMSUserReq) (res *cmsuserdto.UpdateCMSUserRes, err error) {
	user := cmsuserdao.LoadCMSUserFromDB(req.ID)
	if user == nil {
		return nil, errors.New("CMS用户不存在")
	}

	oldStatus := user.Status
	oldName := user.Name
	oldRoleId := user.RoleId
	operatorId := httpserver.GetAuthId(ctx)
	if oldStatus == common.True && req.Status == common.False && user.ID == operatorId {
		return nil, errors.New("不能停用自己的账号")
	}

	// 检查CMS用户名称是否与其他用户重复
	existingUser := cmsuserdao.GetCMSUserByName(req.Name)
	if existingUser != nil && existingUser.ID != req.ID {
		return nil, errors.New("CMS用户名称已存在")
	}

	user.Name = req.Name
	pwdChanged := req.Pwd != ""
	if pwdChanged {
		user.Pwd = req.Pwd
	}
	user.Status = req.Status
	user.Admin = req.Admin
	user.RoleId = req.RoleId
	user.Remark = req.Remark

	err = cmsuserdao.UpdateCMSUser(user)
	if err != nil {
		return nil, err
	}
	if oldName != user.Name {
		cmsuserdao.RemoveCMSLoginUserMissCache(oldName)
	}
	if oldStatus == common.True && req.Status == common.False {
		invalidateCmsToken(user.ID)
	}
	if pwdChanged || user.RoleId != oldRoleId {
		invalidateCmsToken(user.ID)
	}

	return &cmsuserdto.UpdateCMSUserRes{
		Success: true,
	}, nil
}

// DeleteCMSUser 删除CMS用户
func DeleteCMSUser(ctx context.Context, req *cmsuserdto.DeleteCMSUserReq) (res *cmsuserdto.DeleteCMSUserRes, err error) {
	// 检查用户是否存在
	user := cmsuserdao.GetCMSUserById(req.ID)
	if user == nil {
		return nil, errors.New("CMS用户不存在")
	}
	if user.ID == httpserver.GetAuthId(ctx) {
		return nil, errors.New("不能删除自己的账号")
	}

	invalidateCmsToken(req.ID)

	// 删除CMS用户
	err = cmsuserdao.DeleteCMSUser(req.ID)
	if err != nil {
		return nil, err
	}

	return &cmsuserdto.DeleteCMSUserRes{
		Success: true,
	}, nil
}

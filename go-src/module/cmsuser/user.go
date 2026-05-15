package cmsuser

import (
	"context"
	"errors"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dto/cmsuserdto"
	"xr-game-server/entity"
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
	}

	err = cmsuserdao.CreateCMSUser(&user)
	if err != nil {
		return nil, err
	}

	return &cmsuserdto.CreateCMSUserRes{
		ID: user.ID,
	}, nil
}

// UpdateCMSUser 更新CMS用户
func UpdateCMSUser(ctx context.Context, req *cmsuserdto.UpdateCMSUserReq) (res *cmsuserdto.UpdateCMSUserRes, err error) {
	user := cmsuserdao.GetCMSUserById(req.ID)
	if user == nil {
		return nil, errors.New("CMS用户不存在")
	}

	// 检查CMS用户名称是否与其他用户重复
	existingUser := cmsuserdao.GetCMSUserByName(req.Name)
	if existingUser != nil && existingUser.ID != req.ID {
		return nil, errors.New("CMS用户名称已存在")
	}

	user.Name = req.Name
	if req.Pwd != "" {
		user.Pwd = req.Pwd
	}
	user.Status = req.Status
	user.Admin = req.Admin
	user.RoleId = req.RoleId

	err = cmsuserdao.UpdateCMSUser(user)
	if err != nil {
		return nil, err
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

	// 删除CMS用户
	err = cmsuserdao.DeleteCMSUser(req.ID)
	if err != nil {
		return nil, err
	}

	return &cmsuserdto.DeleteCMSUserRes{
		Success: true,
	}, nil
}

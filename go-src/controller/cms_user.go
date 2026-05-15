package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/cmsuserdto"
	"xr-game-server/module/cmsuser"
)

const (
	CMSUserUrl = "/cmsuser"
)

type CMSUserController struct {
}

func initCMSUserController() {
	httpserver.RegCMS(CMSUserUrl, &CMSUserController{})
}

// CMSUserList 获取CMS用户列表
func (c *CMSUserController) CMSUserList(ctx context.Context, req *cmsuserdto.CMSUserListReq) (res *httpserver.CMSQueryResp, err error) {
	return cmsuser.GetCMSUserList(ctx, req)
}

// CreateCMSUser 创建CMS用户
func (c *CMSUserController) CreateCMSUser(ctx context.Context, req *cmsuserdto.CreateCMSUserReq) (res *cmsuserdto.CreateCMSUserRes, err error) {
	return cmsuser.CreateCMSUser(ctx, req)
}

// UpdateCMSUser 更新CMS用户
func (c *CMSUserController) UpdateCMSUser(ctx context.Context, req *cmsuserdto.UpdateCMSUserReq) (res *cmsuserdto.UpdateCMSUserRes, err error) {
	return cmsuser.UpdateCMSUser(ctx, req)
}

// DeleteCMSUser 删除CMS用户
func (c *CMSUserController) DeleteCMSUser(ctx context.Context, req *cmsuserdto.DeleteCMSUserReq) (res *cmsuserdto.DeleteCMSUserRes, err error) {
	return cmsuser.DeleteCMSUser(ctx, req)
}

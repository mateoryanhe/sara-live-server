package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/accountdto"
	"xr-game-server/module/userinfo"
)

const (
	AccountUrl = "/account"
)

type AccountController struct {
}

func initAccountController() {
	httpserver.RegCMS(AccountUrl, &AccountController{})
}

func (a *AccountController) Ban(ctx context.Context, req *accountdto.BanReq) (resp *accountdto.BanRes, e error) {
	return userinfo.Ban(ctx, req)
}

func (a *AccountController) UnBan(ctx context.Context, req *accountdto.UnBanReq) (bool, error) {
	return userinfo.UnBan(ctx, req)
}

func (a *AccountController) CancelUser(ctx context.Context, req *accountdto.CancelReq) (bool, error) {
	return userinfo.CancelUser(ctx, req)
}

func (a *AccountController) UnCancelUser(ctx context.Context, req *accountdto.UnCancelReq) (bool, error) {
	return userinfo.UnCancelUser(ctx, req)
}

func (a *AccountController) QueryUserInfo(ctx context.Context, req *accountdto.QueryUserInfoReq) (res *httpserver.CMSQueryResp, err error) {
	return userinfo.QueryUserInfo(ctx, req)
}

func (a *AccountController) QueryAnchorList(ctx context.Context, req *accountdto.QueryAnchorListReq) (res *httpserver.CMSQueryResp, err error) {
	return userinfo.QueryAnchorList(ctx, req)
}

func (a *AccountController) SetAnchor(ctx context.Context, req *accountdto.SetAnchorReq) (*accountdto.SetAnchorRes, error) {
	return userinfo.SetAnchor(ctx, req)
}

package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/accountdto"
	"xr-game-server/module/account"
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
	return account.Ban(ctx, req)
}

func (a *AccountController) UnBan(ctx context.Context, req *accountdto.UnBanReq) (bool, error) {
	return account.UnBan(ctx, req)
}

func (a *AccountController) CancelUser(ctx context.Context, req *accountdto.CancelReq) (bool, error) {
	return account.CancelUser(ctx, req)
}

func (a *AccountController) UnCancelUser(ctx context.Context, req *accountdto.UnCancelReq) (bool, error) {
	return account.UnCancelUser(ctx, req)
}

func (a *AccountController) QueryUserInfo(ctx context.Context, req *accountdto.QueryUserInfoReq) (res *httpserver.CMSQueryResp, err error) {
	return account.QueryUserInfo(ctx, req)
}

func (a *AccountController) SetAnchor(ctx context.Context, req *accountdto.SetAnchorReq) (*accountdto.SetAnchorRes, error) {
	return account.SetAnchor(ctx, req)
}

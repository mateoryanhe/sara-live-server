package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/authdto"
	"xr-game-server/module/auth"
)

type AuthController struct {
}

func initAuthApi() {
	httpserver.RegNonAuthAPI("/auth", new(AuthController))
}

func (a *AuthController) TestLogin(ctx context.Context, req *authdto.TestLoginReq) (resp *authdto.TestLoginRes, err error) {
	return auth.TestLogin(ctx, req)
}
func (a *AuthController) CMSLogin(ctx context.Context, req *authdto.CMSLoginReq) (res *authdto.CMSLoginRes, err error) {
	return auth.CMSLogin(ctx, req)
}

func (s *AuthController) PushAgainLoginReq(ctx context.Context, req *authdto.PushAgainLoginReq) (res *httpserver.PushResp, err error) {
	return nil, nil
}

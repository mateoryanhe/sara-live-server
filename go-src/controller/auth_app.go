package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/authdto"
	"xr-game-server/module/auth"
)

const AuthAppUrl = "/auth"

type AuthAppController struct{}

func initAuthAppController() {
	httpserver.RegAPI(AuthAppUrl, &AuthAppController{})
}

// PhoneChangePassword App端修改登录密码(需旧密码)
func (c *AuthAppController) PhoneChangePassword(ctx context.Context, req *authdto.PhoneChangePasswordReq) (*authdto.PhoneChangePasswordRes, error) {
	return auth.PhoneChangePassword(ctx, req)
}

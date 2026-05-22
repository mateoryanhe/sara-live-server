package controller

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/authdto"
	"xr-game-server/dto/verificationcodedto"
	"xr-game-server/module/auth"
	"xr-game-server/module/verification_code"
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

func (s *AuthController) SendCode(ctx context.Context, req *verificationcodedto.SendCodeReq) (*verificationcodedto.SendCodeRes, error) {
	req.IP = g.RequestFromCtx(ctx).GetClientIp()
	return verification_code.SendCode(ctx, req)
}

func (s *AuthController) PhoneRegister(ctx context.Context, req *authdto.PhoneRegisterReq) (*authdto.PhoneRegisterRes, error) {
	return auth.PhoneRegister(ctx, req)
}

func (s *AuthController) PhoneLogin(ctx context.Context, req *authdto.PhoneLoginReq) (*authdto.PhoneLoginRes, error) {
	return auth.PhoneLogin(ctx, req)
}

func (s *AuthController) PhoneResetPassword(ctx context.Context, req *authdto.PhoneResetPasswordReq) (*authdto.PhoneResetPasswordRes, error) {
	return auth.PhoneResetPassword(ctx, req)
}

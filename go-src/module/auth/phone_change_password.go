package auth

import (
	"context"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
)

// PhoneChangePassword App端修改登录密码(校验旧密码,无需手机验证码)
func PhoneChangePassword(ctx context.Context, req *authdto.PhoneChangePasswordReq) (*authdto.PhoneChangePasswordRes, error) {
	userId := httpserver.GetAuthId(ctx)
	account := accountdao.GetAccountById(userId)
	if account == nil || account.Password == "" {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	if account.Password != gmd5.MustEncryptString(req.OldPassword) {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	account.SetPassword(gmd5.MustEncryptString(req.NewPassword))
	return &authdto.PhoneChangePasswordRes{Success: true}, nil
}

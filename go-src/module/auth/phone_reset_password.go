package auth

import (
	"context"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"xr-game-server/core/phoneutil"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
	"xr-game-server/module/verification_code"
)

func PhoneResetPassword(ctx context.Context, req *authdto.PhoneResetPasswordReq) (res *authdto.PhoneResetPasswordRes, err error) {
	valid, err := verification_code.VerifyCode(req.PhoneAreaCode, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errercode.CreateCode(errercode.VerifyCodeInvalid)
	}

	phoneAreaCode := phoneutil.NormalizeAreaCode(req.PhoneAreaCode)
	account := accountdao.FindAccountBy(req.Phone, PhoneChannel, phoneAreaCode)
	if account == nil || account.Password == "" {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}

	account.SetPassword(gmd5.MustEncryptString(req.Password))

	return &authdto.PhoneResetPasswordRes{Success: true}, nil
}

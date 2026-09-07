package auth

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
)

// BindEmail 当前登录用户绑定邮箱；邮箱须未被未注销账号使用.
func BindEmail(ctx context.Context, req *authdto.BindEmailReq) (*authdto.BindEmailRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	email := normalizeEmailKey(req.Email)
	if email == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := VerifyEmailCode(ctx, email, req.Code); err != nil {
		return nil, err
	}

	lockKey := fmt.Sprintf("email_bind:%s", email)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	ext := userinfodao.GetUserExtByUserId(userId)
	if ext == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	current := normalizeEmailKey(ext.Email)
	if current != "" {
		if current == email {
			return &authdto.BindEmailRes{Success: true}, nil
		}
		return nil, errercode.CreateCode(errercode.EmailAlreadyBound)
	}
	if err := ensureEmailAvailable(email, userId); err != nil {
		return nil, err
	}

	ext.SetEmail(email)
	userinfodao.PublishUserExt(ext)
	return &authdto.BindEmailRes{Success: true}, nil
}

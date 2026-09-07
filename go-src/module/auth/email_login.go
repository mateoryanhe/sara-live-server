package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/event"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

// EmailLogin 邮箱+验证码登录.
// 优先: EmailChannel 账号 → 已绑定该邮箱的未注销账号 → 新建 EmailChannel 账号.
func EmailLogin(ctx context.Context, req *authdto.EmailLoginReq) (*authdto.EmailLoginRes, error) {
	if err := ensureSimulatorLoginAllowed(req.DeviceInfo); err != nil {
		return nil, err
	}
	email := normalizeEmailKey(req.Email)
	if email == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := VerifyEmailCode(ctx, email, req.Code); err != nil {
		return nil, err
	}

	lockKey := fmt.Sprintf("email_login:%s", email)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	account := accountdao.FindActiveAccount(email, EmailChannel)
	isNewUser := false
	if account == nil {
		account = findActiveAccountByBoundEmail(email)
	}
	if account == nil {
		if err := ensureEmailAvailable(email, 0); err != nil {
			return nil, err
		}
		account = accountdao.RegisterAccount(email, EmailChannel)
		isNewUser = true
	}
	if account == nil || account.ID == 0 {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	if account.Ban && account.BanApplyTime != nil && account.BanApplyTime.After(time.Now()) {
		return nil, errercode.CreateCode(errercode.Ban)
	}

	httpReq := g.RequestFromCtx(ctx)
	clientIP := httpReq.GetClientIp()
	if isNewUser {
		applyRegisterIpInfo(account, clientIP)
	} else {
		applyLoginIpInfo(account, clientIP)
	}
	accountdao.PublishAccountList(account.OpenId, account.Channel)

	userinfodao.GetUserInfoByUserId(account.ID)
	ext := userinfodao.GetUserExtByUserId(account.ID)
	if isNewUser && ext != nil && normalizeEmailKey(ext.Email) == "" {
		ext.SetEmail(email)
		userinfodao.PublishUserExt(ext)
	}

	tokenStr := xrtoken.AddAppToken(account.ID)
	if req.DeviceInfo != nil {
		userlogindevicedao.RefreshLoginDevice(account.ID, req.DeviceInfo)
		userinfodao.SaveRegisterInfo(account.ID, req.DeviceInfo)
	}
	userinfodao.SaveCancelCode(account.ID)
	if isNewUser {
		event.Pub(gameevent.RegisterEvent, gameevent.NewRegisterEventDataFromCtx(ctx, account.ID, time.Now()))
	}

	return &authdto.EmailLoginRes{
		Token:     fmt.Sprintf("%v.%s", account.ID, tokenStr),
		IsNewUser: isNewUser,
	}, nil
}

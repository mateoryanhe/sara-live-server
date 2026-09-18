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

// FirebaseLogin 验证 Firebase ID Token，并按 Firebase UID 登录或注册账号.
func FirebaseLogin(ctx context.Context, req *authdto.FirebaseLoginReq) (*authdto.FirebaseLoginRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := ensureSimulatorLoginAllowed(req.DeviceInfo); err != nil {
		return nil, err
	}
	uid, err := verifyFirebaseUID(ctx, req.IdToken)
	if err != nil {
		g.Log().Warningf(ctx, "firebase id token verification failed: %v", err)
		return nil, errercode.CreateCode(errercode.LoginFail)
	}

	lockKey := fmt.Sprintf("firebase_uid:%s", uid)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	account := findActiveAccountByBoundFirebaseUID(uid)
	isNewUser := false
	if account == nil {
		account = accountdao.FindActiveAccount(uid, FirebaseChannel)
	}
	if account == nil {
		if err := ensureFirebaseUIDAvailable(uid, 0); err != nil {
			return nil, err
		}
		account = accountdao.RegisterAccount(uid, FirebaseChannel)
		isNewUser = true
	}
	if account == nil || account.ID == 0 {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	if account.Ban && account.BanApplyTime != nil && account.BanApplyTime.After(time.Now()) {
		return nil, errercode.CreateCode(errercode.Ban)
	}

	httpReq := g.RequestFromCtx(ctx)
	if isNewUser {
		applyRegisterIpInfo(account, httpReq)
	} else {
		applyLoginIpInfo(account, httpReq)
	}
	accountdao.PublishAccountList(account.OpenId, account.Channel)

	userinfodao.GetUserInfoByUserId(account.ID)
	ext := userinfodao.GetUserExtByUserId(account.ID)
	if ext != nil && normalizeFirebaseUID(ext.FirebaseUID) == "" && account.Channel == FirebaseChannel {
		ext.SetFirebaseUID(uid)
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

	return &authdto.FirebaseLoginRes{
		Token:     fmt.Sprintf("%v.%s", account.ID, tokenStr),
		IsNewUser: isNewUser,
	}, nil
}

package auth

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmlock"
	"time"
	"xr-game-server/core/event"
	"xr-game-server/core/phoneutil"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/verification_code"
)

func PhoneRegister(ctx context.Context, req *authdto.PhoneRegisterReq) (res *authdto.PhoneRegisterRes, err error) {
	phoneAreaCode := phoneutil.NormalizeAreaCode(req.PhoneAreaCode)
	phone := req.Phone
	phoneKey := phoneutil.UniqueKey(phoneAreaCode, phone)
	lockKey := fmt.Sprintf("phone_register:%s:%d", phoneKey, PhoneChannel)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)
	// 验证验证码
	valid, err := verification_code.VerifyCode(phoneAreaCode, phone, req.Code)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errercode.CreateCode(errercode.VerifyCodeInvalid)
	}

	if active := accountdao.FindActivePhoneAccount(phoneAreaCode, phone, PhoneChannel); active != nil && active.Password != "" {
		return nil, errercode.CreateCode(errercode.AccountAlreadyExists)
	}
	account := accountdao.RegisterPhoneAccount(phoneAreaCode, phone, PhoneChannel)

	// 设置密码
	account.SetPassword(gmd5.MustEncryptString(req.Password))

	// 设置IP与国家
	httpReq := g.RequestFromCtx(ctx)
	applyRegisterIpInfo(account, httpReq.GetClientIp())

	// 生成token
	tokenStr := xrtoken.AddAppToken(account.ID)

	// 初始化用户信息
	data := userinfodao.GetUserInfoByUserId(account.ID)
	data.SetPhone(phone)
	userlogindevicedao.RefreshLoginDevice(account.ID, req.DeviceInfo)
	userinfodao.SaveRegisterInfo(account.ID, req.DeviceInfo)
	now := time.Now()
	event.Pub(gameevent.RegisterEvent, gameevent.NewRegisterEventDataFromCtx(ctx, account.ID, now))
	if req.InviteCode != "" {
		inviterId := userinfodao.GetUserIdByShareCode(req.InviteCode)
		if inviterId == 0 {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if inviterId == account.ID {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if data.InviterId == 0 {
			data.SetInviterId(inviterId)
		}
	}

	res = &authdto.PhoneRegisterRes{
		Token: fmt.Sprintf("%v.%s", account.ID, tokenStr),
	}
	return res, nil
}

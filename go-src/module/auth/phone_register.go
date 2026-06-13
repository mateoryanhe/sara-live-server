package auth

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmlock"
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/core/event"
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
	gmlock.Lock(req.Phone)
	defer gmlock.Unlock(req.Phone)
	// 验证验证码
	valid, err := verification_code.VerifyCode(req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errercode.CreateCode(errercode.VerifyCodeInvalid)
	}

	// 检查手机号是否已注册
	account := accountdao.GetAccountBy(req.Phone, PhoneChannel)
	if account.ID != 0 && account.Password != "" {
		return nil, errercode.CreateCode(errercode.AccountAlreadyExists)
	}

	// 设置密码
	account.SetPassword(gmd5.MustEncryptString(req.Password))

	// 设置IP
	httpReq := g.RequestFromCtx(ctx)
	if len(account.IP) == common.Zero {
		account.SetIp(httpReq.GetClientIp())
	}

	// 生成token
	tokenStr := xrtoken.AddAppToken(account.ID)

	// 初始化用户信息
	data := userinfodao.GetUserInfoByUserId(account.ID)
	data.SetPhone(req.Phone)
	userlogindevicedao.RefreshLoginDevice(account.ID, req.DeviceInfo)
	now := time.Now()
	event.Pub(gameevent.RegisterEvent, gameevent.NewRegisterEventData(account.ID, now))
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

package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/common"
	"xr-game-server/core/phoneutil"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
)

func PhoneLogin(ctx context.Context, req *authdto.PhoneLoginReq) (res *authdto.PhoneLoginRes, err error) {
	httpReq := g.RequestFromCtx(ctx)
	phoneAreaCode := phoneutil.NormalizeAreaCode(req.PhoneAreaCode)
	phone := req.Phone
	phoneKey := phoneutil.UniqueKey(phoneAreaCode, phone)
	if err = checkPhoneLoginLimit(phoneKey); err != nil {
		return nil, err
	}

	account := accountdao.FindActivePhoneAccount(phoneAreaCode, phone, PhoneChannel)
	if account == nil || account.Password == "" {
		if blockErr := markPhoneLoginFailure(phoneKey); blockErr != nil {
			return nil, blockErr
		}
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	if account.Cancel {
		return nil, errercode.CreateCode(errercode.AccountCanceled)
	}
	if account.Password != gmd5.MustEncryptString(req.Password) {
		if blockErr := markPhoneLoginFailure(phoneKey); blockErr != nil {
			return nil, blockErr
		}
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	clearPhoneLoginFailure(phoneKey)
	if account.Ban && account.BanApplyTime != nil && account.BanApplyTime.After(time.Now()) {
		return nil, errercode.CreateCode(errercode.Ban)
	}

	if len(account.IP) == common.Zero {
		account.SetIp(httpReq.GetClientIp())
	}
	tokenStr := xrtoken.AddAppToken(account.ID)
	userinfodao.GetUserInfoByUserId(account.ID)
	userlogindevicedao.RefreshLoginDevice(account.ID, req.DeviceInfo)

	return &authdto.PhoneLoginRes{
		Token: fmt.Sprintf("%v.%s", account.ID, tokenStr),
	}, nil
}

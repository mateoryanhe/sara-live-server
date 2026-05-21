package auth

import (
	"context"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
)

func PhoneLogin(ctx context.Context, req *authdto.PhoneLoginReq) (res *authdto.PhoneLoginRes, err error) {
	account := accountdao.GetAccountBy(req.Phone, PhoneChannel)
	if account.Password == "" {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	if account.Password != gmd5.MustEncryptString(req.Password) {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	if account.Ban && account.BanApplyTime != nil && account.BanApplyTime.After(time.Now()) {
		return nil, errercode.CreateCode(errercode.Ban)
	}

	httpReq := g.RequestFromCtx(ctx)
	if len(account.IP) == common.Zero {
		account.SetIp(httpReq.Host)
	}
	tokenStr := xrtoken.AddAppToken(account.ID)
	userinfodao.GetUserInfoByUserId(account.ID)

	return &authdto.PhoneLoginRes{
		Token:  tokenStr,
		AuthId: strconv.FormatUint(account.ID, 10),
	}, nil
}

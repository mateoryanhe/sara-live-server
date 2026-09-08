package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
)

// CoinMerchantLogin 币商用户名+密码登录(仅 CoinMerchantChannel,CMS 创建账号,不自动注册).
func CoinMerchantLogin(ctx context.Context, req *authdto.CoinMerchantLoginReq) (*authdto.CoinMerchantLoginRes, error) {
	if err := ensureSimulatorLoginAllowed(req.DeviceInfo); err != nil {
		return nil, err
	}
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if username == "" || len(password) < 6 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	account := accountdao.FindActiveAccount(username, CoinMerchantChannel)
	if account == nil || account.Password == "" {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	if account.Cancel {
		return nil, errercode.CreateCode(errercode.AccountCanceled)
	}
	if account.Password != gmd5.MustEncryptString(password) {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	if account.Ban && account.BanApplyTime != nil && account.BanApplyTime.After(time.Now()) {
		return nil, errercode.CreateCode(errercode.Ban)
	}

	httpReq := g.RequestFromCtx(ctx)
	applyLoginIpInfo(account, httpReq)
	accountdao.PublishAccountList(account.OpenId, account.Channel)

	tokenStr := xrtoken.AddAppToken(account.ID)
	userinfodao.GetUserInfoByUserId(account.ID)
	if req.DeviceInfo != nil {
		userlogindevicedao.RefreshLoginDevice(account.ID, req.DeviceInfo)
		userinfodao.SaveRegisterInfo(account.ID, req.DeviceInfo)
	}

	return &authdto.CoinMerchantLoginRes{
		Token: fmt.Sprintf("%v.%s", account.ID, tokenStr),
	}, nil
}

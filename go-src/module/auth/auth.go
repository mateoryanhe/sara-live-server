package auth

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
)

const (
	Year100 = 100 * 365 * 24 * time.Hour
)

func TestLogin(ctx context.Context, req *authdto.TestLoginReq) (res *authdto.TestLoginRes, err error) {
	//if !cfg.GetAuthCfg().LoginOff {
	//	return nil, errercode.CreateCode(errercode.TestEnvClose)
	//}
	account := accountdao.FindActiveAccount(req.OpenId, Test)
	isNewUser := false
	if account == nil {
		account = accountdao.RegisterAccount(req.OpenId, Test)
		isNewUser = true
	}
	data := account
	if data.Ban && data.BanApplyTime.After(time.Now()) {
		return nil, errercode.CreateCode(errercode.Ban)
	}
	httpReq := g.RequestFromCtx(ctx)
	clientIP := httpReq.GetClientIp()
	if isNewUser {
		applyRegisterIpInfo(data, clientIP)
	} else {
		applyLoginIpInfo(data, clientIP)
	}
	tokenStr := xrtoken.AddAppToken(data.ID)
	userinfodao.GetUserInfoByUserId(data.ID)
	res = &authdto.TestLoginRes{
		Token:  tokenStr,
		AuthId: strconv.FormatUint(data.ID, 10),
	}
	return res, nil
}

func CMSLogin(ctx context.Context, req *authdto.CMSLoginReq) (res *authdto.CMSLoginRes, err error) {

	data := cmsuserdao.GetCMSUser(req.UserName)
	if data == nil {
		return nil, errercode.CreateCode(errercode.CMSLoginFail)
	}
	if data.Pwd != req.Pwd {
		return nil, errercode.CreateCode(errercode.CMSLoginFail)
	}
	if data.Status == common.False {
		return nil, errercode.CreateCode(errercode.CMSLoginFail)
	}
	tokenStr := xrtoken.AddCmsToken(data.ID)
	perm := cmsuserdao.GetGetPermissionList(data.RoleId)
	return authdto.NewCMSLoginRes(data.ID, tokenStr, data.Admin, perm), nil
}

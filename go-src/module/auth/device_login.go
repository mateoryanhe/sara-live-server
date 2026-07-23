package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
)

func DeviceLogin(_ context.Context, req *authdto.DeviceLoginReq) (*authdto.DeviceLoginRes, error) {
	if req.DeviceInfo == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	deviceId := strings.TrimSpace(req.DeviceInfo.DeviceId)
	if deviceId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	userId := userlogindevicedao.FindUserIdByDeviceId(deviceId)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}

	account := accountdao.GetAccountById(userId)
	if account == nil {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	if account.Cancel {
		return nil, errercode.CreateCode(errercode.AccountCanceled)
	}
	if account.Ban && account.BanApplyTime != nil && account.BanApplyTime.After(time.Now()) {
		return nil, errercode.CreateCode(errercode.Ban)
	}

	tokenStr := xrtoken.AddAppToken(account.ID)
	userinfodao.GetUserInfoByUserId(account.ID)
	userlogindevicedao.RefreshLoginDevice(account.ID, req.DeviceInfo)
	userinfodao.SaveRegisterInfo(account.ID, req.DeviceInfo)

	return &authdto.DeviceLoginRes{
		Token: fmt.Sprintf("%v.%s", account.ID, tokenStr),
	}, nil
}

package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/event"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
	"xr-game-server/dto/authdto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

func DeviceLogin(ctx context.Context, req *authdto.DeviceLoginReq) (*authdto.DeviceLoginRes, error) {
	return loginByDevice(ctx, req.DeviceInfo, DeviceChannel, true)
}

func H5DeviceLogin(ctx context.Context, req *authdto.H5DeviceLoginReq) (*authdto.H5DeviceLoginRes, error) {
	return loginByDevice(ctx, req.DeviceInfo, H5DeviceChannel, false)
}

func loginByDevice(ctx context.Context, deviceInfo *entity.DeviceInfo, channel uint, checkSimulator bool) (*authdto.DeviceLoginRes, error) {
	if deviceInfo == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if checkSimulator {
		if err := ensureSimulatorLoginAllowed(deviceInfo); err != nil {
			return nil, err
		}
	}
	deviceId := strings.TrimSpace(deviceInfo.DeviceId)
	if deviceId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	lockKey := fmt.Sprintf("device_login:%s:%d", deviceId, channel)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	account := accountdao.FindActiveAccount(deviceId, channel)
	isNewUser := false
	if account == nil {
		account = accountdao.RegisterAccount(deviceId, channel)
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

	tokenStr := xrtoken.AddAppToken(account.ID)
	userinfodao.GetUserInfoByUserId(account.ID)
	userlogindevicedao.RefreshLoginDevice(account.ID, deviceInfo)
	userinfodao.SaveRegisterInfo(account.ID, deviceInfo)
	userinfodao.SaveCancelCode(account.ID)
	if isNewUser {
		event.Pub(gameevent.RegisterEvent, gameevent.NewRegisterEventDataFromCtx(ctx, account.ID, time.Now()))
	}

	return &authdto.DeviceLoginRes{
		Token:     fmt.Sprintf("%v.%s", account.ID, tokenStr),
		IsNewUser: isNewUser,
	}, nil
}

package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/constants/common"
	"xr-game-server/core/event"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dao/userlogindevicedao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

func DeviceLogin(ctx context.Context, req *authdto.DeviceLoginReq) (*authdto.DeviceLoginRes, error) {
	if req.DeviceInfo == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	deviceId := strings.TrimSpace(req.DeviceInfo.DeviceId)
	if deviceId == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	lockKey := deviceId + ":" + fmt.Sprint(DeviceChannel)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	account, isNewUser := accountdao.GetDeviceAccount(deviceId, DeviceChannel)
	if account == nil || account.ID == 0 {
		return nil, errercode.CreateCode(errercode.LoginFail)
	}
	if account.Ban && account.BanApplyTime != nil && account.BanApplyTime.After(time.Now()) {
		return nil, errercode.CreateCode(errercode.Ban)
	}

	httpReq := g.RequestFromCtx(ctx)
	if isNewUser && len(account.IP) == common.Zero {
		account.SetIp(httpReq.GetClientIp())
	}

	tokenStr := xrtoken.AddAppToken(account.ID)
	userinfodao.GetUserInfoByUserId(account.ID)
	userlogindevicedao.RefreshLoginDevice(account.ID, req.DeviceInfo)
	userinfodao.SaveRegisterInfo(account.ID, req.DeviceInfo)
	userinfodao.SaveCancelCode(account.ID)
	if isNewUser {
		event.Pub(gameevent.RegisterEvent, gameevent.NewRegisterEventData(account.ID, time.Now()))
	}

	return &authdto.DeviceLoginRes{
		Token:     fmt.Sprintf("%v.%s", account.ID, tokenStr),
		IsNewUser: isNewUser,
	}, nil
}

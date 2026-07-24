package userinfo

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/phoneutil"
	"xr-game-server/core/push"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/dto/userinfodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/auth"
	"xr-game-server/module/verification_code"
)

func CancelUser(ctx context.Context, req *accountdto.CancelReq) (bool, error) {
	AddIdToCache(req.AccountId)
	db := accountdao.GetAccountById(req.AccountId)
	accountCache := accountdao.GetAccountBy(db.OpenId, db.Channel, db.PhoneAreaCode)
	accountCache.SetCancel(true)
	invalidateAppToken(req.AccountId)
	push.Kick(req.AccountId)
	return true, nil
}

func UnCancelUser(ctx context.Context, req *accountdto.UnCancelReq) (bool, error) {
	AddIdToCache(req.AccountId)
	db := accountdao.GetAccountById(req.AccountId)
	accountCache := accountdao.GetAccountBy(db.OpenId, db.Channel, db.PhoneAreaCode)
	accountCache.SetCancel(false)
	return true, nil
}

// CancelAccount App端销户(需登录,无需请求体)
func CancelAccount(ctx context.Context, _ *userinfodto.CancelAccountReq) (*userinfodto.CancelAccountRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	account := accountdao.GetAccountById(userId)
	if account == nil || account.ID == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if account.Cancel {
		return nil, errercode.CreateCode(errercode.AccountCanceled)
	}
	if err := doAppCancelAccount(account); err != nil {
		return nil, err
	}
	return &userinfodto.CancelAccountRes{Success: true}, nil
}

// CancelAccountByPhone 手机号+验证码注销(公开接口,无需登录)
func CancelAccountByPhone(ctx context.Context, req *userinfodto.CancelAccountByPhoneReq) (*userinfodto.CancelAccountByPhoneRes, error) {
	ip := g.RequestFromCtx(ctx).GetClientIp()
	if err := checkCancelByCodeIPLimit(ip); err != nil {
		return nil, err
	}
	valid, err := verification_code.VerifyCode(req.PhoneAreaCode, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, markCancelByCodeFailure(ip, errercode.VerifyCodeInvalid)
	}
	phoneAreaCode := phoneutil.NormalizeAreaCode(req.PhoneAreaCode)
	account := accountdao.FindAccountBy(req.Phone, auth.PhoneChannel, phoneAreaCode)
	if account == nil || account.Password == "" {
		return nil, markCancelByCodeFailure(ip, errercode.InvalidParam)
	}
	if account.Cancel {
		return nil, markCancelByCodeFailure(ip, errercode.AccountCanceled)
	}
	if err := doAppCancelAccount(account); err != nil {
		return nil, markCancelByCodeFailure(ip, errercode.InvalidParam)
	}
	clearCancelByCodeFailure(ip)
	return &userinfodto.CancelAccountByPhoneRes{Success: true}, nil
}

func doAppCancelAccount(account *entity.Account) error {
	if account == nil || account.ID == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	AddIdToCache(account.ID)
	originalOpenId := account.OpenId
	channel := account.Channel
	accountCache := accountdao.GetAccountBy(account.OpenId, account.Channel, account.PhoneAreaCode)
	accountCache.SetCancel(true)
	if channel == auth.DeviceChannel {
		accountdao.FlushDeviceAccountCache(originalOpenId, channel)
	}
	invalidateAppToken(account.ID)
	push.Kick(account.ID)
	return nil
}

// invalidateAppToken 覆盖内存与数据库中的 Token,使旧 Token 立即失效
func invalidateAppToken(userId uint64) {
	if userId == 0 {
		return
	}
	xrtoken.InitAppToken(userId, guid.S(), time.Now().Add(30*time.Minute))
	expireTime := time.Now().Add(-24 * 100 * time.Hour)
	entity.NewAppToken(userId, fmt.Sprintf("%v", userId), expireTime)
}

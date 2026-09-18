package auth

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gmlock"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/authdto"
	"xr-game-server/errercode"
)

// BindFirebase 将验证通过的 Firebase UID 绑定到当前登录账号。
func BindFirebase(ctx context.Context, req *authdto.BindFirebaseReq) (*authdto.BindFirebaseRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	firebaseUID, err := verifyFirebaseUID(ctx, req.IdToken)
	if err != nil {
		g.Log().Warningf(ctx, "firebase bind token verification failed userId=%d: %v", userId, err)
		return nil, errercode.CreateCode(errercode.LoginFail)
	}

	lockKey := fmt.Sprintf("firebase_uid:%s", firebaseUID)
	gmlock.Lock(lockKey)
	defer gmlock.Unlock(lockKey)

	ext := userinfodao.GetUserExtByUserId(userId)
	if ext == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	current := normalizeFirebaseUID(ext.FirebaseUID)
	if current != "" {
		if current == firebaseUID {
			return &authdto.BindFirebaseRes{Success: true}, nil
		}
		return nil, errercode.CreateCode(errercode.FirebaseAlreadyBound)
	}
	if err := ensureFirebaseUIDAvailable(firebaseUID, userId); err != nil {
		return nil, err
	}

	ext.SetFirebaseUID(firebaseUID)
	userinfodao.PublishUserExt(ext)
	return &authdto.BindFirebaseRes{Success: true}, nil
}

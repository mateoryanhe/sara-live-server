package call

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/userchanneltokendao"
	"xr-game-server/dto/calldto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/agora"
)

const channelTokenReuseMinRemain = 2 * time.Hour

// resolveChannelToken 优先复用 user_channel_tokens 中仍有效的 Token
func resolveChannelToken(userId uint64, channelName string) (token string, tokenExpireAt time.Time, err error) {
	now := time.Now()
	if existing := userchanneltokendao.GetByUserChannel(userId, channelName); existing != nil && existing.Token != "" {
		if existing.ExpireAt.After(now.Add(channelTokenReuseMinRemain)) {
			return existing.Token, existing.ExpireAt, nil
		}
	}

	token, expireSeconds, err := agora.BuildCallToken(userId, channelName)
	if err != nil {
		return "", time.Time{}, err
	}
	tokenExpireAt = now.Add(time.Duration(expireSeconds) * time.Second)
	row := entity.NewUserChannelToken(userId, channelName, token, tokenExpireAt)
	userchanneltokendao.AddToCache(row)
	return token, tokenExpireAt, nil
}

// GetCallToken App端上报 channelName 获取通话频道 Token
func GetCallToken(ctx context.Context, req *calldto.CallTokenReq) (*calldto.CallTokenRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	token, tokenExpireAt, err := resolveChannelToken(userId, req.ChannelName)
	if err != nil {
		return nil, err
	}

	appId := ""
	if agoraCfg, err := agora.GetAppId(ctx, nil); err == nil && agoraCfg != nil {
		appId = agoraCfg.AppId
	}

	return &calldto.CallTokenRes{
		Token:       token,
		AppId:       appId,
		ChannelName: req.ChannelName,
		UserAccount: strconv.FormatUint(userId, 10),
		ExpireAt:    tokenExpireAt.Unix(),
	}, nil
}

package agora

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userchanneltokendao"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

const channelTokenReuseMinRemain = 2 * time.Hour

type channelTokenBuilder func() (token string, expireSeconds int64, err error)

func resolveChannelToken(userId uint64, channelName string, build channelTokenBuilder) (token string, expireAt int64, err error) {
	now := time.Now()
	if existing := userchanneltokendao.GetByUserChannel(userId, channelName); existing != nil && existing.Token != "" {
		if existing.ExpireAt.After(now.Add(channelTokenReuseMinRemain)) {
			scheduleChannelTokenRefresh(userId, channelName, build)
			return existing.Token, existing.ExpireAt.Unix(), nil
		}
	}
	return refreshChannelToken(userId, channelName, build)
}

func scheduleChannelTokenRefresh(userId uint64, channelName string, build channelTokenBuilder) {
	xrpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		_, _, _ = refreshChannelToken(userId, channelName, build)
	})
}

func refreshChannelToken(userId uint64, channelName string, build channelTokenBuilder) (token string, expireAt int64, err error) {
	token, expireSeconds, err := build()
	if err != nil {
		return "", 0, err
	}
	tokenExpireAt := time.Now().Add(time.Duration(expireSeconds) * time.Second)
	row := entity.NewUserChannelToken(userId, channelName, token, tokenExpireAt)
	userchanneltokendao.AddToCache(row)
	return token, tokenExpireAt.Unix(), nil
}

// ResolveLiveRoomToken 优先复用缓存中的直播间声网 Token
func ResolveLiveRoomToken(userId, roomId uint64) (token string, expireAt int64, err error) {
	if liveroomdao.GetRoomById(roomId) == nil {
		return "", 0, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	channelName := buildChannelName(roomId)
	return resolveChannelToken(userId, channelName, func() (string, int64, error) {
		return BuildLiveRoomToken(userId, roomId)
	})
}

// ResolveCallToken 优先复用缓存中的通话频道声网 Token
func ResolveCallToken(userId uint64, channelName string) (token string, expireAt int64, err error) {
	return resolveChannelToken(userId, channelName, func() (string, int64, error) {
		return BuildCallToken(userId, channelName)
	})
}

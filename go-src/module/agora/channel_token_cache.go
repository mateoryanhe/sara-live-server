package agora

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/userchanneltokendao"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

type channelTokenBuilder func() (token string, expireSeconds int64, err error)

func getChannelTokenReuseMinRemain() time.Duration {
	return time.Duration(getAgoraCfgCache().TokenRefreshSeconds) * time.Second
}

// resolveChannelToken 查询时百分百刷新 Token:
// 剩余有效期大于提前刷新阈值则先返回缓存并后台静默刷新,否则同步生成
func resolveChannelToken(userId uint64, channelName string, build channelTokenBuilder) (token string, expireAt int64, err error) {
	minRemain := getChannelTokenReuseMinRemain()
	now := time.Now()
	if existing := userchanneltokendao.GetByUserChannel(userId, channelName); existing != nil && existing.Token != "" {
		if existing.ExpireAt.After(now.Add(minRemain)) {
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

// ResolveChannelToken 按频道名查询声网 Token(频道名与角色由业务方决定)
func ResolveChannelToken(userId uint64, channelName string, role uint8) (token string, expireAt int64, err error) {
	if channelName == "" {
		return "", 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return resolveChannelToken(userId, channelName, func() (string, int64, error) {
		return BuildChannelToken(userId, channelName, role)
	})
}

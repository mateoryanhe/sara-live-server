package preload

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/userinfodao"
)

func preloadRecentLoginUsers() {
	ctx := gctx.New()
	userIds := userinfodao.PreloadRecentLoginUserInfos(RecentLoginPreloadLimit)
	if len(userIds) == 0 {
		g.Log().Info(ctx, "preload recent login users skipped: no users")
		return
	}
	userinfodao.PreloadUserExtToCache(userIds)
	userinfodao.PreloadUserCumulativeStatToCache(userIds)
	preloadAppTokenCache(RecentLoginPreloadLimit)
	g.Log().Infof(ctx, "preload recent login users done, count=%d", len(userIds))
}

func preloadAppTokenCache(limit int) {
	tokens := accountdao.ListAppTokensByRecentLogin(limit)
	items := make([]xrtoken.AppTokenCacheItem, 0, len(tokens))
	for _, token := range tokens {
		if token == nil || token.ID == 0 || token.Token == "" {
			continue
		}
		items = append(items, xrtoken.AppTokenCacheItem{
			UserId:   token.ID,
			Token:    token.Token,
			ExpireAt: token.ExpireAt,
		})
	}
	xrtoken.PreloadAppTokenCache(items)
}

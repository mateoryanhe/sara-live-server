package userinfodao

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	userentity "xr-game-server/entity/user"
)

var emailUserIdCacheMgr *cache.RowCache[uint64]

func initEmailUserIdCache() {
	emailUserIdCacheMgr = cache.NewRowCache[uint64]()
}

func normalizeEmailCacheKey(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// GetActiveUserIdByEmail 按绑定邮箱查未注销用户ID(带缓存);不存在返回 0.
func GetActiveUserIdByEmail(email string) uint64 {
	email = normalizeEmailCacheKey(email)
	if email == "" {
		return 0
	}
	if emailUserIdCacheMgr == nil {
		return loadActiveUserIdByEmailFromDB(email)
	}
	return emailUserIdCacheMgr.MustGetRow(gctx.New(), email, func(ctx context.Context) (uint64, error) {
		return loadActiveUserIdByEmailFromDB(email), nil
	})
}

func loadActiveUserIdByEmailFromDB(email string) uint64 {
	email = normalizeEmailCacheKey(email)
	if email == "" {
		return 0
	}
	var userId uint64
	// user_exts.email + accounts.cancel=0
	_ = g.Model(string(userentity.TbUserExt)+" e").Ctx(gctx.New()).
		LeftJoin(string(userentity.TbAccount)+" a", "a.id=e.id").
		Where("e."+string(userentity.UserExtEmail)+" = ?", email).
		Where("IFNULL(a.cancel, 0) = 0").
		Fields("e."+string(db.IdName)).
		Order("e."+string(db.IdName)+" desc").
		Limit(1).
		Scan(&userId)
	return userId
}

// PublishEmailUserIdCache 绑定/更新邮箱后写入 email→userId 缓存
func PublishEmailUserIdCache(email string, userId uint64) {
	email = normalizeEmailCacheKey(email)
	if email == "" || userId == 0 || emailUserIdCacheMgr == nil {
		return
	}
	emailUserIdCacheMgr.PublishRow(gctx.New(), email, userId)
}

// InvalidateEmailUserIdCache 注销或换绑时失效邮箱索引缓存
func InvalidateEmailUserIdCache(email string) {
	email = normalizeEmailCacheKey(email)
	if email == "" || emailUserIdCacheMgr == nil {
		return
	}
	emailUserIdCacheMgr.RemoveRow(gctx.New(), email)
}

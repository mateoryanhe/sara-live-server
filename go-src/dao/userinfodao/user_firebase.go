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

var firebaseUserIdCacheMgr *cache.RowCache[uint64]

func initFirebaseUserIdCache() {
	firebaseUserIdCacheMgr = cache.NewRowCache[uint64]()
}

func normalizeFirebaseUIDCacheKey(firebaseUID string) string {
	return strings.TrimSpace(firebaseUID)
}

// GetActiveUserIdByFirebaseUID 按绑定 Firebase UID 查未注销用户 ID（带缓存）；不存在返回 0。
func GetActiveUserIdByFirebaseUID(firebaseUID string) uint64 {
	firebaseUID = normalizeFirebaseUIDCacheKey(firebaseUID)
	if firebaseUID == "" {
		return 0
	}
	if firebaseUserIdCacheMgr == nil {
		return loadActiveUserIdByFirebaseUIDFromDB(firebaseUID)
	}
	return firebaseUserIdCacheMgr.MustGetRow(gctx.New(), firebaseUID, func(ctx context.Context) (uint64, error) {
		return loadActiveUserIdByFirebaseUIDFromDB(firebaseUID), nil
	})
}

func loadActiveUserIdByFirebaseUIDFromDB(firebaseUID string) uint64 {
	firebaseUID = normalizeFirebaseUIDCacheKey(firebaseUID)
	if firebaseUID == "" {
		return 0
	}
	var userId uint64
	_ = g.Model(string(userentity.TbUserExt)+" e").Ctx(gctx.New()).
		LeftJoin(string(userentity.TbAccount)+" a", "a.id=e.id").
		Where("e."+string(userentity.UserExtFirebaseUID)+" = ?", firebaseUID).
		Where("IFNULL(a.cancel, 0) = 0").
		Fields("e." + string(db.IdName)).
		Order("e." + string(db.IdName) + " desc").
		Limit(1).
		Scan(&userId)
	return userId
}

// PublishFirebaseUserIdCache 绑定后写入 Firebase UID→userId 缓存。
func PublishFirebaseUserIdCache(firebaseUID string, userId uint64) {
	firebaseUID = normalizeFirebaseUIDCacheKey(firebaseUID)
	if firebaseUID == "" || userId == 0 || firebaseUserIdCacheMgr == nil {
		return
	}
	firebaseUserIdCacheMgr.PublishRow(gctx.New(), firebaseUID, userId)
}

// InvalidateFirebaseUserIdCache 注销时失效 Firebase UID 索引缓存。
func InvalidateFirebaseUserIdCache(firebaseUID string) {
	firebaseUID = normalizeFirebaseUIDCacheKey(firebaseUID)
	if firebaseUID == "" || firebaseUserIdCacheMgr == nil {
		return
	}
	firebaseUserIdCacheMgr.RemoveRow(gctx.New(), firebaseUID)
}

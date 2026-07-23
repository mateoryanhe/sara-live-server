package userinfodao

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var (
	userExtCacheMgr    *cache.CacheMgr
	cancelCodeCacheMgr *cache.CacheMgr
)

const cancelCodeMissCacheTime = 10 * time.Minute

func initUserExtDao() {
	userExtCacheMgr = cache.NewCacheMgr()
	cancelCodeCacheMgr = cache.NewCacheMgr()
}

func cancelCodeMissCacheKey(cancelCode string) string {
	return "miss:" + cancelCode
}

func syncCancelCodeCache(cancelCode string, userId uint64) {
	cancelCodeCacheMgr.FlushCache(cancelCode, userId)
	ctx := gctx.New()
	_, _ = cancelCodeCacheMgr.Cache.Remove(ctx, cancelCodeMissCacheKey(cancelCode))
}

// GetUserExtByUserId 获取用户扩展信息,不存在则新建内存对象(异步入库,默认允许上排行榜)
func GetUserExtByUserId(userId uint64) *entity.UserExt {
	cacheData := userExtCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		var data *entity.UserExt
		err = g.Model(string(entity.TbUserExt)).Unscoped().Where(g.Map{
			string(db.IdName): userId,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewUserExt(userId), nil
	})
	return cacheData.(*entity.UserExt)
}

// SaveRegisterInfo 保存注册时的包名与版本号(可为空)
func SaveRegisterInfo(userId uint64, info *entity.DeviceInfo) {
	if userId == 0 || info == nil {
		return
	}
	ext := GetUserExtByUserId(userId)
	ext.SetPackageName(strings.TrimSpace(info.PackageName))
	ext.SetAppVersion(strings.TrimSpace(info.AppVersion))
}

// SaveCancelCode 生成并保存注销码,仅首次写入
func SaveCancelCode(userId uint64) {
	if userId == 0 {
		return
	}
	ext := GetUserExtByUserId(userId)
	if ext.CancelCode != "" {
		return
	}
	cancelCode := guid.S()
	ext.SetCancelCode(cancelCode)
	syncCancelCodeCache(cancelCode, userId)
}

// FindUserIdByCancelCode 根据注销码只读查询用户ID,不存在时负缓存 10 分钟
func FindUserIdByCancelCode(cancelCode string) uint64 {
	cancelCode = strings.TrimSpace(cancelCode)
	if cancelCode == "" {
		return 0
	}

	ctx := gctx.New()
	if ok, _ := cancelCodeCacheMgr.Cache.Contains(ctx, cancelCodeMissCacheKey(cancelCode)); ok {
		return 0
	}

	if cached := cancelCodeCacheMgr.GetFromCache(cancelCode); cached != nil {
		if userId, ok := cached.(uint64); ok && userId > 0 {
			return userId
		}
	}

	var row struct {
		ID uint64
	}
	err := g.Model(string(entity.TbUserExt)).
		Where(string(entity.UserExtCancelCode), cancelCode).
		Limit(1).
		Scan(&row)
	if err != nil || row.ID == 0 {
		_ = cancelCodeCacheMgr.Cache.Set(ctx, cancelCodeMissCacheKey(cancelCode), 1, cancelCodeMissCacheTime)
		return 0
	}

	cancelCodeCacheMgr.FlushCache(cancelCode, row.ID)
	return row.ID
}

// GetFollowCount 获取用户当前关注数(走缓存)
func GetFollowCount(userId uint64) int {
	if userId == 0 {
		return 0
	}
	return int(GetUserExtByUserId(userId).FollowCount)
}

// GetFollowerCount 获取用户当前粉丝数(走缓存)
func GetFollowerCount(userId uint64) int {
	if userId == 0 {
		return 0
	}
	return int(GetUserExtByUserId(userId).FollowerCount)
}

// IncFollowCount 关注成功后增加双方计数
func IncFollowCount(userId, anchorId uint64) {
	if userId == 0 || anchorId == 0 {
		return
	}
	GetUserExtByUserId(userId).AddFollowCount(1)
	GetUserExtByUserId(anchorId).AddFollowerCount(1)
}

// DecFollowCount 取消关注后减少双方计数
func DecFollowCount(userId, anchorId uint64) {
	if userId == 0 || anchorId == 0 {
		return
	}
	GetUserExtByUserId(userId).SubFollowCount(1)
	GetUserExtByUserId(anchorId).SubFollowerCount(1)
}

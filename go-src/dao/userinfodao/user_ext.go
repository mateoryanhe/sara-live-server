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

const cancelCodeValidDuration = 3 * 24 * time.Hour

var userExtCacheMgr *cache.CacheMgr

func initUserExtDao() {
	userExtCacheMgr = cache.NewCacheMgr()
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

// GetUserExtFromMemory 仅从内存缓存读取用户扩展信息,未命中返回 nil
func GetUserExtFromMemory(userId uint64) *entity.UserExt {
	if userId == 0 || userExtCacheMgr == nil {
		return nil
	}
	v := userExtCacheMgr.GetFromCache(userId)
	if v == nil {
		return nil
	}
	ext, ok := v.(*entity.UserExt)
	if !ok || ext == nil {
		return nil
	}
	return ext
}

// IsRechargeWhitelist 用户是否在充值白名单(创建 App 订单后直接到账)
func IsRechargeWhitelist(userId uint64) bool {
	if userId == 0 {
		return false
	}
	return GetUserExtByUserId(userId).RechargeWhitelist
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

// IsCancelCodeValid 注销码是否存在且在有效期内
func IsCancelCodeValid(ext *entity.UserExt) bool {
	if ext == nil || strings.TrimSpace(ext.CancelCode) == "" {
		return false
	}
	if ext.CancelCodeExpireAt == nil || ext.CancelCodeExpireAt.IsZero() {
		return false
	}
	return !time.Now().After(*ext.CancelCodeExpireAt)
}

func refreshCancelCode(ext *entity.UserExt) {
	expireAt := time.Now().Add(cancelCodeValidDuration)
	ext.SetCancelCode(guid.S())
	ext.SetCancelCodeExpireAt(&expireAt)
}

// EnsureCancelCode 注销码为空、过期时间缺失或已过期时重新生成(默认有效期 3 天)
func EnsureCancelCode(userId uint64) *entity.UserExt {
	ext := GetUserExtByUserId(userId)
	if !IsCancelCodeValid(ext) {
		refreshCancelCode(ext)
	}
	return ext
}

// SaveCancelCode 登录时确保注销码有效
func SaveCancelCode(userId uint64) {
	if userId == 0 {
		return
	}
	EnsureCancelCode(userId)
}

// FindUserIdByCancelCode 根据注销码查询用户ID
func FindUserIdByCancelCode(cancelCode string) uint64 {
	cancelCode = strings.TrimSpace(cancelCode)
	if cancelCode == "" {
		return 0
	}
	var row struct {
		ID uint64
	}
	err := g.Model(string(entity.TbUserExt)).
		Where(string(entity.UserExtCancelCode), cancelCode).
		Limit(1).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return 0
	}
	return row.ID
}

// PreloadUserExtToCache 批量预热 user_exts 缓存
func PreloadUserExtToCache(userIds []uint64) {
	if len(userIds) == 0 || userExtCacheMgr == nil {
		return
	}
	ctx := gctx.New()
	rows := make([]*entity.UserExt, 0, len(userIds))
	err := g.Model(string(entity.TbUserExt)).Ctx(ctx).Unscoped().
		WhereIn(string(db.IdName), userIds).
		Scan(&rows)
	if err != nil {
		g.Log().Errorf(ctx, "preload user exts failed: %v", err)
		return
	}
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		userExtCacheMgr.FlushCache(row.ID, row)
	}
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

package userinfodao

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

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

// SaveRegisterInfo 保存注册时的包名与版本号(可为空)
func SaveRegisterInfo(userId uint64, info *entity.DeviceInfo) {
	if userId == 0 || info == nil {
		return
	}
	ext := GetUserExtByUserId(userId)
	ext.SetPackageName(strings.TrimSpace(info.PackageName))
	ext.SetAppVersion(strings.TrimSpace(info.AppVersion))
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

package userinfodao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var userInfoCacheMgr *cache.CacheMgr
var shareCodeUserIdCacheMgr *cache.CacheMgr

func InitUserInfoDao() {
	userInfoCacheMgr = cache.NewCacheMgr()
	shareCodeUserIdCacheMgr = cache.NewCacheMgr()
}

// GetUserIdByShareCode 根据分享码获取玩家ID,不存在则返回 0
func GetUserIdByShareCode(shareCode string) uint64 {
	if shareCode == "" {
		return 0
	}
	cacheData := shareCodeUserIdCacheMgr.GetData(shareCode, func(ctx context.Context) (value interface{}, err error) {
		var userId uint64
		err = g.Model(string(entity.TbUserInfo)).Unscoped().Where(g.Map{
			string(entity.UserInfoShareCode): shareCode,
		}).Fields(string(db.IdName)).Scan(&userId)
		return userId, err
	})
	if cacheData == nil {
		return 0
	}
	if userId, ok := cacheData.(uint64); ok {
		return userId
	}
	return 0
}

// GetUserInfoByUserId 根据用户ID获取用户基础信息,命中不了缓存从数据库拉取,数据库不存在则新建
func GetUserInfoByUserId(userId uint64) *entity.UserInfo {
	cacheData := userInfoCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		var data *entity.UserInfo
		err = g.Model(string(entity.TbUserInfo)).Unscoped().Where(g.Map{
			string(db.IdName): userId,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		//数据库没有,新建一条
		newData := entity.NewUserInfo(userId)
		return newData, nil
	})
	return cacheData.(*entity.UserInfo)
}

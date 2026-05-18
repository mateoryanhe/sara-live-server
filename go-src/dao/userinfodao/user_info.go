package userinfodao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var userInfoCacheMgr *cache.CacheMgr

func InitUserInfoDao() {
	userInfoCacheMgr = cache.NewCacheMgr()
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

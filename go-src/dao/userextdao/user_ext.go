package userextdao

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var userExtCacheMgr *cache.CacheMgr

func InitUserExtDao() {
	userExtCacheMgr = cache.NewCacheMgr()
}

// GetByUserId 获取用户扩展信息,不存在则新建内存对象(异步入库,默认允许上排行榜)
func GetByUserId(userId uint64) *entity.UserExt {
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
	ext := GetByUserId(userId)
	ext.SetPackageName(strings.TrimSpace(info.PackageName))
	ext.SetAppVersion(strings.TrimSpace(info.AppVersion))
}

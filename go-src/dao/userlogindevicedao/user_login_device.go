package userlogindevicedao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var loginDeviceCacheMgr *cache.CacheMgr

func InitUserLoginDeviceDao() {
	loginDeviceCacheMgr = cache.NewCacheMgr()
}

// GetByUserId 获取用户登录设备信息,不存在则新建内存对象(异步入库)
func GetByUserId(userId uint64) *entity.UserLoginDevice {
	cacheData := loginDeviceCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		var data *entity.UserLoginDevice
		err = g.Model(string(entity.TbUserLoginDevice)).Unscoped().Where(g.Map{
			string(db.IdName): userId,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewUserLoginDevice(userId), nil
	})
	return cacheData.(*entity.UserLoginDevice)
}

// RefreshLoginDevice 刷新用户登录设备信息(缓冲写入,不直接入库)
func RefreshLoginDevice(userId uint64, info *entity.DeviceInfo) {
	if info == nil {
		return
	}
	device := entity.NewUserLoginDevice(userId)
	device.Refresh(info)
}

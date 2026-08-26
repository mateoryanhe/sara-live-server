package userlogindevicedao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/user"
)

var loginDeviceCacheMgr *cache.RowCache[*entity.UserLoginDevice]

func InitUserLoginDeviceDao() {
	loginDeviceCacheMgr = cache.NewRowCache[*entity.UserLoginDevice]()
}

// GetByUserId 获取用户登录设备信息,不存在则新建内存对象(异步入库)
func GetByUserId(userId uint64) *entity.UserLoginDevice {
	return loginDeviceCacheMgr.MustGetRow(gctx.New(), userId, func(ctx context.Context) (*entity.UserLoginDevice, error) {
		var data *entity.UserLoginDevice
		_ = g.Model(string(entity.TbUserLoginDevice)).Unscoped().Where(g.Map{
			string(db.IdName): userId,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		return entity.NewUserLoginDevice(userId), nil
	})
}

// PublishLoginDevice 原地修改登录设备后刷新缓存.
func PublishLoginDevice(data *entity.UserLoginDevice) {
	if data == nil || data.ID == 0 || loginDeviceCacheMgr == nil {
		return
	}
	loginDeviceCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

// RefreshLoginDevice 刷新用户登录设备信息(缓冲写入,不直接入库)
func RefreshLoginDevice(userId uint64, info *entity.DeviceInfo) {
	if info == nil {
		return
	}
	device := GetByUserId(userId)
	device.Refresh(info)
	PublishLoginDevice(device)
}

// FindUserIdByDeviceId 根据设备码查找最近在该设备登录的用户ID
func FindUserIdByDeviceId(deviceId string) uint64 {
	deviceId = strings.TrimSpace(deviceId)
	if deviceId == "" {
		return 0
	}
	var row struct {
		ID uint64
	}
	err := g.Model(string(entity.TbUserLoginDevice)).
		Where(string(entity.UserLoginDeviceId), deviceId).
		OrderDesc(string(db.UpdatedAtName)).
		Limit(1).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return 0
	}
	return row.ID
}

// GetUserLoginDeviceFromDB 直查 user_login_devices(不存在返回 nil,不新建)
func GetUserLoginDeviceFromDB(userId uint64) *entity.UserLoginDevice {
	if userId == 0 {
		return nil
	}
	var row entity.UserLoginDevice
	err := g.Model(string(entity.TbUserLoginDevice)).Unscoped().Where(g.Map{
		string(db.IdName): userId,
	}).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

package entity

import (
	"strings"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbUserLoginDevice db.TbName = "user_login_devices"
)

const (
	UserLoginDeviceType   db.TbCol = "device_type"
	UserLoginDeviceModel  db.TbCol = "device_model"
	UserLoginDeviceOsVer  db.TbCol = "os_version"
	UserLoginDeviceAppVer db.TbCol = "app_version"
	UserLoginDeviceId     db.TbCol = "device_id"
)

// UserLoginDevice 用户最近登录设备信息(与用户一一对应,主键ID即用户ID)
type UserLoginDevice struct {
	migrate.OneModel
	DeviceType  string `gorm:"size:16;default:'';comment:设备类型(ios/android)" json:"deviceType"`
	DeviceModel string `gorm:"size:128;default:'';comment:设备型号" json:"deviceModel"`
	OsVersion   string `gorm:"size:64;default:'';comment:系统版本号" json:"osVersion"`
	AppVersion  string `gorm:"size:64;default:'';comment:App/APK版本号" json:"appVersion"`
	DeviceId    string `gorm:"size:128;default:'';comment:设备唯一标识" json:"deviceId"`
}

func NewUserLoginDevice(userId uint64) *UserLoginDevice {
	ret := &UserLoginDevice{}
	ret.ID = userId
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

func (receiver *UserLoginDevice) Refresh(info *DeviceInfo) {
	if info == nil {
		return
	}
	receiver.SetDeviceType(strings.ToLower(strings.TrimSpace(info.DeviceType)))
	receiver.SetDeviceModel(info.DeviceModel)
	receiver.SetOsVersion(info.OsVersion)
	receiver.SetAppVersion(info.AppVersion)
	receiver.SetDeviceId(info.DeviceId)
}

func (receiver *UserLoginDevice) SetDeviceType(deviceType string) {
	receiver.DeviceType = deviceType
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserLoginDevice, UserLoginDeviceType, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: deviceType,
	})
}

func (receiver *UserLoginDevice) SetDeviceModel(deviceModel string) {
	receiver.DeviceModel = deviceModel
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserLoginDevice, UserLoginDeviceModel, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: deviceModel,
	})
}

func (receiver *UserLoginDevice) SetOsVersion(osVersion string) {
	receiver.OsVersion = osVersion
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserLoginDevice, UserLoginDeviceOsVer, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: osVersion,
	})
}

func (receiver *UserLoginDevice) SetAppVersion(appVersion string) {
	receiver.AppVersion = appVersion
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserLoginDevice, UserLoginDeviceAppVer, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: appVersion,
	})
}

func (receiver *UserLoginDevice) SetDeviceId(deviceId string) {
	receiver.DeviceId = deviceId
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserLoginDevice, UserLoginDeviceId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: deviceId,
	})
}

func (receiver *UserLoginDevice) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToQuickChan(TbUserLoginDevice, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserLoginDevice) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToQuickChan(TbUserLoginDevice, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initUserLoginDevice() {
	syndb.RegQuickWithLarge(TbUserLoginDevice, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbUserLoginDevice, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbUserLoginDevice, UserLoginDeviceType)
	syndb.RegQuickWithLarge(TbUserLoginDevice, UserLoginDeviceModel)
	syndb.RegQuickWithLarge(TbUserLoginDevice, UserLoginDeviceOsVer)
	syndb.RegQuickWithLarge(TbUserLoginDevice, UserLoginDeviceAppVer)
	syndb.RegQuickWithLarge(TbUserLoginDevice, UserLoginDeviceId)

	migrate.AutoMigrate(&UserLoginDevice{})
}

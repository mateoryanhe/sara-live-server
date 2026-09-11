package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbSimulatorDeviceWhitelist db.TbName = "simulator_device_whitelists"
)

const (
	SimulatorDeviceWhitelistDeviceId db.TbCol = "device_id"
)

// SimulatorDeviceWhitelist 模拟器登录设备码白名单(CMS 管理,命中则放行拦截)
type SimulatorDeviceWhitelist struct {
	migrate.OneModel
	DeviceId string `gorm:"size:128;uniqueIndex;default:'';comment:设备码(deviceId)" json:"deviceId"`
}

func initSimulatorDeviceWhitelist() {
	migrate.AutoMigrate(&SimulatorDeviceWhitelist{})
}

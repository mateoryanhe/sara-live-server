package simulatordevicewhitelistdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type SimulatorDeviceWhitelistListReq struct {
	g.Meta `path:"/simulatorDeviceWhitelistList" method:"post" summary:"获取模拟器设备白名单列表" tags:"模拟器设备白名单"`
	httpserver.CMSQueryReq
	Key string `json:"key" dc:"设备码模糊搜索"`
}

type SimulatorDeviceWhitelistItem struct {
	ID        string `json:"id"`
	DeviceId  string `json:"deviceId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateSimulatorDeviceWhitelistReq struct {
	g.Meta   `path:"/createSimulatorDeviceWhitelist" method:"post" summary:"新增模拟器设备白名单" tags:"模拟器设备白名单"`
	DeviceId string `json:"deviceId" v:"required#设备码不能为空" dc:"设备码(deviceId)"`
}

type CreateSimulatorDeviceWhitelistRes struct {
	ID string `json:"id"`
}

type UpdateSimulatorDeviceWhitelistReq struct {
	g.Meta   `path:"/updateSimulatorDeviceWhitelist" method:"post" summary:"修改模拟器设备白名单" tags:"模拟器设备白名单"`
	ID       uint64 `json:"id" v:"required#ID不能为空"`
	DeviceId string `json:"deviceId" v:"required#设备码不能为空" dc:"设备码(deviceId)"`
}

type UpdateSimulatorDeviceWhitelistRes struct {
	Success bool `json:"success"`
}

type DeleteSimulatorDeviceWhitelistReq struct {
	g.Meta `path:"/deleteSimulatorDeviceWhitelist" method:"post" summary:"删除模拟器设备白名单" tags:"模拟器设备白名单"`
	ID     uint64 `json:"id" v:"required#ID不能为空"`
}

type DeleteSimulatorDeviceWhitelistRes struct {
	Success bool `json:"success"`
}

package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/user"
)

type H5DeviceLoginReq struct {
	g.Meta     `path:"/h5DeviceLogin" method:"post" summary:"H5设备码快捷登录" tags:"权限"`
	DeviceInfo *entity.DeviceInfo `json:"deviceInfo" v:"required#设备信息不能为空" dc:"设备信息(需包含deviceId,不校验CPU型号)"`
}

type H5DeviceLoginRes = DeviceLoginRes

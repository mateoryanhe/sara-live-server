package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/user"
)

type DeviceLoginReq struct {
	g.Meta     `path:"/deviceLogin" method:"post" summary:"设备码快捷登录" tags:"权限"`
	DeviceInfo *entity.DeviceInfo `json:"deviceInfo" v:"required#设备信息不能为空" dc:"设备信息(需包含deviceId)"`
}

type DeviceLoginRes struct {
	Token     string `json:"token"`
	IsNewUser bool   `json:"isNewUser" dc:"是否首次注册"`
}

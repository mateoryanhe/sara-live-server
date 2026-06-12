package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

type PhoneLoginReq struct {
	g.Meta     `path:"/phoneLogin" method:"post" summary:"手机号登录" tags:"权限"`
	Phone      string             `json:"phone" summary:"手机号"`
	Password   string             `json:"password" summary:"密码"`
	DeviceInfo *entity.DeviceInfo `json:"deviceInfo"  dc:"设备信息"`
}

type PhoneLoginRes struct {
	Token string `json:"token"`
}

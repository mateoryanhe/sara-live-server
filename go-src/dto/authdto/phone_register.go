package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/user"
)

type PhoneRegisterReq struct {
	g.Meta        `path:"/phoneRegister" method:"post" summary:"手机号注册" tags:"权限"`
	PhoneAreaCode string             `json:"phoneAreaCode" v:"required" summary:"手机区号"`
	Phone         string             `json:"phone" v:"required" summary:"手机号"`
	Code          string             `json:"code" v:"required" summary:"验证码"`
	Password      string             `json:"password" v:"required|length:6,32" summary:"密码"`
	InviteCode    string             `json:"inviteCode" summary:"邀请码(邀请人的分享码)"`
	DeviceInfo    *entity.DeviceInfo `json:"deviceInfo"  dc:"设备信息"`
}

type PhoneRegisterRes struct {
	Token string `json:"token"`
}

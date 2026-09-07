package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/user"
)

// EmailLoginReq 邮箱验证码登录(不存在则自动注册)
type EmailLoginReq struct {
	g.Meta     `path:"/emailLogin" method:"post" summary:"邮箱验证码登录" tags:"权限"`
	Email      string             `json:"email" v:"required|email#邮箱不能为空|邮箱格式不正确" dc:"邮箱地址"`
	Code       string             `json:"code" v:"required|length:6,6#验证码不能为空|验证码为6位" dc:"邮箱验证码"`
	DeviceInfo *entity.DeviceInfo `json:"deviceInfo" dc:"设备信息(可选)"`
}

// EmailLoginRes 邮箱验证码登录响应
type EmailLoginRes struct {
	Token     string `json:"token"`
	IsNewUser bool   `json:"isNewUser" dc:"是否首次注册"`
}

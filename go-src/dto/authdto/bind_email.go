package authdto

import "github.com/gogf/gf/v2/frame/g"

// BindEmailReq 登录用户绑定邮箱(需验证码);邮箱须未被任何未注销账号占用
type BindEmailReq struct {
	g.Meta `path:"/bindEmail" method:"post" summary:"绑定邮箱" tags:"权限"`
	Email  string `json:"email" v:"required|email#邮箱不能为空|邮箱格式不正确" dc:"邮箱地址"`
	Code   string `json:"code" v:"required|length:6,6#验证码不能为空|验证码为6位" dc:"邮箱验证码"`
}

// BindEmailRes 绑定邮箱响应
type BindEmailRes struct {
	Success bool `json:"success"`
}

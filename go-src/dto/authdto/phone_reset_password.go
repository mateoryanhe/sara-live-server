package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PhoneResetPasswordReq struct {
	g.Meta        `path:"/phoneResetPassword" method:"post" summary:"手机号忘记密码" tags:"权限"`
	PhoneAreaCode string `json:"phoneAreaCode" v:"required" summary:"手机区号"`
	Phone         string `json:"phone" v:"required" summary:"手机号"`
	Code          string `json:"code" v:"required" summary:"验证码"`
	Password      string `json:"password" v:"required|length:6,32" summary:"新密码"`
}

type PhoneResetPasswordRes struct {
	Success bool `json:"success"`
}

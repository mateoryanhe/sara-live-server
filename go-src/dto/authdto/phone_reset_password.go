package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PhoneResetPasswordReq struct {
	g.Meta   `path:"/phoneResetPassword" method:"post" summary:"手机号忘记密码" tags:"权限"`
	Phone    string `json:"phone" summary:"手机号"`
	Code     string `json:"code" summary:"验证码"`
	Password string `json:"password" summary:"新密码"`
}

type PhoneResetPasswordRes struct {
	Success bool `json:"success"`
}

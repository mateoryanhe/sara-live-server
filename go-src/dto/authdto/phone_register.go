package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PhoneRegisterReq struct {
	g.Meta   `path:"/phoneRegister" method:"post" summary:"手机号注册" tags:"权限"`
	Phone    string `json:"phone" summary:"手机号"`
	Code     string `json:"code" summary:"验证码"`
	Password string `json:"password" summary:"密码"`
}

type PhoneRegisterRes struct {
	AuthId string `json:"authId"`
	Token  string `json:"token"`
}

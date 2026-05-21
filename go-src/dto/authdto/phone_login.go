package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type PhoneLoginReq struct {
	g.Meta   `path:"/phoneLogin" method:"post" summary:"手机号登录" tags:"权限"`
	Phone    string `json:"phone" summary:"手机号"`
	Password string `json:"password" summary:"密码"`
}

type PhoneLoginRes struct {
	AuthId string `json:"authId"`
	Token  string `json:"token"`
}

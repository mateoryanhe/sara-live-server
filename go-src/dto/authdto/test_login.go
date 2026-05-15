package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type TestLoginReq struct {
	g.Meta `path:"/testLogin" method:"post" summary:"测试渠道登录" tags:"权限"`
	OpenId string `json:"openId" summary:"OpenId"`
}

type TestLoginRes struct {
	AuthId string `json:"authId"`
	Token  string `json:"token"`
}

package verificationcodedto

import "github.com/gogf/gf/v2/frame/g"

// SendCodeReq 发送验证码请求
type SendCodeReq struct {
	g.Meta `path:"/sendCode" method:"post" summary:"发送验证码" tags:"验证码"`
	Phone  string `json:"phone" v:"required|phone" dc:"手机号"`
	IP     string `json:"ip"  dc:"客户端IP"`
}

// SendCodeRes 发送验证码响应
type SendCodeRes struct {
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"消息"`
}

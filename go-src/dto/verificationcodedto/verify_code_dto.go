package verificationcodedto

import "github.com/gogf/gf/v2/frame/g"

// VerifyCodeReq 验证验证码请求
type VerifyCodeReq struct {
	g.Meta `path:"/verifyCode" method:"post" summary:"验证验证码" tags:"验证码"`
	Phone  string `json:"phone" v:"required|phone" dc:"手机号"`
	Code   string `json:"code" v:"required" dc:"验证码"`
}

// VerifyCodeRes 验证验证码响应
type VerifyCodeRes struct {
	Success bool   `json:"success" dc:"是否成功"`
	Message string `json:"message" dc:"消息"`
}

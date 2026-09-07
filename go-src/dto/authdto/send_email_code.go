package authdto

import "github.com/gogf/gf/v2/frame/g"

// SendEmailCodeReq App 发送邮箱验证码(无需登录)
type SendEmailCodeReq struct {
	g.Meta `path:"/sendEmailCode" method:"post" summary:"发送邮箱验证码" tags:"验证码"`
	Email  string `json:"email" v:"required|email#邮箱不能为空|邮箱格式不正确" dc:"邮箱地址"`
	Lang   string `json:"lang" dc:"邮件语言(en/es/hi/pt/id);空则用 Accept-Language"`
}

// SendEmailCodeRes 发送邮箱验证码响应
type SendEmailCodeRes struct {
	Success bool `json:"success" dc:"是否成功"`
}

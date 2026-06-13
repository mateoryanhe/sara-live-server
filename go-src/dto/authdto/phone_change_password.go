package authdto

import "github.com/gogf/gf/v2/frame/g"

type PhoneChangePasswordReq struct {
	g.Meta      `path:"/phoneChangePassword" method:"post" summary:"修改登录密码(需旧密码)" tags:"权限"`
	OldPassword string `json:"oldPassword" v:"required#旧密码不能为空" dc:"旧密码"`
	NewPassword string `json:"newPassword" v:"required#新密码不能为空" dc:"新密码"`
}

type PhoneChangePasswordRes struct {
	Success bool `json:"success"`
}

package userinfodto

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GetUserInfoReq 查询当前登录用户的基础信息
type GetUserInfoReq struct {
	g.Meta `path:"/get" method:"post" summary:"获取用户基础信息" tags:"用户基础信息"`
}

type GetUserInfoRes struct {
	UserId   uint64 `json:"userId"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Remark   string `json:"remark"`
}

// UpdateNicknameReq 修改昵称
type UpdateNicknameReq struct {
	g.Meta   `path:"/updateNickname" method:"post" summary:"修改用户昵称" tags:"用户基础信息"`
	Nickname string `json:"nickname" v:"required#昵称不能为空" dc:"用户昵称"`
}

type UpdateNicknameRes struct {
	Nickname string `json:"nickname"`
}

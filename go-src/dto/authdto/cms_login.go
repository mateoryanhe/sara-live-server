package authdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

type CMSLoginReq struct {
	g.Meta   `path:"/cmsLogin" method:"post" summary:"cms登录" tags:"权限"`
	UserName string `json:"userName" summary:"用户名"`
	Pwd      string `json:"pwd" summary:"密码"`
}

type CMSLoginRes struct {
	AuthId  uint64               `json:"authId"`
	Token   string               `json:"token"`
	Admin   bool                 `json:"admin" summary:"超级管理员"`
	Modules []*entity.Permission `json:"modules" summary:"模块列表"`
}

func NewCMSLoginRes(userId uint64, token string, admin bool, module []*entity.Permission) *CMSLoginRes {
	return &CMSLoginRes{
		AuthId:  userId,
		Token:   token,
		Admin:   admin,
		Modules: module,
	}
}

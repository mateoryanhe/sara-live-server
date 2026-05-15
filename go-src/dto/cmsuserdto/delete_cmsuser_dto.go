package cmsuserdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type DeleteCMSUserReq struct {
	g.Meta `path:"/deleteCMSUser" method:"post" summary:"删除CMS用户" tags:"CMS用户管理"`
	ID     uint64 `json:"id" v:"required#CMS用户ID不能为空" dc:"CMS用户ID"`
}

type DeleteCMSUserRes struct {
	Success bool `json:"success"`
}

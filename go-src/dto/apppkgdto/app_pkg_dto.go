package apppkgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// AppPkgListReq CMS分页查询App包配置(内存)
type AppPkgListReq struct {
	g.Meta `path:"/appPkgList" method:"post" summary:"获取App包配置列表" tags:"App包管理"`
	httpserver.CMSQueryReq
	PackageName string `json:"packageName" dc:"包名(模糊匹配)"`
}

// AppPkgListRes 列表项
type AppPkgListRes struct {
	ID          string `json:"id"`
	PackageName string `json:"packageName"`
	SecretKey   string `json:"secretKey"`
	Remark      string `json:"remark"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateAppPkgReq 创建App包配置
type CreateAppPkgReq struct {
	g.Meta      `path:"/createAppPkg" method:"post" summary:"创建App包配置" tags:"App包管理"`
	PackageName string `json:"packageName" v:"required#包名不能为空"`
	SecretKey   string `json:"secretKey" v:"required#密钥不能为空"`
	Remark      string `json:"remark"`
}

type CreateAppPkgRes struct {
	ID string `json:"id"`
}

// UpdateAppPkgReq 更新App包配置
type UpdateAppPkgReq struct {
	g.Meta      `path:"/updateAppPkg" method:"post" summary:"修改App包配置" tags:"App包管理"`
	ID          uint64 `json:"id,string" v:"required#ID不能为空"`
	PackageName string `json:"packageName" v:"required#包名不能为空"`
	SecretKey   string `json:"secretKey" v:"required#密钥不能为空"`
	Remark      string `json:"remark"`
}

type UpdateAppPkgRes struct {
	Success bool `json:"success"`
}

// DeleteAppPkgReq 删除App包配置
type DeleteAppPkgReq struct {
	g.Meta `path:"/deleteAppPkg" method:"post" summary:"删除App包配置" tags:"App包管理"`
	ID     uint64 `json:"id,string" v:"required#ID不能为空"`
}

type DeleteAppPkgRes struct {
	Success bool `json:"success"`
}

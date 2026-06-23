package shortvideodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type ShortVideoCategoryListReq struct {
	g.Meta `path:"/shortVideoCategoryList" method:"post" summary:"获取短视频分类列表" tags:"短视频分类"`
	httpserver.CMSQueryReq
}

type ShortVideoCategoryListRes struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Sort      int    `json:"sort"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateShortVideoCategoryReq struct {
	g.Meta `path:"/createShortVideoCategory" method:"post" summary:"创建短视频分类" tags:"短视频分类"`
	Name   string `json:"name" v:"required|length:1,64#分类名称不能为空|分类名称长度需在1-64字符" dc:"分类名称"`
	Sort   int    `json:"sort" dc:"排序值(越大越靠前)"`
}

type CreateShortVideoCategoryRes struct {
	ID string `json:"id" dc:"分类ID"`
}

type UpdateShortVideoCategoryReq struct {
	g.Meta `path:"/updateShortVideoCategory" method:"post" summary:"修改短视频分类" tags:"短视频分类"`
	ID     uint64 `json:"id" v:"required#分类ID不能为空" dc:"分类ID"`
	Name   string `json:"name" v:"required|length:1,64#分类名称不能为空|分类名称长度需在1-64字符" dc:"分类名称"`
	Sort   int    `json:"sort" dc:"排序值(越大越靠前)"`
}

type UpdateShortVideoCategoryRes struct {
	Success bool `json:"success"`
}

type DeleteShortVideoCategoryReq struct {
	g.Meta `path:"/deleteShortVideoCategory" method:"post" summary:"删除短视频分类" tags:"短视频分类"`
	ID     uint64 `json:"id" v:"required#分类ID不能为空" dc:"分类ID"`
}

type DeleteShortVideoCategoryRes struct {
	Success bool `json:"success"`
}

type AppShortVideoCategoryListReq struct {
	g.Meta `path:"/appShortVideoCategoryList" method:"post" summary:"App查询短视频分类列表" tags:"短视频"`
}

type AppShortVideoCategoryItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Sort int    `json:"sort"`
}

type AppShortVideoCategoryListRes struct {
	List []*AppShortVideoCategoryItem `json:"list"`
}

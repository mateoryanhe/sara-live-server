package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

// AppShortVideoListReq App端分页查询短视频列表(仅已上架,按点赞数排序,走缓存)
type AppShortVideoListReq struct {
	g.Meta   `path:"/appShortVideoList" method:"post" summary:"App分页查询短视频列表(已上架)" tags:"短视频"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// AppShortVideoItem App端短视频列表元素
type AppShortVideoItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Video       string `json:"video" dc:"视频完整URL"`
	Cover       string `json:"cover" dc:"封面完整URL"`
	Description string `json:"description"`
	IsPaid      uint8  `json:"isPaid" dc:"是否付费(0免费,1付费)"`
	LikeCount   uint64 `json:"likeCount"`
}

// AppShortVideoListRes App端短视频分页列表响应
type AppShortVideoListRes struct {
	Total    int                  `json:"total" dc:"总条数"`
	Page     int                  `json:"page" dc:"当前页码"`
	PageSize int                  `json:"pageSize" dc:"每页数量"`
	List     []*AppShortVideoItem `json:"list" dc:"短视频列表"`
}

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
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Video            string  `json:"video" dc:"视频完整URL"`
	Cover            string  `json:"cover" dc:"封面完整URL"`
	IsPaid           uint8   `json:"isPaid" dc:"是否付费(0免费,1付费)"`
	DiamondPerMinute float64 `json:"diamondPerMinute" dc:"每分钟钻石数"`
	CategoryId       int     `json:"categoryId" dc:"视频分类ID"`
	Source           uint8   `json:"source" dc:"视频来源(1原创,2转发,3AI生成)"`
	AuthorId         string  `json:"authorId" dc:"作者用户ID"`
	AuthorNickname   string  `json:"authorNickname" dc:"作者昵称"`
	LikeCount        uint64  `json:"likeCount"`
	ViewCount        uint64  `json:"viewCount" dc:"观看人数"`
}

// AppShortVideoViewListReq App端分页查询短视频列表(仅已上架,按观看人数排序,走缓存)
type AppShortVideoViewListReq struct {
	g.Meta   `path:"/appShortVideoViewList" method:"post" summary:"App分页查询短视频列表(已上架,按观看人数排序)" tags:"短视频"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// AppShortVideoPublishListReq App端分页查询短视频列表(仅已上架,按发布时间降序,走缓存)
type AppShortVideoPublishListReq struct {
	g.Meta   `path:"/appShortVideoPublishList" method:"post" summary:"App分页查询短视频列表(已上架,按发布时间降序)" tags:"短视频"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// AppShortVideoListRes App端短视频分页列表响应
type AppShortVideoListRes struct {
	Total    int                  `json:"total" dc:"总条数"`
	Page     int                  `json:"page" dc:"当前页码"`
	PageSize int                  `json:"pageSize" dc:"每页数量"`
	List     []*AppShortVideoItem `json:"list" dc:"短视频列表"`
}

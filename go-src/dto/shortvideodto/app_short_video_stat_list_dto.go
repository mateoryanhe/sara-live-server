package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

// AppShortVideoStatListReq App端分页查询本人短视频统计数据(直接查库,按发布时间降序)
type AppShortVideoStatListReq struct {
	g.Meta   `path:"/appShortVideoStatList" method:"post" summary:"App分页查询本人短视频统计数据(按发布时间降序)" tags:"短视频"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// AppShortVideoStatItem App端短视频统计项
type AppShortVideoStatItem struct {
	VideoId            string  `json:"videoId" dc:"短视频ID"`
	Title              string  `json:"title" dc:"视频标题"`
	LikeCount          uint64  `json:"likeCount" dc:"点赞累计数量"`
	ViewCount          uint64  `json:"viewCount" dc:"观看人数(去重)"`
	WatchCount         uint64  `json:"watchCount" dc:"观看次数(累计)"`
	TotalDiamondIncome float64 `json:"totalDiamondIncome" dc:"累计钻石收益"`
	PublishedAt        int64   `json:"publishedAt" dc:"发布时间(毫秒时间戳)"`
}

// AppShortVideoStatListRes App端短视频统计分页响应
type AppShortVideoStatListRes struct {
	Page     int                      `json:"page" dc:"当前页码"`
	PageSize int                      `json:"pageSize" dc:"每页数量"`
	List     []*AppShortVideoStatItem `json:"list" dc:"统计列表"`
}

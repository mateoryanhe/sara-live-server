package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

// AppShortVideoWatchListReq App端分页查询短视频观看记录
type AppShortVideoWatchListReq struct {
	g.Meta   `path:"/appShortVideoWatchList" method:"post" summary:"App分页查询短视频观看记录" tags:"短视频"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// AppShortVideoWatchListRes App端观看记录列表响应(不含 total)
type AppShortVideoWatchListRes struct {
	Page     int                  `json:"page" dc:"当前页码"`
	PageSize int                  `json:"pageSize" dc:"每页数量"`
	List     []*AppShortVideoItem `json:"list" dc:"短视频列表"`
}

package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

// AppShortVideoWatchListReq App端分页查询短视频观看记录
type AppShortVideoWatchListReq struct {
	g.Meta   `path:"/appShortVideoWatchList" method:"post" summary:"App分页查询短视频观看记录" tags:"短视频"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// AppShortVideoWatchListItem App端短视频观看记录
type AppShortVideoWatchListItem struct {
	ID             string `json:"id"`
	VideoId        string `json:"videoId"`
	Title          string `json:"title"`
	Cover          string `json:"cover" dc:"封面完整URL"`
	AuthorId       string `json:"authorId" dc:"作者用户ID"`
	AuthorNickname string `json:"authorNickname" dc:"作者昵称"`
	AuthorAvatar   string `json:"authorAvatar" dc:"作者头像URL"`
	BilledSeconds  uint64 `json:"billedSeconds" dc:"已累计计费观看秒数"`
	WatchSeconds   uint64 `json:"watchSeconds" dc:"累计观看时长(秒,含免费与付费)"`
	UpdatedAt      string `json:"updatedAt" dc:"最近观看时间"`
}

// AppShortVideoWatchListRes App端短视频观看记录分页响应
type AppShortVideoWatchListRes struct {
	Total    int                           `json:"total" dc:"总条数"`
	Page     int                           `json:"page" dc:"当前页码"`
	PageSize int                           `json:"pageSize" dc:"每页数量"`
	List     []*AppShortVideoWatchListItem `json:"list" dc:"观看记录列表"`
}

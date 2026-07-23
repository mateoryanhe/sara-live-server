package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

// AppShortVideoListReq App端分页查询短视频列表(仅已上架,按点赞数排序,走缓存)
type AppShortVideoListReq struct {
	g.Meta   `path:"/appShortVideoList" method:"post" summary:"App端分页查询短视频列表(仅已上架,按点赞数排序,走内存缓存)" tags:"短视频"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

const (
	AppShortVideoScrollNext = 1  // 向下(下一个)
	AppShortVideoScrollPrev = -1 // 向上(上一个)

	AppShortVideoSortLike    = 1 // 按点赞数(默认)
	AppShortVideoSortView    = 2 // 按观看人数
	AppShortVideoSortPublish = 3 // 按发布时间
)

// AppShortVideoScrollReq 以当前观看视频为锚点,向上/向下拉取n个短视频
type AppShortVideoScrollReq struct {
	g.Meta    `path:"/appShortVideoScroll" method:"post" summary:"App以当前视频为锚点拉取上下n个短视频" tags:"短视频"`
	VideoId   string `json:"videoId" v:"required#当前视频ID不能为空" dc:"当前观看的视频ID"`
	Direction int    `json:"direction" v:"required|in:1,-1#滑动方向不能为空|direction取值无效(1向下下一个,-1向上上一个)" dc:"滑动方向(1向下下一个,-1向上上一个)"`
	Count     int    `json:"count" dc:"拉取数量(默认10,最大50)"`
	SortType  int    `json:"sortType" v:"in:0,1,2,3#排序类型无效" dc:"列表排序(0或1点赞数,2观看人数,3发布时间)"`
}

// AppShortVideoScrollRes App端短视频滑动拉取响应
type AppShortVideoScrollRes struct {
	List    []*AppShortVideoItem `json:"list" dc:"短视频列表(不含锚点视频本身)"`
	HasMore bool                 `json:"hasMore" dc:"该方向是否还有更多视频"`
}

// AppShortVideoItem App端短视频列表元素
type AppShortVideoItem struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Video            string  `json:"video" dc:"视频完整URL"`
	Cover            string  `json:"cover" dc:"封面完整URL"`
	IsPaid           uint8   `json:"isPaid" dc:"是否付费(0免费,1付费)"`
	PayDiamond       float64 `json:"payDiamond" dc:"付费钻石(一次性)"`
	CategoryId       int     `json:"categoryId" dc:"视频分类ID"`
	Source           uint8   `json:"source" dc:"视频来源(1原创,2转发,3AI生成)"`
	AuthorId         string  `json:"authorId" dc:"作者用户ID"`
	AuthorType       uint8   `json:"authorType" dc:"作者类型(0App用户,1CMS用户)"`
	AuthorNickname   string  `json:"authorNickname" dc:"作者昵称"`
	AuthorAvatar     string  `json:"authorAvatar" dc:"作者头像URL"`
	LikeCount        uint64  `json:"likeCount"`
	ViewCount        uint64  `json:"viewCount" dc:"观看人数(去重)"`
	WatchCount       uint64  `json:"watchCount" dc:"观看次数(累计)"`
	Duration         uint32  `json:"duration" dc:"视频时长(秒)"`
	FreeWatchSeconds uint32  `json:"freeWatchSeconds" dc:"免费观看时长(秒)"`
	FreeTime         uint64  `json:"freeTime" dc:"剩余免费观看时长(秒,来自观看记录)"`
	HasPaid          bool    `json:"hasPaid" dc:"当前用户是否已付费观看"`
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

// AppShortVideoUploadRecordListReq App端分页查询当前用户短视频上传记录(走内存缓存,按创建时间降序)
type AppShortVideoUploadRecordListReq struct {
	g.Meta   `path:"/appShortVideoUploadRecordList" method:"post" summary:"App分页查询短视频上传记录" tags:"短视频"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// AppShortVideoPendingReviewListReq App端分页查询本人发布的全部短视频(审核中优先)
type AppShortVideoPendingReviewListReq struct {
	g.Meta   `path:"/appShortVideoPendingReviewList" method:"post" summary:"App分页查询本人发布的短视频(审核中优先)" tags:"短视频"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// AppShortVideoUploadRecordItem App端短视频上传记录
type AppShortVideoUploadRecordItem struct {
	ID                 string  `json:"id"`
	Title              string  `json:"title"`
	Video              string  `json:"video" dc:"视频完整URL"`
	Cover              string  `json:"cover" dc:"封面完整URL"`
	Status             uint8   `json:"status" dc:"状态(0下架,1上架)"`
	CategoryId         int     `json:"categoryId" dc:"视频分类ID"`
	Source             uint8   `json:"source" dc:"视频来源(1原创,2转发,3AI生成)"`
	AuthorId           string  `json:"authorId" dc:"作者用户ID"`
	AuthorNickname     string  `json:"authorNickname" dc:"作者昵称"`
	AuthorAvatar       string  `json:"authorAvatar" dc:"作者头像URL"`
	LikeCount          uint64  `json:"likeCount"`
	ViewCount          uint64  `json:"viewCount" dc:"观看人数(去重)"`
	WatchCount         uint64  `json:"watchCount" dc:"观看次数(累计)"`
	Duration           uint32  `json:"duration" dc:"视频时长(秒)"`
	TotalDiamondIncome float64 `json:"totalDiamondIncome" dc:"累计钻石收益"`
	CreatedAt          string  `json:"createdAt" dc:"上传时间"`
	UpdatedAt          string  `json:"updatedAt" dc:"审核时间"`
}

// AppShortVideoUploadRecordListRes App端短视频上传记录分页响应
type AppShortVideoUploadRecordListRes struct {
	Total    int                              `json:"total" dc:"总条数"`
	Page     int                              `json:"page" dc:"当前页码"`
	PageSize int                              `json:"pageSize" dc:"每页数量"`
	List     []*AppShortVideoUploadRecordItem `json:"list" dc:"上传记录列表"`
}

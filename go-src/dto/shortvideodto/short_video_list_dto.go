package shortvideodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type ShortVideoListReq struct {
	g.Meta `path:"/shortVideoList" method:"post" summary:"获取短视频列表" tags:"短视频"`
	httpserver.CMSQueryReq
	Title          string `json:"title" dc:"标题(模糊匹配)"`
	AuthorNickname string `json:"authorNickname" dc:"作者昵称(模糊匹配)"`
	AuthorId       string `json:"authorId" dc:"作者用户ID(精确匹配,可选)"`
	StatusFilter   int    `json:"statusFilter" dc:"状态过滤(0=全部, 1=只看下架, 2=只看上架)"`
	SortField      string `json:"sortField" dc:"排序字段(空=创建时间倒序, viewCount=观看人数升序, totalDiamondIncome=钻石收益升序)"`
}

type ShortVideoListRes struct {
	ID                 string  `json:"id"`
	Title              string  `json:"title"`
	Video              string  `json:"video" dc:"视频完整URL(列表展示)"`
	VideoName          string  `json:"videoName" dc:"视频资源文件名(编辑保存用)"`
	Cover              string  `json:"cover" dc:"封面完整URL(列表展示)"`
	CoverName          string  `json:"coverName" dc:"封面资源文件名(编辑保存用)"`
	Sort               int     `json:"sort"`
	Status             uint8   `json:"status"`
	IsPaid             uint8   `json:"isPaid"`
	PayDiamond         float64 `json:"payDiamond"`
	CategoryId         int     `json:"categoryId"`
	Source             uint8   `json:"source"`
	AuthorId           string  `json:"authorId"`
	AuthorType         uint8   `json:"authorType" dc:"作者类型(0App用户,1CMS用户,仅展示不可编辑)"`
	AuthorNickname     string  `json:"authorNickname"`
	LikeCount          uint64  `json:"likeCount" dc:"点赞数"`
	ViewCount          uint64  `json:"viewCount" dc:"观看人数(去重)"`
	WatchCount         uint64  `json:"watchCount" dc:"观看次数(累计)"`
	TotalDiamondIncome float64 `json:"totalDiamondIncome" dc:"累计钻石收益"`
	Duration           uint32  `json:"duration" dc:"视频时长(秒)"`
	FreeWatchSeconds   uint32  `json:"freeWatchSeconds" dc:"免费观看时长(秒)"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

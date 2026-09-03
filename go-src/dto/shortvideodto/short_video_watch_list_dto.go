package shortvideodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type ShortVideoWatchListReq struct {
	g.Meta `path:"/shortVideoWatchList" method:"post" summary:"分页查询短视频观看记录" tags:"短视频"`
	httpserver.CMSQueryReq
	UserId    string `json:"userId" dc:"用户ID(空=全部)"`
	StartTime int64  `json:"startTime" dc:"时间起(秒,0=不过滤; onlyPaid=true 时按付费时间,否则按更新时间)"`
	EndTime   int64  `json:"endTime" dc:"时间止(秒,0=不过滤; onlyPaid=true 时按付费时间,否则按更新时间)"`
	OnlyPaid  bool   `json:"onlyPaid" dc:"仅付费购买记录(paid_time 非空)"`
}

type ShortVideoWatchListItem struct {
	ID         string  `json:"id"`
	UserId     string  `json:"userId"`
	Nickname   string  `json:"nickname"`
	VideoId    string  `json:"videoId"`
	VideoTitle string  `json:"videoTitle"`
	PayDiamond float64 `json:"payDiamond" dc:"视频当前标价钻石(仅供参考)"`
	PaidTime   string  `json:"paidTime" dc:"付费时间(空表示未付费)"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

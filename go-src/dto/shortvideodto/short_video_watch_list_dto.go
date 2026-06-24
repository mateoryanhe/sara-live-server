package shortvideodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type ShortVideoWatchListReq struct {
	g.Meta `path:"/shortVideoWatchList" method:"post" summary:"分页查询短视频观看记录" tags:"短视频"`
	httpserver.CMSQueryReq
	UserId    string `json:"userId" dc:"用户ID(空=全部)"`
	StartTime int64  `json:"startTime" dc:"更新时间起(秒,0=不过滤)"`
	EndTime   int64  `json:"endTime" dc:"更新时间止(秒,0=不过滤)"`
}

type ShortVideoWatchListItem struct {
	ID            string `json:"id"`
	UserId        string `json:"userId"`
	Nickname      string `json:"nickname"`
	VideoId       string `json:"videoId"`
	VideoTitle    string `json:"videoTitle"`
	BilledSeconds uint64 `json:"billedSeconds,string"`
	WatchSeconds  uint64 `json:"watchSeconds,string"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

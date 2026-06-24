package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

// WatchShortVideoStartReq App端开始观看短视频
type WatchShortVideoStartReq struct {
	g.Meta  `path:"/watchShortVideoStart" method:"post" summary:"开始观看短视频" tags:"短视频"`
	VideoId uint64 `json:"videoId,string" v:"required#视频ID不能为空" dc:"短视频ID"`
}

type WatchShortVideoStartRes struct {
}

// WatchShortVideoEndReq App端结束观看短视频
type WatchShortVideoEndReq struct {
	g.Meta  `path:"/watchShortVideoEnd" method:"post" summary:"结束观看短视频" tags:"短视频"`
	VideoId uint64 `json:"videoId,string" v:"required#视频ID不能为空" dc:"短视频ID"`
}

type WatchShortVideoEndRes struct {
}

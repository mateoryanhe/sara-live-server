package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type LikeShortVideoReq struct {
	g.Meta  `path:"/likeShortVideo" method:"post" summary:"点赞短视频" tags:"短视频"`
	VideoId uint64 `json:"videoId" v:"required#视频ID不能为空" dc:"短视频ID"`
}

type LikeShortVideoRes struct {
	LikeCount uint64 `json:"likeCount" dc:"点赞累计数量"`
}

package liveroomdto

import "github.com/gogf/gf/v2/frame/g"

type GetLiveRoomShareTextReq struct {
	g.Meta `path:"/getLiveRoomShareText" method:"post" summary:"直播间分享文案" tags:"直播间"`
	RoomId uint64 `json:"roomId" v:"required|min:1#房间ID不能为空|房间ID非法" dc:"直播间ID"`
}

type GetLiveRoomShareTextRes struct {
	ShareText string `json:"shareText" dc:"分享文案(按请求语言)"`
	AfLink    string `json:"afLink" dc:"AppsFlyer 链接(已拼接 inviteCode、roomId)"`
	ShareCode string `json:"shareCode" dc:"两行:文案+\\n+AF链接"`
}
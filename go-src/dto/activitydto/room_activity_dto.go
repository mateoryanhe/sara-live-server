package activitydto

import "github.com/gogf/gf/v2/frame/g"

// GetRoomActivityReq App 按房间拉取活动入口配置.
type GetRoomActivityReq struct {
	g.Meta `path:"/getRoomActivity" method:"post" summary:"房间活动配置" tags:"活动"`
	RoomId uint64 `json:"roomId" v:"required|min:1#房间ID不能为空|房间ID无效" dc:"房间ID(直播间ID)"`
}

// GetRoomActivityRes 活动按钮/列表地址 + 推荐项 ID(优先房间推荐,否则列表首条).
type GetRoomActivityRes struct {
	ButtonIcon string `json:"buttonIcon" dc:"活动按钮 Icon URL"`
	ListUrl    string `json:"listUrl" dc:"活动列表 URL"`
	ItemId     string `json:"itemId" dc:"推荐项 ID(房间推荐第一条;无则活动列表第一条)"`
}

package calldto

import "github.com/gogf/gf/v2/frame/g"

// BatchInviteLiveRoomCallReq 主播批量邀请观众进行直播间通话
type BatchInviteLiveRoomCallReq struct {
	g.Meta  `path:"/batchInviteLiveRoomCall" method:"post" summary:"主播批量邀请直播间通话" tags:"通话"`
	UserIds []string `json:"userIds" v:"required|min-length:1#用户ID列表不能为空|用户ID列表不能为空" dc:"被邀请观众用户ID列表"`
}

// BatchInviteLiveRoomCallRes 主播批量邀请直播间通话响应
type BatchInviteLiveRoomCallRes struct {
	InvitedCount  int      `json:"invitedCount" dc:"成功推送邀请数"`
	SkippedUserIds []string `json:"skippedUserIds" dc:"跳过的用户ID(无效/自己/不在房/不存在)"`
}

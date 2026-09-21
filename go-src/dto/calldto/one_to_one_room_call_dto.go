package calldto

import "github.com/gogf/gf/v2/frame/g"

// OneToOneRoomCallReq 1v1房间视频通话呼叫。
type OneToOneRoomCallReq struct {
	g.Meta   `path:"/oneToOneRoomCall" method:"post" summary:"1v1房间视频通话呼叫" tags:"通话"`
	TargetId uint64 `json:"targetId" v:"required|min:1#目标用户ID不能为空|目标用户ID无效" dc:"目标用户ID(通话接收方)"`
}

// OneToOneRoomCallRes 1v1房间视频通话呼叫响应。
type OneToOneRoomCallRes = LiveRoomCallRes

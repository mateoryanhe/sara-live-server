package datasyncdto

import (
	"github.com/gogf/gf/v2/frame/g"
	liveentity "xr-game-server/entity/live"
)

type SyncVoiceChatMicSeatTierReq struct {
	g.Meta `path:"/syncVoiceChatMicSeatTier" method:"post" summary:"同步语聊麦位挡位到目标环境" tags:"数据同步"`
}

type ReceiveVoiceChatMicSeatTierReq struct {
	g.Meta `path:"/receiveVoiceChatMicSeatTier" method:"post" summary:"接收语聊麦位挡位同步" tags:"数据同步"`
	Rows   []*liveentity.VoiceChatMicSeatTier `json:"rows"`
}

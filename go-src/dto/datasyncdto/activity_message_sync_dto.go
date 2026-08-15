package datasyncdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/message"
)

type SyncActivityMessageReq struct {
	g.Meta `path:"/syncActivityMessage" method:"post" summary:"同步活动消息到目标环境" tags:"数据同步"`
	IDs    []uint64 `json:"ids" v:"required|min-length:1#请选择要同步的活动消息" dc:"要同步的活动消息ID列表"`
}

type SyncActivityMessageRes struct {
	Success   bool   `json:"success"`
	RowCount  int    `json:"rowCount"`
	FileCount int    `json:"fileCount"`
	Message   string `json:"message"`
}

type ReceiveActivityMessageReq struct {
	g.Meta `path:"/receiveActivityMessage" method:"post" summary:"接收活动消息同步" tags:"数据同步"`
	Rows   []*entity.ActivityMessage `json:"rows"`
	Files  []*SyncFileItem           `json:"files"`
}

type ReceiveActivityMessageRes struct {
	Success   bool `json:"success"`
	RowCount  int  `json:"rowCount"`
	FileCount int  `json:"fileCount"`
}

package liveroomdto

import "github.com/gogf/gf/v2/frame/g"

// BatchAnchorStatusReq 批量查询主播状态(开播/WebSocket在线等)
type BatchAnchorStatusReq struct {
	g.Meta     `path:"/batchAnchorStatus" method:"post" summary:"批量查询主播状态" tags:"直播间"`
	AnchorIds  []uint64 `json:"anchorIds,string" v:"required|min-length:1#主播ID列表不能为空" dc:"主播ID列表(最多100个)"`
}

// AnchorStatusItem 单个主播状态
type AnchorStatusItem struct {
	AnchorId     string `json:"anchorId" dc:"主播用户ID(同直播间ID)"`
	Status       uint8  `json:"status" dc:"开播状态(0=未开播/已下播,1=直播中;与roomList.status一致)"`
	OnlineStatus uint8  `json:"onlineStatus" dc:"主播WebSocket在线(1=在线,2=离线;与oneToOneRoomList一致)"`
}

// BatchAnchorStatusRes 批量主播状态
type BatchAnchorStatusRes struct {
	List []*AnchorStatusItem `json:"list" dc:"与请求anchorIds顺序一致(去重后)"`
}
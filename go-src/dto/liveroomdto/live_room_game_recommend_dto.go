package liveroomdto

import "github.com/gogf/gf/v2/frame/g"

// GetLiveRoomGameRecommendListReq App 查询直播间推荐游戏列表
type GetLiveRoomGameRecommendListReq struct {
	g.Meta     `path:"/gameRecommendList" method:"post" summary:"查询直播间推荐游戏列表" tags:"直播间"`
	LiveRoomId uint64 `json:"liveRoomId" v:"required|min:1#直播间ID不能为空|直播间ID无效" dc:"直播间ID"`
}

// LiveRoomGameRecommendItem 直播间推荐游戏项
type LiveRoomGameRecommendItem struct {
	GameCode string `json:"gameCode" dc:"游戏编码"`
	Name     string `json:"name" dc:"名称(第三方)"`
	NameEn   string `json:"nameEn" dc:"英文名称(第三方)"`
	Cover    string `json:"cover" dc:"封面完整URL(第三方)"`
	Category string `json:"category" dc:"分类(第三方)"`
	Platform string `json:"platform" dc:"平台(第三方)"`
}

// GetLiveRoomGameRecommendListRes App 直播间推荐游戏列表响应
type GetLiveRoomGameRecommendListRes struct {
	List []*LiveRoomGameRecommendItem `json:"list" dc:"推荐游戏列表"`
}

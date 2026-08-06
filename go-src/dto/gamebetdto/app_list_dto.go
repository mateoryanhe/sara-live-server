package gamebetdto

import "github.com/gogf/gf/v2/frame/g"

// AppGameBetListReq App 分页查询游戏下注记录
type AppGameBetListReq struct {
	g.Meta   `path:"/appGameBetList" method:"post" summary:"App分页查询游戏下注记录" tags:"游戏"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"忽略,后端固定50"`
}

// AppGameBetListItem App 游戏下注记录项
type AppGameBetListItem struct {
	Id           string  `json:"id" dc:"记录ID"`
	GameCode     string  `json:"gameCode" dc:"游戏编码"`
	NameEn       string  `json:"nameEn" dc:"英文名称"`
	Cover        string  `json:"cover" dc:"封面完整URL"`
	Amount       float64 `json:"amount" dc:"下注金额"`
	PlatformType string  `json:"platformType" dc:"平台类型,枚举: ZY"`
	OrderId      string  `json:"orderId" dc:"订单ID"`
	LiveRoomId   string  `json:"liveRoomId" dc:"下注时所在直播间ID"`
	LiveRecordId string  `json:"liveRecordId" dc:"下注时所在直播记录ID"`
	CreatedAt    int64   `json:"createdAt" dc:"下注时间(毫秒时间戳)"`
}

// AppGameBetListRes App 游戏下注记录分页响应
type AppGameBetListRes struct {
	Page     int                   `json:"page" dc:"当前页码"`
	PageSize int                   `json:"pageSize" dc:"每页数量(固定50)"`
	List     []*AppGameBetListItem `json:"list" dc:"下注记录列表"`
}

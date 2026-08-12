package messagedto

import "github.com/gogf/gf/v2/frame/g"

// AppPersonalSystemMessageListReq App端查询个人系统消息列表
type AppPersonalSystemMessageListReq struct {
	g.Meta `path:"/personalSystemMessageList" method:"post" summary:"查询个人系统消息列表" tags:"系统消息"`
}

type AppPersonalSystemMessageItem struct {
	Id            uint64 `json:"id,string" dc:"个人系统消息ID"`
	MessageTypeId uint32 `json:"messageTypeId" dc:"消息类型ID"`
	Params        string `json:"params" dc:"扩展参数(JSON)"`
	IconEn        string `json:"iconEn" dc:"图标URL(英文)"`
	IconEs        string `json:"iconEs" dc:"图标URL(西班牙语)"`
	IconPt        string `json:"iconPt" dc:"图标URL(葡萄牙语)"`
	IconHi        string `json:"iconHi" dc:"图标URL(印地语)"`
	IconId        string `json:"iconId" dc:"图标URL(印尼语)"`
	TitleEn       string `json:"titleEn" dc:"标题(英文)"`
	TitleEs       string `json:"titleEs" dc:"标题(西班牙语)"`
	TitlePt       string `json:"titlePt" dc:"标题(葡萄牙语)"`
	TitleHi       string `json:"titleHi" dc:"标题(印地语)"`
	TitleId       string `json:"titleId" dc:"标题(印尼语)"`
	ContentEn     string `json:"contentEn" dc:"内容(英文)"`
	ContentEs     string `json:"contentEs" dc:"内容(西班牙语)"`
	ContentPt     string `json:"contentPt" dc:"内容(葡萄牙语)"`
	ContentHi     string `json:"contentHi" dc:"内容(印地语)"`
	ContentId     string `json:"contentId" dc:"内容(印尼语)"`
	CreatedAt     string `json:"createdAt" dc:"创建时间"`
}

type AppPersonalSystemMessageListRes struct {
	List []*AppPersonalSystemMessageItem `json:"list" dc:"个人系统消息列表"`
}

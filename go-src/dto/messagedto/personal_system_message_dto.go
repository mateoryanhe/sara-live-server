package messagedto

import "github.com/gogf/gf/v2/frame/g"

// AppPersonalSystemMessageListReq App端查询个人系统消息列表
type AppPersonalSystemMessageListReq struct {
	g.Meta `path:"/personalSystemMessageList" method:"post" summary:"查询个人系统消息列表" tags:"系统消息"`
}

type AppPersonalSystemMessageItem struct {
	Id        uint64 `json:"id,string" dc:"个人系统消息ID"`
	Icon      string `json:"icon" dc:"图标URL"`
	TitleEn   string `json:"titleEn" dc:"标题(英文)"`
	TitleEs   string `json:"titleEs" dc:"标题(西班牙语)"`
	TitlePt   string `json:"titlePt" dc:"标题(葡萄牙语)"`
	TitleHi   string `json:"titleHi" dc:"标题(印地语)"`
	ContentEn string `json:"contentEn" dc:"内容(英文)"`
	ContentEs string `json:"contentEs" dc:"内容(西班牙语)"`
	ContentPt string `json:"contentPt" dc:"内容(葡萄牙语)"`
	ContentHi string `json:"contentHi" dc:"内容(印地语)"`
	Params    string `json:"params" dc:"扩展参数(JSON)"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
}

type AppPersonalSystemMessageListRes struct {
	List []*AppPersonalSystemMessageItem `json:"list" dc:"个人系统消息列表"`
}

package messagedto

import "github.com/gogf/gf/v2/frame/g"

// AppPersonalSystemMessageListReq App端查询个人系统消息列表
type AppPersonalSystemMessageListReq struct {
	g.Meta `path:"/personalSystemMessageList" method:"post" summary:"查询个人系统消息列表" tags:"系统消息"`
}

type AppPersonalSystemMessageItem struct {
	Id        uint64 `json:"id,string"`
	Icon      string `json:"icon"`
	TitleEn   string `json:"titleEn"`
	TitleEs   string `json:"titleEs"`
	TitlePt   string `json:"titlePt"`
	TitleHi   string `json:"titleHi"`
	ContentEn string `json:"contentEn"`
	ContentEs string `json:"contentEs"`
	ContentPt string `json:"contentPt"`
	ContentHi string `json:"contentHi"`
	Params    string `json:"params"`
	CreatedAt string `json:"createdAt"`
}

type AppPersonalSystemMessageListRes struct {
	List []*AppPersonalSystemMessageItem `json:"list"`
}

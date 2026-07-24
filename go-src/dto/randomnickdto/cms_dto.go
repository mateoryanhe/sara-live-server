package randomnickdto

import "github.com/gogf/gf/v2/frame/g"

type GetRandomNicknameCfgReq struct {
	g.Meta `path:"/getRandomNicknameCfg" method:"post" summary:"查询随机昵称库概览" tags:"随机昵称"`
}

type RandomNicknameLangItem struct {
	Lang      uint8    `json:"lang"`
	LangCode  string   `json:"langCode"`
	LangLabel string   `json:"langLabel"`
	Count     int      `json:"count"`
	Samples   []string `json:"samples"`
}

type GetRandomNicknameCfgRes struct {
	UseDB bool                      `json:"useDB" dc:"是否使用数据库昵称库"`
	Langs []*RandomNicknameLangItem `json:"langs"`
}

type ImportRandomNicknamesReq struct {
	g.Meta  `path:"/importRandomNicknames" method:"post" summary:"批量导入随机昵称" tags:"随机昵称"`
	Lang    uint8  `json:"lang" v:"required|in:1,2,3,4#请选择语言|语言无效" dc:"语言(1英2西3印4葡)"`
	Content string `json:"content" v:"required#请输入昵称内容" dc:"昵称文本,一行一个"`
	Replace bool   `json:"replace" dc:"true=覆盖该语言,false=追加去重"`
}

type ImportRandomNicknamesRes struct {
	Imported int `json:"imported" dc:"本次导入条数"`
	Total    int `json:"total" dc:"该语言导入后总数(内存)"`
}

type ClearRandomNicknamesReq struct {
	g.Meta `path:"/clearRandomNicknames" method:"post" summary:"清空指定语言随机昵称" tags:"随机昵称"`
	Lang   uint8 `json:"lang" v:"required|in:1,2,3,4#请选择语言|语言无效"`
}

type ClearRandomNicknamesRes struct {
	Success bool `json:"success"`
	Total   int  `json:"total"`
}

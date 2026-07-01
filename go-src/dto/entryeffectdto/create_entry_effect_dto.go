package entryeffectdto

import "github.com/gogf/gf/v2/frame/g"

type CreateEntryEffectReq struct {
	g.Meta     `path:"/createEntryEffect" method:"post" summary:"创建进场特效" tags:"进场特效"`
	Name       string `json:"name"       v:"required|length:1,64#名称不能为空|名称长度需在1到64之间" dc:"名称"`
	LevelStart int    `json:"levelStart" v:"required|min:1#等级开始不能为空|等级开始必须大于0" dc:"等级开始"`
	LevelEnd   int    `json:"levelEnd"   v:"required|min:1#等级结束不能为空|等级结束必须大于0" dc:"等级结束"`
	Animation  string `json:"animation"  v:"required|max-length:255#动画资源不能为空|动画URL最长255字符" dc:"动画资源文件名"`
}

type CreateEntryEffectRes struct {
	ID string `json:"id" dc:"进场特效ID"`
}

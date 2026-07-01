package entryeffectdto

import "github.com/gogf/gf/v2/frame/g"

type DeleteEntryEffectReq struct {
	g.Meta `path:"/deleteEntryEffect" method:"post" summary:"删除进场特效" tags:"进场特效"`
	ID     uint64 `json:"id" v:"required#进场特效ID不能为空" dc:"进场特效ID"`
}

type DeleteEntryEffectRes struct {
	Success bool `json:"success"`
}

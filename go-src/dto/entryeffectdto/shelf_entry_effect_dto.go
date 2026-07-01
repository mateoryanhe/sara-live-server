package entryeffectdto

import "github.com/gogf/gf/v2/frame/g"

type OnShelfEntryEffectReq struct {
	g.Meta `path:"/onShelfEntryEffect" method:"post" summary:"上架进场特效" tags:"进场特效"`
	ID     uint64 `json:"id" v:"required#进场特效ID不能为空" dc:"进场特效ID"`
}

type OnShelfEntryEffectRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type OffShelfEntryEffectReq struct {
	g.Meta `path:"/offShelfEntryEffect" method:"post" summary:"下架进场特效" tags:"进场特效"`
	ID     uint64 `json:"id" v:"required#进场特效ID不能为空" dc:"进场特效ID"`
}

type OffShelfEntryEffectRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

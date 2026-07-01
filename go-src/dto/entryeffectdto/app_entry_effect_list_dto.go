package entryeffectdto

import "github.com/gogf/gf/v2/frame/g"

type AppEntryEffectListReq struct {
	g.Meta `path:"/appEntryEffectList" method:"post" summary:"App查询进场特效列表(已上架)" tags:"进场特效"`
}

type AppEntryEffectItem struct {
	ID         uint64 `json:"id,string"`
	Name       string `json:"name"`
	LevelStart int    `json:"levelStart"`
	LevelEnd   int    `json:"levelEnd"`
	Animation  string `json:"animation"`
}

type AppEntryEffectListRes struct {
	List []*AppEntryEffectItem `json:"list" dc:"进场特效列表"`
}

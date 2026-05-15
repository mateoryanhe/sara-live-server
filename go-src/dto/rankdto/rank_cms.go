package rankdto

import (
	"github.com/gogf/gf/v2/frame/g"
)

type AddRankDataReq struct {
	g.Meta `path:"/AddRankData" method:"post" tags:"排行榜" summary:"加事件"`
	TypeId uint32 `json:"typeId" v:"required#请选择排行榜"`
	Val    uint64 `json:"val" v:"required#请输入数值"`
}

type ReduceRankReq struct {
	g.Meta `path:"/ReduceRank" method:"post" tags:"排行榜" summary:"减事件"`
	TypeId uint32 `json:"typeId" v:"required#请选择排行榜"`
	Val    uint64 `json:"val" v:"required#请输入数值"`
}

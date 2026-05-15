package rankdto

import "github.com/gogf/gf/v2/frame/g"

type GetRankReq struct {
	g.Meta `path:"/getRank" method:"post" tags:"排行榜" summary:"获取排行榜数据"`
	TypeId uint32 `json:"typeId" v:"required#请选择排行榜"`
}

type GetRankRes struct {
	Data           []*RankDto `json:"data"`
	TypeId         uint32     `json:"typeId"`
	My             *RankDto   `json:"my"`
	SettlementTime string     `json:"settlementTime"`
}

func NewGetRankRes(typeId uint32) *GetRankRes {
	return &GetRankRes{
		TypeId: typeId,
		Data:   []*RankDto{},
		My:     nil,
	}
}

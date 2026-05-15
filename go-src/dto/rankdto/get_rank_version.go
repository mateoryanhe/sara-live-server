package rankdto

import "github.com/gogf/gf/v2/frame/g"

type GetRankVersionReq struct {
	g.Meta `path:"/getRankVersion" method:"post" tags:"排行榜" summary:"获取排行榜版本"`
	TypeId uint32 `json:"typeId" v:"required#请选择排行榜"`
}

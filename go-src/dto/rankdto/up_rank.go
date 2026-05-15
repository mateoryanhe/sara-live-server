package rankdto

import "github.com/gogf/gf/v2/frame/g"

// 上排行榜
type UpRankReq struct {
	g.Meta `path:"/UpRank" method:"post" tags:"排行榜" summary:"上排行榜"`
	RoleId uint64 `json:"roleId" v:"required#请选择角色" dc:"角色"`
}

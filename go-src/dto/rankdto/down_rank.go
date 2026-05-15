package rankdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

type DownRankReq struct {
	g.Meta   `path:"/DownRank" method:"post" tags:"排行榜" summary:"下排行榜"`
	RoleId   uint64     `json:"roleId" v:"required#请选择角色" dc:"角色"`
	LockTime *time.Time `json:"lockTime"  dc:"禁止时间"`
}

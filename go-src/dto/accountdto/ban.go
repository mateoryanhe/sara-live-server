package accountdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

type BanReq struct {
	g.Meta       `path:"/ban" method:"post" summary:"封号" tags:"账号"`
	AccountId    uint64     `json:"accountId" dc:"账号id"`
	BanApplyTime *time.Time `json:"banApplyTime" dc:"封禁时间"`
}

type UnBanReq struct {
	g.Meta    `path:"/unBan" method:"post" summary:"解封" tags:"账号"`
	AccountId uint64 `json:"accountId" dc:"账号id"`
}

type BanRes struct {
}

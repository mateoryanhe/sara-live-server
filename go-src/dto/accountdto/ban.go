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

// BanAnchorReq CMS封禁主播(含App推送,写入 live_rooms)
type BanAnchorReq struct {
	g.Meta       `path:"/banAnchor" method:"post" summary:"封禁主播" tags:"账号"`
	AccountId    uint64     `json:"accountId" v:"required#账号ID不能为空" dc:"主播账号id"`
	BanApplyTime *time.Time `json:"banApplyTime" v:"required#封禁截止时间不能为空" dc:"封禁截止时间"`
	BanReason    string     `json:"banReason" v:"required|length:1,512#封禁原因不能为空|封禁原因长度需在1到512之间" dc:"封禁原因"`
}

// UnBanAnchorReq CMS解封主播直播间
type UnBanAnchorReq struct {
	g.Meta    `path:"/unBanAnchor" method:"post" summary:"解封主播" tags:"账号"`
	AccountId uint64 `json:"accountId" v:"required#账号ID不能为空" dc:"主播账号id"`
}

type UnBanReq struct {
	g.Meta    `path:"/unBan" method:"post" summary:"解封" tags:"账号"`
	AccountId uint64 `json:"accountId" dc:"账号id"`
}

type BanRes struct {
}

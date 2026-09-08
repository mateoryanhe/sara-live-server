package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type GuildListReq struct {
	g.Meta `path:"/guildList" method:"post" summary:"获取直播工会列表" tags:"直播工会"`
	httpserver.CMSQueryReq
	Name string `json:"name" dc:"工会名称"`
}

// GuildListForVisibilityReq 可见性管理页拉全部上架工会(不按可见性表过滤)
type GuildListForVisibilityReq struct {
	g.Meta `path:"/guildListForVisibility" method:"post" summary:"可见性管理拉取全部上架工会" tags:"直播工会"`
	httpserver.CMSQueryReq
	Name string `json:"name" dc:"工会名称"`
}

type GuildListRes struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	LeaderId             string  `json:"leaderId"`
	LeaderName           string  `json:"leaderName"`
	Description          string  `json:"description"`
	Status               uint8   `json:"status"`
	GuildType            uint8   `json:"guildType" dc:"工会类型(0普通,1币商)"`
	SharePercent         float64 `json:"sharePercent" dc:"分佣比例(%)"`
	CreatorId            string  `json:"creatorId" dc:"创建者CMS用户ID"`
	CreatorName          string  `json:"creatorName" dc:"创建者CMS用户名"`
	UnsettledTotalIncome float64 `json:"unsettledTotalIncome" dc:"未结算工会总收益"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}

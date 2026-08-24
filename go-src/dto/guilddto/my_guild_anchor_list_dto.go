package guilddto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type GetMyGuildAnchorListReq struct {
	g.Meta `path:"/getMyGuildAnchorList" method:"post" summary:"获取当前CMS用户管理的工会名下主播列表(含未结算收益)" tags:"直播工会"`
	httpserver.CMSQueryReq
	GuildId uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
}

type GetMyOwnedGuildAnchorListReq struct {
	g.Meta `path:"/getMyOwnedGuildAnchorList" method:"post" summary:"获取当前CMS用户管理的全部工会名下主播列表" tags:"直播工会"`
	httpserver.CMSQueryReq
	Key string `json:"key" dc:"查询关键字(用户ID模糊/昵称/手机号)"`
}

type MyGuildAnchorListItem struct {
	ID                   string  `json:"id"`
	Nickname             string  `json:"nickname"`
	Avatar               string  `json:"avatar"`
	UnsettledTotalIncome float64 `json:"unsettledTotalIncome" dc:"未结算总收益"`
}

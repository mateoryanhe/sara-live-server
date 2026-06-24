package accountdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
	"xr-game-server/core/httpserver"
)

// QueryAnchorListReq CMS分页查询主播列表
type QueryAnchorListReq struct {
	g.Meta `path:"/getAnchorList" method:"post" summary:"获取主播列表" tags:"账号"`
	httpserver.CMSQueryReq
	Key string `json:"key" dc:"查询关键字(用户ID/昵称/手机号/分享码)"`
}

// AnchorListItem 主播列表项(基于 user_infos)
type AnchorListItem struct {
	ID                           uint64     `json:"id,string"`
	Nickname                     string     `json:"nickname"`
	Phone                        string     `json:"phone"`
	Avatar                       string     `json:"avatar"`
	GuildId                      uint64     `json:"guildId,string"`
	IP                           string     `json:"ip" dc:"登录IP"`
	RoomTitle                    string     `json:"roomTitle"`
	LiveStatus                   uint8      `json:"liveStatus" dc:"直播状态(0未开播,1直播中)"`
	TotalIncome                  float64    `json:"totalIncome" dc:"直播收益"`
	TotalGiftIncome              float64    `json:"totalGiftIncome" dc:"累计礼物收益"`
	TotalPaidDanmakuIncome       float64    `json:"totalPaidDanmakuIncome" dc:"累计付费弹幕收益"`
	TotalPrivateRoomTicketIncome float64    `json:"totalPrivateRoomTicketIncome" dc:"累计私密直播间门票收益"`
	TotalPrivateRoomWatchIncome  float64    `json:"totalPrivateRoomWatchIncome" dc:"累计私密房观看收益"`
	Ban                          bool       `json:"ban" dc:"是否封禁"`
	BanApplyTime                 *time.Time `json:"banApplyTime" dc:"封禁截止时间"`
	BanReason                    string     `json:"banReason" dc:"封禁原因"`
	CreatedAt                    *time.Time `json:"createdAt"`
	RegisteredAt                 *time.Time `json:"registeredAt" dc:"账号注册时间"`
}

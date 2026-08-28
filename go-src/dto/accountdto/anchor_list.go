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
	Key          string `json:"key" dc:"查询关键字(用户ID模糊/昵称/手机号/分享码)"`
	GuildId      uint64 `json:"guildId,string" dc:"工会ID(可选,大于0时按工会过滤)"`
	PlatformOnly bool   `json:"platformOnly" dc:"仅平台主播(guild_id=0)"`
	GuildOnly    bool   `json:"guildOnly" dc:"仅工会主播(guild_id>0)"`
	LiveStatus   *uint8 `json:"liveStatus" dc:"直播状态过滤(0未开播,1直播中,不传则不过滤)"`
}

// AnchorListItem 主播列表项(基于 user_infos)
type AnchorListItem struct {
	ID                           uint64     `json:"id,string"`
	Nickname                     string     `json:"nickname"`
	Phone                        string     `json:"phone"`
	Avatar                       string     `json:"avatar"`
	UserType                     uint8      `json:"userType" dc:"用户类型(1=普通主播,7=高级主播)"`
	GuildId                      uint64     `json:"guildId,string"`
	GuildName                    string     `json:"guildName" dc:"工会名称"`
	IP                           string     `json:"ip" dc:"登录IP"`
	RoomTitle                    string     `json:"roomTitle"`
	RoomCover                    string     `json:"roomCover" dc:"直播间封面URL"`
	RoomId                       uint64     `json:"roomId,string" dc:"直播间ID"`
	Category                     uint8      `json:"category" dc:"分类(1=hot,2=game,3=私密)"`
	PrivateInviteType            uint8      `json:"privateInviteType" dc:"私密邀请类型(1=接受所有人,3=拒绝所有人)"`
	Ticket                       float64    `json:"ticket" dc:"门票价格(钻石,私密直播间)"`
	Billing                      float64    `json:"billing" dc:"计费价格(每分钟钻石,私密直播间)"`
	LiveStatus                   uint8      `json:"liveStatus" dc:"直播状态(0未开播,1直播中)"`
	TotalIncome                  float64    `json:"totalIncome" dc:"未结算总直播收益"`
	TotalGiftIncome              float64    `json:"totalGiftIncome" dc:"未结算礼物收益"`
	TotalPaidDanmakuIncome       float64    `json:"totalPaidDanmakuIncome" dc:"未结算付费弹幕收益"`
	TotalPrivateRoomTicketIncome float64    `json:"totalPrivateRoomTicketIncome" dc:"未结算私密直播间门票收益"`
	TotalPrivateRoomWatchIncome  float64    `json:"totalPrivateRoomWatchIncome" dc:"未结算私密房观看收益"`
	TotalVideoCallIncome         float64    `json:"totalVideoCallIncome" dc:"未结算直播间视频通话收益"`
	TotalVideoCallTicketIncome   float64    `json:"totalVideoCallTicketIncome" dc:"未结算直播间视频通话门票收益"`
	TotalVideoCallBillingIncome  float64    `json:"totalVideoCallBillingIncome" dc:"未结算直播间视频通话计费收益"`
	Ban                          bool       `json:"ban" dc:"是否封禁"`
	BanApplyTime                 *time.Time `json:"banApplyTime" dc:"封禁截止时间"`
	BanReason                    string     `json:"banReason" dc:"封禁原因"`
	Status                       uint8      `json:"status" dc:"直播间状态(0-下架,1-上架)"`
	CreatedAt                    *time.Time `json:"createdAt"`
	RegisteredAt                 *time.Time `json:"registeredAt" dc:"账号注册时间"`
}

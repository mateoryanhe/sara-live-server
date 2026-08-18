package guilddto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dto/accountdto"
)

// GetGuildDetailReq CMS获取工会详情收益(基本信息由前端页面传入)
type GetGuildDetailReq struct {
	g.Meta  `path:"/getGuildDetail" method:"post" summary:"CMS获取工会详情收益" tags:"直播工会"`
	GuildId uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
}

// GetGuildIncomeArchivesReq CMS获取工会下架归档(直查DB)
type GetGuildIncomeArchivesReq struct {
	g.Meta  `path:"/getGuildIncomeArchives" method:"post" summary:"CMS获取工会下架归档" tags:"直播工会"`
	GuildId uint64 `json:"guildId" v:"required#工会ID不能为空" dc:"工会ID"`
}

// GuildDetailItem 工会基本信息
type GuildDetailItem struct {
	ID          uint64     `json:"id,string"`
	Name        string     `json:"name"`
	LeaderId    uint64     `json:"leaderId,string"`
	LeaderName  string     `json:"leaderName"`
	Description string     `json:"description"`
	Status      uint8      `json:"status"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

// GuildIncomeArchiveItem 工会下架未结算收益归档
type GuildIncomeArchiveItem struct {
	ID      uint64 `json:"id,string"`
	GuildId uint64 `json:"guildId,string"`
	accountdto.LiveRoomIncomeAmountsItem
	CreatedAt *time.Time `json:"createdAt"`
}

// GetGuildDetailRes CMS工会详情收益(缓存优先,否则直查DB)
type GetGuildDetailRes struct {
	IncomeUnsettled *accountdto.LiveRoomIncomeUnsettledItem `json:"incomeUnsettled"`
	IncomeSettled   *accountdto.LiveRoomIncomeSettledItem   `json:"incomeSettled"`
	IncomeTotal     *accountdto.LiveRoomIncomeTotalItem     `json:"incomeTotal"`
}

// GetGuildIncomeArchivesRes CMS工会下架归档(直查DB)
type GetGuildIncomeArchivesRes struct {
	List []*GuildIncomeArchiveItem `json:"list"`
}

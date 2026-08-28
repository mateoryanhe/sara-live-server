package accountdto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// GetAnchorDetailReq CMS获取主播详情
type GetAnchorDetailReq struct {
	g.Meta   `path:"/getAnchorDetail" method:"post" summary:"CMS获取主播详情" tags:"账号"`
	AnchorId uint64 `json:"anchorId" v:"required#主播ID不能为空" dc:"主播用户ID"`
}

// LiveRoomIncomeAmountsItem 直播间收益金额字段
type LiveRoomIncomeAmountsItem struct {
	TotalIncome                  float64 `json:"totalIncome"`
	TotalGiftIncome              float64 `json:"totalGiftIncome"`
	TotalPaidDanmakuIncome       float64 `json:"totalPaidDanmakuIncome"`
	TotalPrivateRoomTicketIncome float64 `json:"totalPrivateRoomTicketIncome"`
	TotalPrivateRoomWatchIncome  float64 `json:"totalPrivateRoomWatchIncome"`
	TotalVideoCallIncome         float64 `json:"totalVideoCallIncome"`
	TotalVideoCallTicketIncome   float64 `json:"totalVideoCallTicketIncome"`
	TotalVideoCallBillingIncome  float64 `json:"totalVideoCallBillingIncome"`
	TotalShortVideoIncome        float64 `json:"totalShortVideoIncome"`
	TotalLiveDuration            float64 `json:"totalLiveDuration"`
}

// LiveRoomIncomeUnsettledItem 未结算收益
type LiveRoomIncomeUnsettledItem struct {
	LiveRoomIncomeAmountsItem
	UpdatedAt *time.Time `json:"updatedAt"`
}

// LiveRoomIncomeSettledItem 已结算收益
type LiveRoomIncomeSettledItem struct {
	LiveRoomIncomeAmountsItem
	SettlementSalary      float64    `json:"settlementSalary"`
	SettlementShareAmount float64    `json:"settlementShareAmount"`
	UpdatedAt             *time.Time `json:"updatedAt"`
}

// LiveRoomIncomeTotalItem 生涯累计收益
type LiveRoomIncomeTotalItem struct {
	LiveRoomIncomeAmountsItem
	SettlementSalary      float64    `json:"settlementSalary"`
	SettlementShareAmount float64    `json:"settlementShareAmount"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

// AnchorLiveRoomDetailItem 直播间详情
type AnchorLiveRoomDetailItem struct {
	ID                uint64     `json:"id,string"`
	GuildId           uint64     `json:"guildId,string"`
	Title             string     `json:"title"`
	Cover             string     `json:"cover"`
	Notice            string     `json:"notice"`
	LiveRecordId      uint64     `json:"liveRecordId,string"`
	HeartTime         *time.Time `json:"heartTime"`
	Ban               bool       `json:"ban"`
	BanApplyTime      *time.Time `json:"banApplyTime"`
	BanReason         string     `json:"banReason"`
	Status            uint8      `json:"status"`
	LiveStatus        uint8      `json:"liveStatus"`
	Category          uint8      `json:"category"`
	PrivateInviteType uint8      `json:"privateInviteType"`
	Ticket            float64    `json:"ticket"`
	Billing           float64    `json:"billing"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
}

// LiveRoomIncomeArchiveItem 下架未结算收益归档
type LiveRoomIncomeArchiveItem struct {
	ID               uint64     `json:"id,string"`
	RoomId           uint64     `json:"roomId,string"`
	GuildId          uint64     `json:"guildId,string"`
	LiveRoomIncomeAmountsItem
	SettlementSalary float64    `json:"settlementSalary"`
	CreatedAt        *time.Time `json:"createdAt"`
}

// GetAnchorDetailRes CMS主播详情
type GetAnchorDetailRes struct {
	Anchor          *AnchorListItem              `json:"anchor"`
	LiveRoom        *AnchorLiveRoomDetailItem    `json:"liveRoom"`
	IncomeUnsettled *LiveRoomIncomeUnsettledItem `json:"incomeUnsettled"`
	IncomeSettled   *LiveRoomIncomeSettledItem   `json:"incomeSettled"`
	IncomeTotal     *LiveRoomIncomeTotalItem     `json:"incomeTotal"`
	IncomeArchives  []*LiveRoomIncomeArchiveItem `json:"incomeArchives"`
}

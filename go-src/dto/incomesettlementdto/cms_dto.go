package incomesettlementdto

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CMSAnchorIncomeSettlementLogListReq CMS分页查询主播结算流水
type CMSAnchorIncomeSettlementLogListReq struct {
	g.Meta `path:"/cmsAnchorIncomeSettlementLogList" method:"post" summary:"CMS查询主播结算流水" tags:"主播结算流水"`
	httpserver.CMSQueryReq
	RoomId    string   `json:"roomId"    dc:"直播间ID(可选,留空查全部,兼容旧版单选)"`
	AnchorIds []string `json:"anchorIds" dc:"主播ID列表(可选,多选)"`
	StartTime int64    `json:"startTime" dc:"创建时间起(秒, 0=不过滤)"`
	EndTime   int64    `json:"endTime"   dc:"创建时间止(秒, 0=不过滤)"`
}

// CMSGuildIncomeSettlementLogListReq CMS分页查询工会结算流水
type CMSGuildIncomeSettlementLogListReq struct {
	g.Meta `path:"/cmsGuildIncomeSettlementLogList" method:"post" summary:"CMS查询工会结算流水" tags:"工会结算流水"`
	httpserver.CMSQueryReq
	GuildId                  string `json:"guildId"   dc:"工会ID(可选,留空查全部)"`
	StartTime                int64  `json:"startTime" dc:"创建时间起(秒, 0=不过滤)"`
	EndTime                  int64  `json:"endTime"   dc:"创建时间止(秒, 0=不过滤)"`
	Status                   *uint8 `json:"status"    dc:"状态(可选,0未审核1审核通过2转账成功)"`
	OrderByReceivableUsdDesc bool   `json:"orderByReceivableUsdDesc" dc:"按可收金额USD降序"`
	IncludeTransferInfo      bool   `json:"includeTransferInfo" dc:"是否附带工会收款信息"`
}

// CMSBatchApproveGuildSettlementReq CMS批量审核工会结算流水
type CMSBatchApproveGuildSettlementReq struct {
	g.Meta `path:"/cmsBatchApproveGuildSettlement" method:"post" summary:"CMS批量审核工会结算流水" tags:"工会结算流水"`
	Ids    []string `json:"ids" v:"required#请选择结算流水"`
}

// CMSBatchApproveGuildSettlementRes CMS批量审核结果
type CMSBatchApproveGuildSettlementRes struct {
	SuccessCount int `json:"successCount"`
	FailCount    int `json:"failCount"`
}

// CMSBatchTransferGuildSettlementReq CMS批量转账(预留)
type CMSBatchTransferGuildSettlementReq struct {
	g.Meta `path:"/cmsBatchTransferGuildSettlement" method:"post" summary:"CMS批量转账工会结算(预留)" tags:"工会结算流水"`
	Ids    []string `json:"ids" v:"required#请选择结算流水"`
}

// CMSBatchTransferGuildSettlementRes CMS批量转账预留结果
type CMSBatchTransferGuildSettlementRes struct {
	Reserved bool   `json:"reserved"`
	Message  string `json:"message"`
}

// CMSIncomeSettlementLogItem CMS结算流水列表项(主播/工会共用收益快照字段)
type CMSIncomeSettlementLogItem struct {
	Id                           uint64     `json:"id,string"`
	RoomId                       uint64     `json:"roomId,string"`
	RoomNickname                 string     `json:"roomNickname"`
	RoomAvatar                   string     `json:"roomAvatar"`
	GuildId                      uint64     `json:"guildId,string"`
	GuildName                    string     `json:"guildName"`
	TotalIncome                  float64    `json:"totalIncome"`
	TotalGiftIncome              float64    `json:"totalGiftIncome"`
	TotalPaidDanmakuIncome       float64    `json:"totalPaidDanmakuIncome"`
	TotalPrivateRoomTicketIncome float64    `json:"totalPrivateRoomTicketIncome"`
	TotalPrivateRoomWatchIncome  float64    `json:"totalPrivateRoomWatchIncome"`
	TotalVideoCallIncome         float64    `json:"totalVideoCallIncome"`
	TotalVideoCallTicketIncome   float64    `json:"totalVideoCallTicketIncome"`
	TotalVideoCallBillingIncome  float64    `json:"totalVideoCallBillingIncome"`
	TotalShortVideoIncome        float64    `json:"totalShortVideoIncome"`
	TotalGameIncome              float64    `json:"totalGameIncome"`
	TotalLiveDuration            float64    `json:"totalLiveDuration"`
	SettlementSalary             float64    `json:"settlementSalary"`
	SettlementShareAmount        float64    `json:"settlementShareAmount"`
	SettlementShareAmountUsd     float64    `json:"settlementShareAmountUsd"`
	SettlementReceivableUsd      float64    `json:"settlementReceivableUsd"`
	AnchorSharePercent           float64    `json:"anchorSharePercent"`
	GuildSharePercent            float64    `json:"guildSharePercent"`
	Status                       uint8      `json:"status"`
	TransferAt                   *time.Time `json:"transferAt"`
	TransferCurrency             string     `json:"transferCurrency"`
	TransferPayeeName            string     `json:"transferPayeeName"`
	TransferBankName             string     `json:"transferBankName"`
	TransferAccountNo            string     `json:"transferAccountNo"`
	TransferBankCode             string     `json:"transferBankCode"`
	CreatedAt                    *time.Time `json:"createdAt"`
}

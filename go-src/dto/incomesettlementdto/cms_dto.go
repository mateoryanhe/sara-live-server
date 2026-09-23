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
	RoomId                   string   `json:"roomId"    dc:"直播间ID(可选,留空查全部,兼容旧版单选)"`
	AnchorIds                []string `json:"anchorIds" dc:"主播ID列表(可选,多选)"`
	StartTime                int64    `json:"startTime" dc:"创建时间起(秒, 0=不过滤)"`
	EndTime                  int64    `json:"endTime"   dc:"创建时间止(秒, 0=不过滤)"`
	Status                   *uint8   `json:"status" dc:"平台主播代付状态(可选,0审核中1审核通过2转账成功3代付中)"`
	DirectPayout             *bool    `json:"directPayout" dc:"是否仅查询平台主播直接代付单"`
	OrderByReceivableUsdDesc bool     `json:"orderByReceivableUsdDesc" dc:"按可收金额USD降序"`
	IncludeTransferInfo      bool     `json:"includeTransferInfo" dc:"是否附带平台主播收款信息"`
}

type CMSBatchApproveAnchorSettlementReq struct {
	g.Meta `path:"/cmsBatchApproveAnchorSettlement" method:"post" summary:"CMS批量审核平台主播结算" tags:"主播结算流水"`
	Ids    []string `json:"ids" v:"required#请选择结算流水"`
}

type CMSBatchApproveAnchorSettlementRes struct {
	SuccessCount int `json:"successCount"`
	FailCount    int `json:"failCount"`
}

type CMSBatchTransferAnchorSettlementReq struct {
	g.Meta `path:"/cmsBatchTransferAnchorSettlement" method:"post" summary:"CMS批量HaiPay代付平台主播结算" tags:"主播结算流水"`
	Ids    []string `json:"ids" v:"required#请选择结算流水"`
}

type CMSBatchTransferAnchorSettlementRes struct {
	SuccessCount int    `json:"successCount"`
	FailCount    int    `json:"failCount"`
	Message      string `json:"message"`
}

// CMSGuildIncomeSettlementLogListReq CMS分页查询工会结算流水
type CMSGuildIncomeSettlementLogListReq struct {
	g.Meta `path:"/cmsGuildIncomeSettlementLogList" method:"post" summary:"CMS查询工会结算流水" tags:"工会结算流水"`
	httpserver.CMSQueryReq
	GuildId                  string `json:"guildId"   dc:"工会ID(可选,留空查全部)"`
	GuildType                *uint8 `json:"guildType" v:"in:0,1#工会类型仅支持0普通或1币商" dc:"工会类型(可选,0普通工会,1币商工会)"`
	StartTime                int64  `json:"startTime" dc:"创建时间起(秒, 0=不过滤)"`
	EndTime                  int64  `json:"endTime"   dc:"创建时间止(秒, 0=不过滤)"`
	TransferStartTime        int64  `json:"transferStartTime" dc:"代付时间起(秒, 0=不过滤)"`
	TransferEndTime          int64  `json:"transferEndTime" dc:"代付时间止(秒, 0=不过滤)"`
	PayoutOnly               bool   `json:"payoutOnly" dc:"是否只查询已产生代付过程的结算单"`
	Status                   *uint8 `json:"status"    dc:"状态(可选,0审核中1审核通过2转账成功3代付中)"`
	OrderByReceivableUsdDesc bool   `json:"orderByReceivableUsdDesc" dc:"按可收金额USD降序"`
	IncludeDetail            bool   `json:"includeDetail" dc:"是否附带结算快照与代付过程"`
	IncludeTransfer          bool   `json:"includeTransfer" dc:"是否附带代付过程"`
	IncludeTransferInfo      bool   `json:"includeTransferInfo" dc:"是否附带工会收款信息"`
}

// CMSGuildIncomeSettlementLogDetailReq CMS查询工会结算主单完整明细。
type CMSGuildIncomeSettlementLogDetailReq struct {
	g.Meta `path:"/cmsGuildIncomeSettlementLogDetail" method:"post" summary:"CMS查询工会结算明细" tags:"工会结算流水"`
	Id     string `json:"id" v:"required#请选择结算流水" dc:"工会结算主单ID"`
}

type CMSGuildIncomeSettlementLogDetailRes struct {
	Item *CMSIncomeSettlementLogItem `json:"item"`
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

// CMSReopenGuildSettlementApprovalReq CMS将审核通过的工会结算退回审核中
type CMSReopenGuildSettlementApprovalReq struct {
	g.Meta `path:"/cmsReopenGuildSettlementApproval" method:"post" summary:"CMS将工会结算退回审核中" tags:"工会结算流水"`
	Id     string `json:"id" v:"required#请选择结算流水" dc:"工会结算流水ID"`
}

type CMSReopenGuildSettlementApprovalRes struct {
	Success bool `json:"success"`
}

// CMSCopyGuildSettlementPayoutReq CMS复制已成功的工会代付订单
type CMSCopyGuildSettlementPayoutReq struct {
	g.Meta `path:"/cmsCopyGuildSettlementPayout" method:"post" summary:"CMS复制工会代付订单" tags:"工会结算流水"`
	Id     string `json:"id" v:"required#请选择结算流水" dc:"已转账成功的工会结算流水ID"`
}

type CMSCopyGuildSettlementPayoutRes struct {
	Id string `json:"id" dc:"新工会结算流水ID，同时作为下一次HaiPay代付orderId"`
}

// CMSUpdateGuildSettlementReceivableUsdReq CMS修改审核中工会结算的可收金额
type CMSUpdateGuildSettlementReceivableUsdReq struct {
	g.Meta                  `path:"/cmsUpdateGuildSettlementReceivableUsd" method:"post" summary:"CMS修改工会结算可收金额" tags:"工会结算流水"`
	Id                      string  `json:"id" v:"required#请选择结算流水" dc:"工会结算流水ID"`
	SettlementReceivableUsd float64 `json:"settlementReceivableUsd" v:"required|min:0.0001#请输入可收金额|可收金额必须大于0" dc:"修改后的结算可收金额(USD)"`
}

type CMSUpdateGuildSettlementReceivableUsdRes struct {
	Success bool `json:"success"`
}

// CMSBatchTransferGuildSettlementReq CMS批量代付工会结算
type CMSBatchTransferGuildSettlementReq struct {
	g.Meta `path:"/cmsBatchTransferGuildSettlement" method:"post" summary:"CMS批量HaiPay代付工会结算" tags:"工会结算流水"`
	Ids    []string `json:"ids" v:"required#请选择结算流水"`
}

// CMSBatchTransferGuildSettlementRes CMS批量代付结果
type CMSBatchTransferGuildSettlementRes struct {
	SuccessCount int    `json:"successCount"`
	FailCount    int    `json:"failCount"`
	Message      string `json:"message"`
}

// CMSIncomeSettlementLogItem CMS结算流水列表项(主播/工会共用收益快照字段)
type CMSIncomeSettlementLogItem struct {
	Id                          uint64     `json:"id,string"`
	RoomId                      uint64     `json:"roomId,string"`
	RoomNickname                string     `json:"roomNickname"`
	RoomAvatar                  string     `json:"roomAvatar"`
	GuildId                     uint64     `json:"guildId,string"`
	GuildName                   string     `json:"guildName"`
	DirectPayout                bool       `json:"directPayout"`
	TotalIncome                 float64    `json:"totalIncome"`
	TotalSocialIncome           float64    `json:"totalSocialIncome"`
	TotalGiftIncome             float64    `json:"totalGiftIncome"`
	TotalPaidDanmakuIncome      float64    `json:"totalPaidDanmakuIncome"`
	TotalVideoCallIncome        float64    `json:"totalVideoCallIncome"`
	TotalVideoCallTicketIncome  float64    `json:"totalVideoCallTicketIncome"`
	TotalVideoCallBillingIncome float64    `json:"totalVideoCallBillingIncome"`
	TotalShortVideoIncome       float64    `json:"totalShortVideoIncome"`
	TotalGameIncome             float64    `json:"totalGameIncome"`
	TotalLiveDuration           float64    `json:"totalLiveDuration"`
	SettlementSalary            float64    `json:"settlementSalary"`
	SettlementShareAmount       float64    `json:"settlementShareAmount"`
	SettlementShareAmountUsd    float64    `json:"settlementShareAmountUsd"`
	SettlementReceivableUsd     float64    `json:"settlementReceivableUsd"`
	AnchorSharePercent          float64    `json:"anchorSharePercent"`
	GuildSharePercent           float64    `json:"guildSharePercent"`
	HasSalary                   bool       `json:"hasSalary"`
	AnchorSocialSharePercent    float64    `json:"anchorSocialSharePercent"`
	GuildSocialSharePercent     float64    `json:"guildSocialSharePercent"`
	AnchorGameSharePercent      float64    `json:"anchorGameSharePercent"`
	GuildGameSharePercent       float64    `json:"guildGameSharePercent"`
	AnchorSocialShareAmount     float64    `json:"anchorSocialShareAmount"`
	GuildSocialShareAmount      float64    `json:"guildSocialShareAmount"`
	AnchorGameShareAmountGold   float64    `json:"anchorGameShareAmountGold"`
	GuildGameShareAmountGold    float64    `json:"guildGameShareAmountGold"`
	SettlementRuleType          uint8      `json:"settlementRuleType"`
	GoldToDiamondRate           int        `json:"goldToDiamondRate"`
	UsdToGoldRate               int        `json:"usdToGoldRate"`
	GameShareAmountDiamond      float64    `json:"gameShareAmountDiamond"`
	TotalSettlementDiamond      float64    `json:"totalSettlementDiamond"`
	Status                      uint8      `json:"status"`
	TransferAt                  *time.Time `json:"transferAt"`
	TransferOrderId             string     `json:"transferOrderId"`
	TransferPlatformNo          string     `json:"transferPlatformNo"`
	TransferLocalAmount         float64    `json:"transferLocalAmount"`
	TransferFailMsg             string     `json:"transferFailMsg"`
	TransferCurrency            string     `json:"transferCurrency"`
	TransferPayeeName           string     `json:"transferPayeeName"`
	TransferBankName            string     `json:"transferBankName"`
	TransferAccountNo           string     `json:"transferAccountNo"`
	TransferBankCode            string     `json:"transferBankCode"`
	CreatedAt                   *time.Time `json:"createdAt"`
	UpdatedAt                   *time.Time `json:"updatedAt"`
}

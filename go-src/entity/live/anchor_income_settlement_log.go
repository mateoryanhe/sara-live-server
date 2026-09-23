package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbAnchorIncomeSettlementLog db.TbName = "anchor_income_settlement_logs"
)

const (
	AnchorIncomeSettlementLogRoomId                    db.TbCol = "room_id"
	AnchorIncomeSettlementLogSettlementSalary          db.TbCol = "settlement_salary"
	AnchorIncomeSettlementLogSettlementShareAmount     db.TbCol = "settlement_share_amount"
	AnchorIncomeSettlementLogSettlementShareAmountUsd  db.TbCol = "settlement_share_amount_usd"
	AnchorIncomeSettlementLogAnchorSharePercent        db.TbCol = "anchor_share_percent"
	AnchorIncomeSettlementLogHasSalary                 db.TbCol = "has_salary"
	AnchorIncomeSettlementLogAnchorSocialSharePercent  db.TbCol = "anchor_social_share_percent"
	AnchorIncomeSettlementLogGuildSocialSharePercent   db.TbCol = "guild_social_share_percent"
	AnchorIncomeSettlementLogAnchorGameSharePercent    db.TbCol = "anchor_game_share_percent"
	AnchorIncomeSettlementLogGuildGameSharePercent     db.TbCol = "guild_game_share_percent"
	AnchorIncomeSettlementLogAnchorSocialShareAmount   db.TbCol = "anchor_social_share_amount"
	AnchorIncomeSettlementLogGuildSocialShareAmount    db.TbCol = "guild_social_share_amount"
	AnchorIncomeSettlementLogAnchorGameShareAmountGold db.TbCol = "anchor_game_share_amount_gold"
	AnchorIncomeSettlementLogGuildGameShareAmountGold  db.TbCol = "guild_game_share_amount_gold"
	AnchorIncomeSettlementLogSettlementRuleType        db.TbCol = "settlement_rule_type"
	AnchorIncomeSettlementLogDirectPayout              db.TbCol = "direct_payout"
	AnchorIncomeSettlementLogStatus                    db.TbCol = "status"
	AnchorIncomeSettlementLogTransferAt                db.TbCol = "transfer_at"
	AnchorIncomeSettlementLogTransferOrderId           db.TbCol = "transfer_order_id"
	AnchorIncomeSettlementLogTransferPlatformNo        db.TbCol = "transfer_platform_no"
	AnchorIncomeSettlementLogTransferLocalAmount       db.TbCol = "transfer_local_amount"
	AnchorIncomeSettlementLogTransferCurrency          db.TbCol = "transfer_currency"
	AnchorIncomeSettlementLogTransferFailMsg           db.TbCol = "transfer_fail_msg"
	AnchorIncomeSettlementLogGoldToDiamondRate         db.TbCol = "gold_to_diamond_rate"
	AnchorIncomeSettlementLogUsdToGoldRate             db.TbCol = "usd_to_gold_rate"
	AnchorIncomeSettlementLogGameShareAmountDiamond    db.TbCol = "game_share_amount_diamond"
	AnchorIncomeSettlementLogTotalSettlementDiamond    db.TbCol = "total_settlement_diamond"
)

const (
	AnchorIncomeSettlementStatusPending      uint8 = 0 // 审核中
	AnchorIncomeSettlementStatusApproved     uint8 = 1 // 审核通过
	AnchorIncomeSettlementStatusTransferred  uint8 = 2 // 转账成功
	AnchorIncomeSettlementStatusTransferring uint8 = 3 // 已提交代付,等待回调

	AnchorIncomeSettlementRuleLegacy uint8 = 0
	AnchorIncomeSettlementRuleTiered uint8 = 1
)

// AnchorIncomeSettlementBreakdown 保留周结算时使用的原始单位与档位结果。
type AnchorIncomeSettlementBreakdown struct {
	SettlementRuleType        uint8
	HasSalary                 bool
	AnchorSocialSharePercent  float64
	GuildSocialSharePercent   float64
	AnchorGameSharePercent    float64
	GuildGameSharePercent     float64
	AnchorSocialShareAmount   float64
	GuildSocialShareAmount    float64
	AnchorGameShareAmountGold float64
	GuildGameShareAmountGold  float64
}

// AnchorIncomeSettlementLog 主播周结算成功日志(每次结算一条,历史留存)
type AnchorIncomeSettlementLog struct {
	migrate.OneModel
	CreatedAt time.Time `gorm:"index:idx_ais_room_created,priority:2" json:"-"`
	RoomId    uint64    `gorm:"index:idx_ais_room_created,priority:1;default:0;comment:直播间ID(==主播用户ID)" json:"roomId"`
	LiveRoomIncomeAmounts
	SettlementSalary          float64    `gorm:"type:decimal(16,4);default:0;comment:本次结算薪资" json:"settlementSalary"`
	SettlementShareAmount     float64    `gorm:"type:decimal(16,4);default:0;comment:本次结算分佣金额" json:"settlementShareAmount"`
	SettlementShareAmountUsd  float64    `gorm:"type:decimal(16,4);default:0;comment:本次结算分佣金额(USD)" json:"settlementShareAmountUsd"`
	AnchorSharePercent        float64    `gorm:"type:decimal(6,2);default:0;comment:本次结算主播分佣比例(%)" json:"anchorSharePercent"`
	HasSalary                 bool       `gorm:"default:0;comment:本结算周期是否按有底薪规则" json:"hasSalary"`
	AnchorSocialSharePercent  float64    `gorm:"type:decimal(6,2);default:0;comment:主播社交提成比(%)" json:"anchorSocialSharePercent"`
	GuildSocialSharePercent   float64    `gorm:"type:decimal(6,2);default:0;comment:工会社交提成比(%)" json:"guildSocialSharePercent"`
	AnchorGameSharePercent    float64    `gorm:"type:decimal(6,2);default:0;comment:主播游戏提成比(%)" json:"anchorGameSharePercent"`
	GuildGameSharePercent     float64    `gorm:"type:decimal(6,2);default:0;comment:工会游戏提成比(%)" json:"guildGameSharePercent"`
	AnchorSocialShareAmount   float64    `gorm:"type:decimal(20,4);default:0;comment:主播社交分佣(钻石)" json:"anchorSocialShareAmount"`
	GuildSocialShareAmount    float64    `gorm:"type:decimal(20,4);default:0;comment:工会社交分佣(钻石)" json:"guildSocialShareAmount"`
	AnchorGameShareAmountGold float64    `gorm:"type:decimal(20,4);default:0;comment:主播游戏分佣(金币)" json:"anchorGameShareAmountGold"`
	GuildGameShareAmountGold  float64    `gorm:"type:decimal(20,4);default:0;comment:工会游戏分佣(金币)" json:"guildGameShareAmountGold"`
	SettlementRuleType        uint8      `gorm:"default:0;comment:结算规则(0旧版 1分档结算)" json:"settlementRuleType"`
	DirectPayout              bool       `gorm:"index;default:0;comment:是否平台主播直接代付" json:"directPayout"`
	SettlementReceivableUsd   float64    `gorm:"type:decimal(16,4);default:0;comment:主播最终可收金额(USD)" json:"settlementReceivableUsd"`
	Status                    uint8      `gorm:"index;default:0;comment:代付状态(0审核中1审核通过2转账成功3代付中)" json:"status"`
	TransferAt                *time.Time `gorm:"comment:转账时间" json:"transferAt"`
	TransferOrderId           string     `gorm:"size:64;default:'';index;comment:HaiPay代付商户单号" json:"transferOrderId"`
	TransferPlatformNo        string     `gorm:"size:64;default:'';comment:HaiPay平台单号" json:"transferPlatformNo"`
	TransferLocalAmount       float64    `gorm:"type:decimal(18,4);default:0;comment:代付当地币金额" json:"transferLocalAmount"`
	TransferCurrency          string     `gorm:"size:16;default:'';comment:代付币种" json:"transferCurrency"`
	TransferFailMsg           string     `gorm:"size:512;default:'';comment:代付失败原因" json:"transferFailMsg"`
	GoldToDiamondRate         int        `gorm:"default:0;comment:代付换算时1金币兑换钻石数" json:"goldToDiamondRate"`
	UsdToGoldRate             int        `gorm:"default:0;comment:代付换算时1USD兑换金币数" json:"usdToGoldRate"`
	GameShareAmountDiamond    float64    `gorm:"type:decimal(20,4);default:0;comment:主播游戏分佣换算钻石" json:"gameShareAmountDiamond"`
	TotalSettlementDiamond    float64    `gorm:"type:decimal(20,4);default:0;comment:主播代付换算总钻石" json:"totalSettlementDiamond"`
}

// NewAnchorIncomeSettlementLog 新建一条主播结算日志并入库
func NewAnchorIncomeSettlementLog(roomId uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, anchorSharePercent float64) *AnchorIncomeSettlementLog {
	return newAnchorIncomeSettlementLogWithBreakdown(roomId, a, salary, shareAmount, shareAmountUsd, anchorSharePercent, nil, false)
}

// NewAnchorIncomeSettlementLogWithBreakdown 新建带社交/游戏分项的主播结算日志。
func NewAnchorIncomeSettlementLogWithBreakdown(roomId uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, anchorSharePercent float64, breakdown *AnchorIncomeSettlementBreakdown) *AnchorIncomeSettlementLog {
	return newAnchorIncomeSettlementLogWithBreakdown(roomId, a, salary, shareAmount, shareAmountUsd, anchorSharePercent, breakdown, false)
}

// NewPlatformAnchorIncomeSettlementLogWithBreakdown 新建平台主播可直接代付的周结算单。
func NewPlatformAnchorIncomeSettlementLogWithBreakdown(roomId uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, anchorSharePercent float64, breakdown *AnchorIncomeSettlementBreakdown) *AnchorIncomeSettlementLog {
	return newAnchorIncomeSettlementLogWithBreakdown(roomId, a, salary, shareAmount, 0, anchorSharePercent, breakdown, true)
}

func newAnchorIncomeSettlementLogWithBreakdown(roomId uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, anchorSharePercent float64, breakdown *AnchorIncomeSettlementBreakdown, directPayout bool) *AnchorIncomeSettlementLog {
	if a == nil {
		a = &LiveRoomIncomeAmounts{}
	}
	ret := &AnchorIncomeSettlementLog{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	ret.RoomId = roomId
	ret.LiveRoomIncomeAmounts = *a
	ret.SettlementSalary = salary
	ret.SettlementShareAmount = shareAmount
	ret.SettlementShareAmountUsd = shareAmountUsd
	ret.AnchorSharePercent = anchorSharePercent
	ret.DirectPayout = directPayout
	ret.Status = AnchorIncomeSettlementStatusPending
	if breakdown != nil {
		ret.SettlementRuleType = breakdown.SettlementRuleType
		ret.HasSalary = breakdown.HasSalary
		ret.AnchorSocialSharePercent = breakdown.AnchorSocialSharePercent
		ret.GuildSocialSharePercent = breakdown.GuildSocialSharePercent
		ret.AnchorGameSharePercent = breakdown.AnchorGameSharePercent
		ret.GuildGameSharePercent = breakdown.GuildGameSharePercent
		ret.AnchorSocialShareAmount = breakdown.AnchorSocialShareAmount
		ret.GuildSocialShareAmount = breakdown.GuildSocialShareAmount
		ret.AnchorGameShareAmountGold = breakdown.AnchorGameShareAmountGold
		ret.GuildGameShareAmountGold = breakdown.GuildGameShareAmountGold
	}

	syndb.AddData(TbAnchorIncomeSettlementLog, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbAnchorIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogRoomId, &syndb.ColData{IdVal: ret.ID, ColVal: roomId})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogDirectPayout, &syndb.ColData{IdVal: ret.ID, ColVal: directPayout})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogStatus, &syndb.ColData{IdVal: ret.ID, ColVal: AnchorIncomeSettlementStatusPending})
	writeAnchorIncomeSettlementLogAmounts(ret.ID, a, salary, shareAmount, shareAmountUsd, anchorSharePercent)
	writeAnchorIncomeSettlementBreakdown(ret.ID, breakdown)
	return ret
}

// SetPayoutConversion 保存主播代付实际使用的换算快照。
func (r *AnchorIncomeSettlementLog) SetPayoutConversion(goldToDiamondRate, usdToGoldRate int, gameShareDiamond, totalDiamond, receivableUsd float64) {
	if r == nil {
		return
	}
	r.GoldToDiamondRate = goldToDiamondRate
	r.UsdToGoldRate = usdToGoldRate
	r.GameShareAmountDiamond = gameShareDiamond
	r.TotalSettlementDiamond = totalDiamond
	r.SettlementReceivableUsd = receivableUsd
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGoldToDiamondRate, &syndb.ColData{IdVal: r.ID, ColVal: goldToDiamondRate})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogUsdToGoldRate, &syndb.ColData{IdVal: r.ID, ColVal: usdToGoldRate})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGameShareAmountDiamond, &syndb.ColData{IdVal: r.ID, ColVal: gameShareDiamond})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTotalSettlementDiamond, &syndb.ColData{IdVal: r.ID, ColVal: totalDiamond})
	syndb.AddData(TbAnchorIncomeSettlementLog, LiveRoomIncomeSettlementReceivableUsd, &syndb.ColData{IdVal: r.ID, ColVal: receivableUsd})
	syndb.AddData(TbAnchorIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
}

func (r *AnchorIncomeSettlementLog) SetStatus(v uint8) {
	if r == nil {
		return
	}
	r.Status = v
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogStatus, &syndb.ColData{IdVal: r.ID, ColVal: v})
	syndb.AddData(TbAnchorIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
}

func (r *AnchorIncomeSettlementLog) SetTransferAt(v *time.Time) {
	if r == nil {
		return
	}
	r.TransferAt = v
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferAt, &syndb.ColData{IdVal: r.ID, ColVal: v})
	syndb.AddData(TbAnchorIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
}

func (r *AnchorIncomeSettlementLog) SetTransferPayout(orderId, platformNo, currency string, localAmount float64) {
	if r == nil {
		return
	}
	r.TransferOrderId = orderId
	r.TransferPlatformNo = platformNo
	r.TransferCurrency = currency
	r.TransferLocalAmount = localAmount
	r.TransferFailMsg = ""
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferOrderId, &syndb.ColData{IdVal: r.ID, ColVal: orderId})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferPlatformNo, &syndb.ColData{IdVal: r.ID, ColVal: platformNo})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferCurrency, &syndb.ColData{IdVal: r.ID, ColVal: currency})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferLocalAmount, &syndb.ColData{IdVal: r.ID, ColVal: localAmount})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferFailMsg, &syndb.ColData{IdVal: r.ID, ColVal: ""})
	syndb.AddData(TbAnchorIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
}

func (r *AnchorIncomeSettlementLog) SetTransferFailMsg(msg string) {
	if r == nil {
		return
	}
	r.TransferFailMsg = msg
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferFailMsg, &syndb.ColData{IdVal: r.ID, ColVal: msg})
	syndb.AddData(TbAnchorIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
}

func initAnchorIncomeSettlementLog() {
	syndb.RegQuick(TbAnchorIncomeSettlementLog, db.CreatedAtName)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, db.UpdatedAtName)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogRoomId)
	regLiveRoomIncomeCols(TbAnchorIncomeSettlementLog)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogSettlementSalary)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogSettlementShareAmount)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogSettlementShareAmountUsd)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorSharePercent)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogHasSalary)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorSocialSharePercent)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGuildSocialSharePercent)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorGameSharePercent)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGuildGameSharePercent)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorSocialShareAmount)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGuildSocialShareAmount)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogAnchorGameShareAmountGold)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGuildGameShareAmountGold)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogSettlementRuleType)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogDirectPayout)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, LiveRoomIncomeSettlementReceivableUsd)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogStatus)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferAt)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferOrderId)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferPlatformNo)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferLocalAmount)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferCurrency)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTransferFailMsg)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGoldToDiamondRate)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogUsdToGoldRate)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogGameShareAmountDiamond)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogTotalSettlementDiamond)
	migrate.AutoMigrate(&AnchorIncomeSettlementLog{})
}

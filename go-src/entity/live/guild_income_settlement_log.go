package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbGuildIncomeSettlementLog      db.TbName = "guild_income_settlement_logs"
	TbGuildIncomeSettlementDetail   db.TbName = "guild_income_settlement_details"
	TbGuildIncomeSettlementTransfer db.TbName = "guild_income_settlement_transfers"
)

const (
	GuildIncomeSettlementLogGuildId                   db.TbCol = "guild_id"
	GuildIncomeSettlementLogSettlementSalary          db.TbCol = "settlement_salary"
	GuildIncomeSettlementLogSettlementShareAmount     db.TbCol = "settlement_share_amount"
	GuildIncomeSettlementLogSettlementShareAmountUsd  db.TbCol = "settlement_share_amount_usd"
	GuildIncomeSettlementLogGuildSharePercent         db.TbCol = "guild_share_percent"
	GuildIncomeSettlementLogStatus                    db.TbCol = "status"
	GuildIncomeSettlementLogTransferAt                db.TbCol = "transfer_at"
	GuildIncomeSettlementLogTransferOrderId           db.TbCol = "transfer_order_id"
	GuildIncomeSettlementLogTransferPlatformNo        db.TbCol = "transfer_platform_no"
	GuildIncomeSettlementLogTransferLocalAmount       db.TbCol = "transfer_local_amount"
	GuildIncomeSettlementLogTransferCurrency          db.TbCol = "transfer_currency"
	GuildIncomeSettlementLogTransferFailMsg           db.TbCol = "transfer_fail_msg"
	GuildIncomeSettlementLogSettlementRuleType        db.TbCol = "settlement_rule_type"
	GuildIncomeSettlementLogAnchorSocialShareAmount   db.TbCol = "anchor_social_share_amount"
	GuildIncomeSettlementLogGuildSocialShareAmount    db.TbCol = "guild_social_share_amount"
	GuildIncomeSettlementLogAnchorGameShareAmountGold db.TbCol = "anchor_game_share_amount_gold"
	GuildIncomeSettlementLogGuildGameShareAmountGold  db.TbCol = "guild_game_share_amount_gold"
	GuildIncomeSettlementLogGoldToDiamondRate         db.TbCol = "gold_to_diamond_rate"
	GuildIncomeSettlementLogUsdToGoldRate             db.TbCol = "usd_to_gold_rate"
	GuildIncomeSettlementLogGameShareAmountDiamond    db.TbCol = "game_share_amount_diamond"
	GuildIncomeSettlementLogTotalSettlementDiamond    db.TbCol = "total_settlement_diamond"
	GuildIncomeSettlementDetailSettlementId           db.TbCol = "settlement_id"
	GuildIncomeSettlementTransferSettlementId         db.TbCol = "settlement_id"
)

const (
	GuildIncomeSettlementStatusPending      uint8 = 0 // 审核中
	GuildIncomeSettlementStatusApproved     uint8 = 1 // 审核通过
	GuildIncomeSettlementStatusTransferred  uint8 = 2 // 转账成功
	GuildIncomeSettlementStatusTransferring uint8 = 3 // 已提交代付,等待回调

	GuildIncomeSettlementRuleLegacy uint8 = 0 // 旧逻辑/币商逻辑,结算时已得到USD
	GuildIncomeSettlementRuleTiered uint8 = 1 // 普通工会分项结算,代付时再换算USD
)

// GuildIncomeSettlementBreakdown 普通工会周结算的原始单位分项。
type GuildIncomeSettlementBreakdown struct {
	SettlementRuleType        uint8
	AnchorSocialShareAmount   float64
	GuildSocialShareAmount    float64
	AnchorGameShareAmountGold float64
	GuildGameShareAmountGold  float64
}

// GuildIncomeSettlementLog 工会结算主单。数据库主表只保存列表、审核和排序所需字段；
// 结算快照与代付过程分别落到一对一明细表，结构体继续作为业务聚合对象使用。
type GuildIncomeSettlementLog struct {
	migrate.OneModel
	CreatedAt               time.Time `gorm:"index:idx_gis_guild_status_created,priority:3" json:"-"`
	GuildId                 uint64    `gorm:"index:idx_gis_guild_status_created,priority:1;default:0;comment:工会ID" json:"guildId"`
	SettlementReceivableUsd float64   `gorm:"type:decimal(16,4);default:0;comment:结算可收金额(USD)" json:"settlementReceivableUsd"`
	Status                  uint8     `gorm:"index:idx_gis_guild_status_created,priority:2;default:0;comment:状态(0审核中1审核通过2转账成功3代付中)" json:"status"`

	LiveRoomIncomeAmounts     `gorm:"-"`
	SettlementSalary          float64    `gorm:"-" json:"settlementSalary"`
	SettlementShareAmount     float64    `gorm:"-" json:"settlementShareAmount"`
	SettlementShareAmountUsd  float64    `gorm:"-" json:"settlementShareAmountUsd"`
	GuildSharePercent         float64    `gorm:"-" json:"guildSharePercent"`
	SettlementRuleType        uint8      `gorm:"-" json:"settlementRuleType"`
	AnchorSocialShareAmount   float64    `gorm:"-" json:"anchorSocialShareAmount"`
	GuildSocialShareAmount    float64    `gorm:"-" json:"guildSocialShareAmount"`
	AnchorGameShareAmountGold float64    `gorm:"-" json:"anchorGameShareAmountGold"`
	GuildGameShareAmountGold  float64    `gorm:"-" json:"guildGameShareAmountGold"`
	GoldToDiamondRate         int        `gorm:"-" json:"goldToDiamondRate"`
	UsdToGoldRate             int        `gorm:"-" json:"usdToGoldRate"`
	GameShareAmountDiamond    float64    `gorm:"-" json:"gameShareAmountDiamond"`
	TotalSettlementDiamond    float64    `gorm:"-" json:"totalSettlementDiamond"`
	TransferAt                *time.Time `gorm:"-" json:"transferAt"`
	TransferOrderId           string     `gorm:"-" json:"transferOrderId"`
	TransferPlatformNo        string     `gorm:"-" json:"transferPlatformNo"`
	TransferLocalAmount       float64    `gorm:"-" json:"transferLocalAmount"`
	TransferCurrency          string     `gorm:"-" json:"transferCurrency"`
	TransferFailMsg           string     `gorm:"-" json:"transferFailMsg"`
}

// GuildIncomeSettlementDetail 工会结算计算快照，与主单一对一。
type GuildIncomeSettlementDetail struct {
	SettlementId uint64 `gorm:"column:settlement_id;primaryKey;autoIncrement:false;comment:工会结算主单ID" json:"settlementId,string"`
	LiveRoomIncomeAmounts
	SettlementSalary          float64 `gorm:"type:decimal(16,4);default:0;comment:结算薪资" json:"settlementSalary"`
	SettlementShareAmount     float64 `gorm:"type:decimal(16,4);default:0;comment:结算分佣金额" json:"settlementShareAmount"`
	SettlementShareAmountUsd  float64 `gorm:"type:decimal(16,4);default:0;comment:结算分佣金额(USD)" json:"settlementShareAmountUsd"`
	GuildSharePercent         float64 `gorm:"type:decimal(6,2);default:0;comment:本次结算工会分佣比例(%)" json:"guildSharePercent"`
	SettlementRuleType        uint8   `gorm:"default:0;comment:结算规则(0旧版或币商 1普通工会分项结算)" json:"settlementRuleType"`
	AnchorSocialShareAmount   float64 `gorm:"type:decimal(20,4);default:0;comment:主播社交分佣合计(钻石)" json:"anchorSocialShareAmount"`
	GuildSocialShareAmount    float64 `gorm:"type:decimal(20,4);default:0;comment:工会社交分佣合计(钻石)" json:"guildSocialShareAmount"`
	AnchorGameShareAmountGold float64 `gorm:"type:decimal(20,4);default:0;comment:主播游戏分佣合计(金币)" json:"anchorGameShareAmountGold"`
	GuildGameShareAmountGold  float64 `gorm:"type:decimal(20,4);default:0;comment:工会游戏分佣合计(金币)" json:"guildGameShareAmountGold"`
	GoldToDiamondRate         int     `gorm:"default:0;comment:代付换算时1金币兑换钻石数" json:"goldToDiamondRate"`
	UsdToGoldRate             int     `gorm:"default:0;comment:代付换算时1USD兑换金币数" json:"usdToGoldRate"`
	GameShareAmountDiamond    float64 `gorm:"type:decimal(20,4);default:0;comment:游戏分佣换算钻石" json:"gameShareAmountDiamond"`
	TotalSettlementDiamond    float64 `gorm:"type:decimal(20,4);default:0;comment:代付换算总钻石" json:"totalSettlementDiamond"`
}

func (GuildIncomeSettlementDetail) TableName() string {
	return string(TbGuildIncomeSettlementDetail)
}

// GuildIncomeSettlementTransfer 工会结算代付过程，与主单一对一。
type GuildIncomeSettlementTransfer struct {
	SettlementId        uint64     `gorm:"column:settlement_id;primaryKey;autoIncrement:false;comment:工会结算主单ID" json:"settlementId,string"`
	TransferAt          *time.Time `gorm:"comment:转账时间" json:"transferAt"`
	TransferOrderId     string     `gorm:"size:64;default:'';index;comment:HaiPay代付商户单号" json:"transferOrderId"`
	TransferPlatformNo  string     `gorm:"size:64;default:'';comment:HaiPay平台单号" json:"transferPlatformNo"`
	TransferLocalAmount float64    `gorm:"type:decimal(18,4);default:0;comment:代付当地币金额" json:"transferLocalAmount"`
	TransferCurrency    string     `gorm:"size:16;default:'';comment:代付币种" json:"transferCurrency"`
	TransferFailMsg     string     `gorm:"size:512;default:'';comment:代付失败原因" json:"transferFailMsg"`
}

func (GuildIncomeSettlementTransfer) TableName() string {
	return string(TbGuildIncomeSettlementTransfer)
}

// ApplyDetail 将持久化明细装配回业务聚合对象。
func (r *GuildIncomeSettlementLog) ApplyDetail(v *GuildIncomeSettlementDetail) {
	if r == nil || v == nil || v.SettlementId != r.ID {
		return
	}
	r.LiveRoomIncomeAmounts = v.LiveRoomIncomeAmounts
	r.SettlementSalary = v.SettlementSalary
	r.SettlementShareAmount = v.SettlementShareAmount
	r.SettlementShareAmountUsd = v.SettlementShareAmountUsd
	r.GuildSharePercent = v.GuildSharePercent
	r.SettlementRuleType = v.SettlementRuleType
	r.AnchorSocialShareAmount = v.AnchorSocialShareAmount
	r.GuildSocialShareAmount = v.GuildSocialShareAmount
	r.AnchorGameShareAmountGold = v.AnchorGameShareAmountGold
	r.GuildGameShareAmountGold = v.GuildGameShareAmountGold
	r.GoldToDiamondRate = v.GoldToDiamondRate
	r.UsdToGoldRate = v.UsdToGoldRate
	r.GameShareAmountDiamond = v.GameShareAmountDiamond
	r.TotalSettlementDiamond = v.TotalSettlementDiamond
}

// ApplyTransfer 将持久化代付过程装配回业务聚合对象。
func (r *GuildIncomeSettlementLog) ApplyTransfer(v *GuildIncomeSettlementTransfer) {
	if r == nil || v == nil || v.SettlementId != r.ID {
		return
	}
	r.TransferAt = v.TransferAt
	r.TransferOrderId = v.TransferOrderId
	r.TransferPlatformNo = v.TransferPlatformNo
	r.TransferLocalAmount = v.TransferLocalAmount
	r.TransferCurrency = v.TransferCurrency
	r.TransferFailMsg = v.TransferFailMsg
}

// NewGuildIncomeSettlementLog 新建一条工会结算日志并入库。
func NewGuildIncomeSettlementLog(guildId uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, receivableUsd, guildSharePercent float64) *GuildIncomeSettlementLog {
	return NewGuildIncomeSettlementLogWithBreakdown(guildId, a, salary, shareAmount, shareAmountUsd, receivableUsd, guildSharePercent, nil)
}

// NewGuildIncomeSettlementLogWithBreakdown 新建带原始单位分项的工会结算日志。
func NewGuildIncomeSettlementLogWithBreakdown(guildId uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, receivableUsd, guildSharePercent float64, breakdown *GuildIncomeSettlementBreakdown) *GuildIncomeSettlementLog {
	if a == nil {
		a = &LiveRoomIncomeAmounts{}
	}
	ret := &GuildIncomeSettlementLog{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	ret.GuildId = guildId
	ret.LiveRoomIncomeAmounts = *a
	ret.SettlementSalary = salary
	ret.SettlementShareAmount = shareAmount
	ret.SettlementShareAmountUsd = shareAmountUsd
	ret.SettlementReceivableUsd = receivableUsd
	ret.GuildSharePercent = guildSharePercent
	ret.Status = GuildIncomeSettlementStatusPending
	if breakdown != nil {
		ret.SettlementRuleType = breakdown.SettlementRuleType
		ret.AnchorSocialShareAmount = breakdown.AnchorSocialShareAmount
		ret.GuildSocialShareAmount = breakdown.GuildSocialShareAmount
		ret.AnchorGameShareAmountGold = breakdown.AnchorGameShareAmountGold
		ret.GuildGameShareAmountGold = breakdown.GuildGameShareAmountGold
	}

	syndb.AddData(TbGuildIncomeSettlementLog, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbGuildIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogGuildId, &syndb.ColData{IdVal: ret.ID, ColVal: guildId})
	syndb.AddData(TbGuildIncomeSettlementLog, LiveRoomIncomeSettlementReceivableUsd, &syndb.ColData{IdVal: ret.ID, ColVal: receivableUsd})
	writeGuildIncomeSettlementLogAmounts(ret.ID, a, salary, shareAmount, shareAmountUsd, guildSharePercent)
	writeGuildIncomeSettlementBreakdown(ret.ID, breakdown)
	ret.SetStatus(GuildIncomeSettlementStatusPending)
	return ret
}

func (r *GuildIncomeSettlementLog) touchMain() {
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(TbGuildIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
}

// SetPayoutConversion 保存本次代付实际使用的换算快照。
func (r *GuildIncomeSettlementLog) SetPayoutConversion(goldToDiamondRate, usdToGoldRate int, gameShareDiamond, totalDiamond, receivableUsd float64) {
	if r == nil {
		return
	}
	r.GoldToDiamondRate = goldToDiamondRate
	r.UsdToGoldRate = usdToGoldRate
	r.GameShareAmountDiamond = gameShareDiamond
	r.TotalSettlementDiamond = totalDiamond
	r.SettlementReceivableUsd = receivableUsd
	syndb.AddData(TbGuildIncomeSettlementDetail, GuildIncomeSettlementLogGoldToDiamondRate, &syndb.ColData{IdVal: r.ID, ColVal: goldToDiamondRate})
	syndb.AddData(TbGuildIncomeSettlementDetail, GuildIncomeSettlementLogUsdToGoldRate, &syndb.ColData{IdVal: r.ID, ColVal: usdToGoldRate})
	syndb.AddData(TbGuildIncomeSettlementDetail, GuildIncomeSettlementLogGameShareAmountDiamond, &syndb.ColData{IdVal: r.ID, ColVal: gameShareDiamond})
	syndb.AddData(TbGuildIncomeSettlementDetail, GuildIncomeSettlementLogTotalSettlementDiamond, &syndb.ColData{IdVal: r.ID, ColVal: totalDiamond})
	syndb.AddData(TbGuildIncomeSettlementLog, LiveRoomIncomeSettlementReceivableUsd, &syndb.ColData{IdVal: r.ID, ColVal: receivableUsd})
	r.touchMain()
}

func (r *GuildIncomeSettlementLog) SetStatus(v uint8) {
	if r == nil {
		return
	}
	r.Status = v
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogStatus, &syndb.ColData{IdVal: r.ID, ColVal: v})
	r.touchMain()
}

func (r *GuildIncomeSettlementLog) SetSettlementReceivableUsd(v float64) {
	if r == nil {
		return
	}
	r.SettlementReceivableUsd = v
	syndb.AddData(TbGuildIncomeSettlementLog, LiveRoomIncomeSettlementReceivableUsd, &syndb.ColData{IdVal: r.ID, ColVal: v})
	r.touchMain()
}

func (r *GuildIncomeSettlementLog) SetTransferAt(v *time.Time) {
	if r == nil {
		return
	}
	r.TransferAt = v
	syndb.AddData(TbGuildIncomeSettlementTransfer, GuildIncomeSettlementLogTransferAt, &syndb.ColData{IdVal: r.ID, ColVal: v})
	r.touchMain()
}

func (r *GuildIncomeSettlementLog) SetTransferPayout(orderId, platformNo, currency string, localAmount float64) {
	if r == nil {
		return
	}
	r.TransferOrderId = orderId
	r.TransferPlatformNo = platformNo
	r.TransferCurrency = currency
	r.TransferLocalAmount = localAmount
	r.TransferFailMsg = ""
	syndb.AddData(TbGuildIncomeSettlementTransfer, GuildIncomeSettlementLogTransferOrderId, &syndb.ColData{IdVal: r.ID, ColVal: orderId})
	syndb.AddData(TbGuildIncomeSettlementTransfer, GuildIncomeSettlementLogTransferPlatformNo, &syndb.ColData{IdVal: r.ID, ColVal: platformNo})
	syndb.AddData(TbGuildIncomeSettlementTransfer, GuildIncomeSettlementLogTransferCurrency, &syndb.ColData{IdVal: r.ID, ColVal: currency})
	syndb.AddData(TbGuildIncomeSettlementTransfer, GuildIncomeSettlementLogTransferLocalAmount, &syndb.ColData{IdVal: r.ID, ColVal: localAmount})
	syndb.AddData(TbGuildIncomeSettlementTransfer, GuildIncomeSettlementLogTransferFailMsg, &syndb.ColData{IdVal: r.ID, ColVal: ""})
	r.touchMain()
}

func (r *GuildIncomeSettlementLog) SetTransferFailMsg(msg string) {
	if r == nil {
		return
	}
	r.TransferFailMsg = msg
	syndb.AddData(TbGuildIncomeSettlementTransfer, GuildIncomeSettlementLogTransferFailMsg, &syndb.ColData{IdVal: r.ID, ColVal: msg})
	r.touchMain()
}

func regGuildIncomeSettlementDetailColumn(column db.TbCol) {
	syndb.RegWithIDName(TbGuildIncomeSettlementDetail, column, GuildIncomeSettlementDetailSettlementId)
}

func regGuildIncomeSettlementTransferColumn(column db.TbCol) {
	syndb.RegWithIDName(TbGuildIncomeSettlementTransfer, column, GuildIncomeSettlementTransferSettlementId)
}

func initGuildIncomeSettlementLog() {
	for _, column := range []db.TbCol{
		db.CreatedAtName,
		db.UpdatedAtName,
		GuildIncomeSettlementLogGuildId,
		LiveRoomIncomeSettlementReceivableUsd,
		GuildIncomeSettlementLogStatus,
	} {
		syndb.RegQuick(TbGuildIncomeSettlementLog, column)
	}
	for _, column := range []db.TbCol{
		LiveRoomIncomeTotalIncome,
		LiveRoomIncomeTotalSocialIncome,
		LiveRoomIncomeTotalGiftIncome,
		LiveRoomIncomeTotalPaidDanmakuIncome,
		LiveRoomIncomeTotalVideoCallIncome,
		LiveRoomIncomeTotalVideoCallTicketIncome,
		LiveRoomIncomeTotalVideoCallBillingIncome,
		LiveRoomIncomeTotalShortVideoIncome,
		LiveRoomIncomeTotalGameIncome,
		LiveRoomIncomeTotalLiveDuration,
		GuildIncomeSettlementLogSettlementSalary,
		GuildIncomeSettlementLogSettlementShareAmount,
		GuildIncomeSettlementLogSettlementShareAmountUsd,
		GuildIncomeSettlementLogGuildSharePercent,
		GuildIncomeSettlementLogSettlementRuleType,
		GuildIncomeSettlementLogAnchorSocialShareAmount,
		GuildIncomeSettlementLogGuildSocialShareAmount,
		GuildIncomeSettlementLogAnchorGameShareAmountGold,
		GuildIncomeSettlementLogGuildGameShareAmountGold,
		GuildIncomeSettlementLogGoldToDiamondRate,
		GuildIncomeSettlementLogUsdToGoldRate,
		GuildIncomeSettlementLogGameShareAmountDiamond,
		GuildIncomeSettlementLogTotalSettlementDiamond,
	} {
		regGuildIncomeSettlementDetailColumn(column)
	}
	for _, column := range []db.TbCol{
		GuildIncomeSettlementLogTransferAt,
		GuildIncomeSettlementLogTransferOrderId,
		GuildIncomeSettlementLogTransferPlatformNo,
		GuildIncomeSettlementLogTransferLocalAmount,
		GuildIncomeSettlementLogTransferCurrency,
		GuildIncomeSettlementLogTransferFailMsg,
	} {
		regGuildIncomeSettlementTransferColumn(column)
	}
	migrate.AutoMigrate(
		&GuildIncomeSettlementLog{},
		&GuildIncomeSettlementDetail{},
		&GuildIncomeSettlementTransfer{},
	)
}

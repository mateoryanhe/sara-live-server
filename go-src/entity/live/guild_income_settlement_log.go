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
	GuildIncomeSettlementLogSettlementSalary         db.TbCol = "settlement_salary"
	GuildIncomeSettlementLogGuildSharePercent        db.TbCol = "guild_share_percent"
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
	GuildIncomeSettlementLogAnchorPayoutTotalDiamond  db.TbCol = "anchor_payout_total_diamond"
	GuildIncomeSettlementLogGuildPayoutTotalDiamond   db.TbCol = "guild_payout_total_diamond"
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
// 普通工会与币商工会使用不同的物理表，由 Storage 区分。
type GuildIncomeSettlementLog struct {
	migrate.OneModel
	Storage                 GuildIncomeSettlementStorage `gorm:"-" json:"-"`
	CreatedAt               time.Time                    `gorm:"index:idx_gis_guild_status_created,priority:3" json:"-"`
	GuildId                 uint64                       `gorm:"index:idx_gis_guild_status_created,priority:1;default:0;comment:工会ID" json:"guildId"`
	SettlementReceivableUsd float64   `gorm:"type:decimal(16,4);default:0;comment:结算可收金额(USD)" json:"settlementReceivableUsd"`
	Status                  uint8     `gorm:"index:idx_gis_guild_status_created,priority:2;default:0;comment:状态(0审核中1审核通过2转账成功3代付中)" json:"status"`

	LiveRoomIncomeAmounts     `gorm:"-"`
	SettlementSalary  float64 `gorm:"-" json:"settlementSalary"`
	GuildSharePercent float64 `gorm:"-" json:"guildSharePercent"`
	SettlementRuleType        uint8      `gorm:"-" json:"settlementRuleType"`
	AnchorSocialShareAmount   float64    `gorm:"-" json:"anchorSocialShareAmount"`
	GuildSocialShareAmount    float64    `gorm:"-" json:"guildSocialShareAmount"`
	AnchorGameShareAmountGold float64    `gorm:"-" json:"anchorGameShareAmountGold"`
	GuildGameShareAmountGold  float64    `gorm:"-" json:"guildGameShareAmountGold"`
	GoldToDiamondRate         int        `gorm:"-" json:"goldToDiamondRate"`
	UsdToGoldRate             int        `gorm:"-" json:"usdToGoldRate"`
	GameShareAmountDiamond     float64    `gorm:"-" json:"gameShareAmountDiamond"`
	TotalSettlementDiamond     float64    `gorm:"-" json:"totalSettlementDiamond"`
	AnchorPayoutTotalDiamond   float64    `gorm:"-" json:"anchorPayoutTotalDiamond"`
	GuildPayoutTotalDiamond    float64    `gorm:"-" json:"guildPayoutTotalDiamond"`
	TransferAt                 *time.Time `gorm:"-" json:"transferAt"`
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
	SettlementSalary  float64 `gorm:"type:decimal(16,4);default:0;comment:结算薪资(钻石)" json:"settlementSalary"`
	GuildSharePercent float64 `gorm:"type:decimal(6,2);default:0;comment:币商工会结算比例(%,普通工会为0)" json:"guildSharePercent"`
	SettlementRuleType        uint8   `gorm:"default:0;comment:结算规则(0旧版或币商 1普通工会分项结算)" json:"settlementRuleType"`
	AnchorSocialShareAmount   float64 `gorm:"type:decimal(20,4);default:0;comment:主播社交分佣合计(钻石)" json:"anchorSocialShareAmount"`
	GuildSocialShareAmount    float64 `gorm:"type:decimal(20,4);default:0;comment:工会社交分佣合计(钻石)" json:"guildSocialShareAmount"`
	AnchorGameShareAmountGold float64 `gorm:"type:decimal(20,4);default:0;comment:主播游戏分佣合计(金币)" json:"anchorGameShareAmountGold"`
	GuildGameShareAmountGold  float64 `gorm:"type:decimal(20,4);default:0;comment:工会游戏分佣合计(金币)" json:"guildGameShareAmountGold"`
	GoldToDiamondRate         int     `gorm:"default:0;comment:代付换算时1金币兑换钻石数" json:"goldToDiamondRate"`
	UsdToGoldRate             int     `gorm:"default:0;comment:代付换算时1USD兑换金币数" json:"usdToGoldRate"`
	GameShareAmountDiamond     float64 `gorm:"type:decimal(20,4);default:0;comment:游戏分佣换算钻石" json:"gameShareAmountDiamond"`
	TotalSettlementDiamond     float64 `gorm:"type:decimal(20,4);default:0;comment:代付换算总钻石" json:"totalSettlementDiamond"`
	AnchorPayoutTotalDiamond   float64 `gorm:"type:decimal(20,4);default:0;comment:旗下主播代付总额(钻石,分项结算快照)" json:"anchorPayoutTotalDiamond"`
	GuildPayoutTotalDiamond    float64 `gorm:"type:decimal(20,4);default:0;comment:工会自身代付总额(钻石,分项结算快照)" json:"guildPayoutTotalDiamond"`
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

func (GuildIncomeSettlementLog) TableName() string {
	return string(TbGuildIncomeSettlementLog)
}

type coinMerchantGuildIncomeSettlementLog struct {
	GuildIncomeSettlementLog
}

func (coinMerchantGuildIncomeSettlementLog) TableName() string {
	return string(TbCoinMerchantGuildIncomeSettlementLog)
}

type coinMerchantGuildIncomeSettlementDetail struct {
	GuildIncomeSettlementDetail
}

func (coinMerchantGuildIncomeSettlementDetail) TableName() string {
	return string(TbCoinMerchantGuildIncomeSettlementDetail)
}

type coinMerchantGuildIncomeSettlementTransfer struct {
	GuildIncomeSettlementTransfer
}

func (coinMerchantGuildIncomeSettlementTransfer) TableName() string {
	return string(TbCoinMerchantGuildIncomeSettlementTransfer)
}

// ApplyDetail 将持久化明细装配回业务聚合对象。
func (r *GuildIncomeSettlementLog) ApplyDetail(v *GuildIncomeSettlementDetail) {
	if r == nil || v == nil || v.SettlementId != r.ID {
		return
	}
	r.LiveRoomIncomeAmounts = v.LiveRoomIncomeAmounts
	r.SettlementSalary = v.SettlementSalary
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
	r.AnchorPayoutTotalDiamond = v.AnchorPayoutTotalDiamond
	r.GuildPayoutTotalDiamond = v.GuildPayoutTotalDiamond
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

// NewGuildIncomeSettlementLogWithBreakdown 新建普通工会结算日志（分项/阶梯规则）。
func NewGuildIncomeSettlementLogWithBreakdown(guildId uint64, a *LiveRoomIncomeAmounts, salary, receivableUsd, guildSharePercent float64, breakdown *GuildIncomeSettlementBreakdown) *GuildIncomeSettlementLog {
	return newGuildIncomeSettlementLogWithBreakdown(GuildIncomeSettlementStorageNormal, guildId, a, salary, receivableUsd, guildSharePercent, breakdown)
}

// NewCoinMerchantGuildIncomeSettlementLogWithBreakdown 新建币商工会结算日志（独立物理表）。
func NewCoinMerchantGuildIncomeSettlementLogWithBreakdown(guildId uint64, a *LiveRoomIncomeAmounts, salary, receivableUsd, guildSharePercent float64, breakdown *GuildIncomeSettlementBreakdown) *GuildIncomeSettlementLog {
	return newGuildIncomeSettlementLogWithBreakdown(GuildIncomeSettlementStorageCoinMerchant, guildId, a, salary, receivableUsd, guildSharePercent, breakdown)
}

func newGuildIncomeSettlementLogWithBreakdown(storage GuildIncomeSettlementStorage, guildId uint64, a *LiveRoomIncomeAmounts, salary, receivableUsd, guildSharePercent float64, breakdown *GuildIncomeSettlementBreakdown) *GuildIncomeSettlementLog {
	if a == nil {
		a = &LiveRoomIncomeAmounts{}
	}
	ret := &GuildIncomeSettlementLog{Storage: storage}
	tables := storage.tables()
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	ret.GuildId = guildId
	ret.LiveRoomIncomeAmounts = *a
	ret.SettlementSalary = salary
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

	syndb.AddData(tables.Log, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(tables.Log, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(tables.Log, GuildIncomeSettlementLogGuildId, &syndb.ColData{IdVal: ret.ID, ColVal: guildId})
	syndb.AddData(tables.Log, LiveRoomIncomeSettlementReceivableUsd, &syndb.ColData{IdVal: ret.ID, ColVal: receivableUsd})
	writeGuildIncomeSettlementLogAmounts(tables.Detail, ret.ID, a, salary, guildSharePercent)
	writeGuildIncomeSettlementBreakdown(tables.Detail, ret.ID, breakdown)
	ret.SetStatus(GuildIncomeSettlementStatusPending)
	return ret
}

func (r *GuildIncomeSettlementLog) touchMain() {
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(r.logTable(), db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
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
	syndb.AddData(r.detailTable(), GuildIncomeSettlementLogGoldToDiamondRate, &syndb.ColData{IdVal: r.ID, ColVal: goldToDiamondRate})
	syndb.AddData(r.detailTable(), GuildIncomeSettlementLogUsdToGoldRate, &syndb.ColData{IdVal: r.ID, ColVal: usdToGoldRate})
	syndb.AddData(r.detailTable(), GuildIncomeSettlementLogGameShareAmountDiamond, &syndb.ColData{IdVal: r.ID, ColVal: gameShareDiamond})
	syndb.AddData(r.detailTable(), GuildIncomeSettlementLogTotalSettlementDiamond, &syndb.ColData{IdVal: r.ID, ColVal: totalDiamond})
	if r.SettlementRuleType == GuildIncomeSettlementRuleTiered {
		anchorTotal, guildTotal := CalcTieredGuildPayoutDiamondSplit(
			r.SettlementSalary,
			r.AnchorSocialShareAmount,
			r.GuildSocialShareAmount,
			r.AnchorGameShareAmountGold,
			r.GuildGameShareAmountGold,
			goldToDiamondRate,
		)
		r.AnchorPayoutTotalDiamond = anchorTotal
		r.GuildPayoutTotalDiamond = guildTotal
		syndb.AddData(r.detailTable(), GuildIncomeSettlementLogAnchorPayoutTotalDiamond, &syndb.ColData{IdVal: r.ID, ColVal: anchorTotal})
		syndb.AddData(r.detailTable(), GuildIncomeSettlementLogGuildPayoutTotalDiamond, &syndb.ColData{IdVal: r.ID, ColVal: guildTotal})
	}
	syndb.AddData(r.logTable(), LiveRoomIncomeSettlementReceivableUsd, &syndb.ColData{IdVal: r.ID, ColVal: receivableUsd})
	r.touchMain()
}

// TieredGuildConversionSnapshotPersisted 普通工会分项结算的代付换算快照是否已完整落库（查询侧只读这些字段，不再重算）。
func (r *GuildIncomeSettlementLog) TieredGuildConversionSnapshotPersisted() bool {
	if r == nil || r.SettlementRuleType != GuildIncomeSettlementRuleTiered {
		return true
	}
	return r.SettlementReceivableUsd > 0 &&
		r.TotalSettlementDiamond > 0 &&
		r.GoldToDiamondRate > 0 &&
		r.UsdToGoldRate > 0
}

// PersistTieredPayoutDiamondSplit 仅补写主播/工会代付钻石拆分列（使用已落库的分项与汇率，供历史数据或写入路径一次性补齐）。
func (r *GuildIncomeSettlementLog) PersistTieredPayoutDiamondSplit() {
	if r == nil || r.ID == 0 || r.SettlementRuleType != GuildIncomeSettlementRuleTiered {
		return
	}
	if r.GoldToDiamondRate <= 0 {
		return
	}
	anchorTotal, guildTotal := CalcTieredGuildPayoutDiamondSplit(
		r.SettlementSalary,
		r.AnchorSocialShareAmount,
		r.GuildSocialShareAmount,
		r.AnchorGameShareAmountGold,
		r.GuildGameShareAmountGold,
		r.GoldToDiamondRate,
	)
	r.AnchorPayoutTotalDiamond = anchorTotal
	r.GuildPayoutTotalDiamond = guildTotal
	syndb.AddData(r.detailTable(), GuildIncomeSettlementLogAnchorPayoutTotalDiamond, &syndb.ColData{IdVal: r.ID, ColVal: anchorTotal})
	syndb.AddData(r.detailTable(), GuildIncomeSettlementLogGuildPayoutTotalDiamond, &syndb.ColData{IdVal: r.ID, ColVal: guildTotal})
	r.touchMain()
}

func (r *GuildIncomeSettlementLog) SetStatus(v uint8) {
	if r == nil {
		return
	}
	r.Status = v
	syndb.AddData(r.logTable(), GuildIncomeSettlementLogStatus, &syndb.ColData{IdVal: r.ID, ColVal: v})
	r.touchMain()
}

func (r *GuildIncomeSettlementLog) SetSettlementReceivableUsd(v float64) {
	if r == nil {
		return
	}
	r.SettlementReceivableUsd = v
	syndb.AddData(r.logTable(), LiveRoomIncomeSettlementReceivableUsd, &syndb.ColData{IdVal: r.ID, ColVal: v})
	r.touchMain()
}

func (r *GuildIncomeSettlementLog) SetTransferAt(v *time.Time) {
	if r == nil {
		return
	}
	r.TransferAt = v
	syndb.AddData(r.transferTable(), GuildIncomeSettlementLogTransferAt, &syndb.ColData{IdVal: r.ID, ColVal: v})
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
	syndb.AddData(r.transferTable(), GuildIncomeSettlementLogTransferOrderId, &syndb.ColData{IdVal: r.ID, ColVal: orderId})
	syndb.AddData(r.transferTable(), GuildIncomeSettlementLogTransferPlatformNo, &syndb.ColData{IdVal: r.ID, ColVal: platformNo})
	syndb.AddData(r.transferTable(), GuildIncomeSettlementLogTransferCurrency, &syndb.ColData{IdVal: r.ID, ColVal: currency})
	syndb.AddData(r.transferTable(), GuildIncomeSettlementLogTransferLocalAmount, &syndb.ColData{IdVal: r.ID, ColVal: localAmount})
	syndb.AddData(r.transferTable(), GuildIncomeSettlementLogTransferFailMsg, &syndb.ColData{IdVal: r.ID, ColVal: ""})
	r.touchMain()
}

func (r *GuildIncomeSettlementLog) SetTransferFailMsg(msg string) {
	if r == nil {
		return
	}
	r.TransferFailMsg = msg
	syndb.AddData(r.transferTable(), GuildIncomeSettlementLogTransferFailMsg, &syndb.ColData{IdVal: r.ID, ColVal: msg})
	r.touchMain()
}

func regGuildIncomeSettlementDetailColumn(tb db.TbName, column db.TbCol) {
	syndb.RegWithIDName(tb, column, GuildIncomeSettlementDetailSettlementId)
}

func regGuildIncomeSettlementTransferColumn(tb db.TbName, column db.TbCol) {
	syndb.RegWithIDName(tb, column, GuildIncomeSettlementTransferSettlementId)
}

func initGuildIncomeSettlementLogSyndb(storage GuildIncomeSettlementStorage) {
	tables := storage.tables()
	for _, column := range []db.TbCol{
		db.CreatedAtName,
		db.UpdatedAtName,
		GuildIncomeSettlementLogGuildId,
		LiveRoomIncomeSettlementReceivableUsd,
		GuildIncomeSettlementLogStatus,
	} {
		syndb.RegQuick(tables.Log, column)
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
		GuildIncomeSettlementLogAnchorPayoutTotalDiamond,
		GuildIncomeSettlementLogGuildPayoutTotalDiamond,
	} {
		regGuildIncomeSettlementDetailColumn(tables.Detail, column)
	}
	for _, column := range []db.TbCol{
		GuildIncomeSettlementLogTransferAt,
		GuildIncomeSettlementLogTransferOrderId,
		GuildIncomeSettlementLogTransferPlatformNo,
		GuildIncomeSettlementLogTransferLocalAmount,
		GuildIncomeSettlementLogTransferCurrency,
		GuildIncomeSettlementLogTransferFailMsg,
	} {
		regGuildIncomeSettlementTransferColumn(tables.Transfer, column)
	}
}

func initGuildIncomeSettlementLog() {
	initGuildIncomeSettlementLogSyndb(GuildIncomeSettlementStorageNormal)
	initGuildIncomeSettlementLogSyndb(GuildIncomeSettlementStorageCoinMerchant)
	migrate.AutoMigrate(
		&GuildIncomeSettlementLog{},
		&GuildIncomeSettlementDetail{},
		&GuildIncomeSettlementTransfer{},
		&coinMerchantGuildIncomeSettlementLog{},
		&coinMerchantGuildIncomeSettlementDetail{},
		&coinMerchantGuildIncomeSettlementTransfer{},
	)
}

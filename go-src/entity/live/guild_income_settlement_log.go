package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbGuildIncomeSettlementLog db.TbName = "guild_income_settlement_logs"
)

const (
	GuildIncomeSettlementLogGuildId               db.TbCol = "guild_id"
	GuildIncomeSettlementLogSettlementSalary      db.TbCol = "settlement_salary"
	GuildIncomeSettlementLogSettlementShareAmount    db.TbCol = "settlement_share_amount"
	GuildIncomeSettlementLogSettlementShareAmountUsd db.TbCol = "settlement_share_amount_usd"
	GuildIncomeSettlementLogGuildSharePercent        db.TbCol = "guild_share_percent"
)

// GuildIncomeSettlementLog 工会周结算成功日志(每次结算一条,历史留存)
type GuildIncomeSettlementLog struct {
	migrate.OneModel
	GuildId uint64 `gorm:"index;default:0;comment:工会ID" json:"guildId"`
	LiveRoomIncomeAmounts
	SettlementSalary          float64 `gorm:"type:decimal(16,4);default:0;comment:结算薪资" json:"settlementSalary"`
	SettlementShareAmount     float64 `gorm:"type:decimal(16,4);default:0;comment:结算分佣金额" json:"settlementShareAmount"`
	SettlementShareAmountUsd  float64 `gorm:"type:decimal(16,4);default:0;comment:结算分佣金额(USD)" json:"settlementShareAmountUsd"`
	SettlementReceivableUsd   float64 `gorm:"type:decimal(16,4);default:0;comment:结算可收金额(USD)=流水分佣+开播薪资" json:"settlementReceivableUsd"`
	GuildSharePercent         float64 `gorm:"type:decimal(6,2);default:0;comment:本次结算工会分佣比例(%)" json:"guildSharePercent"`
}

// NewGuildIncomeSettlementLog 新建一条工会结算日志并入库
func NewGuildIncomeSettlementLog(guildId uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, receivableUsd, guildSharePercent float64) *GuildIncomeSettlementLog {
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

	syndb.AddData(TbGuildIncomeSettlementLog, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbGuildIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogGuildId, &syndb.ColData{IdVal: ret.ID, ColVal: guildId})
	writeGuildIncomeSettlementLogAmounts(ret.ID, a, salary, shareAmount, shareAmountUsd, receivableUsd, guildSharePercent)
	return ret
}

func initGuildIncomeSettlementLog() {
	syndb.RegQuick(TbGuildIncomeSettlementLog, db.CreatedAtName)
	syndb.RegQuick(TbGuildIncomeSettlementLog, db.UpdatedAtName)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogGuildId)
	regLiveRoomIncomeCols(TbGuildIncomeSettlementLog)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogSettlementSalary)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogSettlementShareAmount)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogSettlementShareAmountUsd)
	syndb.RegQuick(TbGuildIncomeSettlementLog, LiveRoomIncomeSettlementReceivableUsd)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogGuildSharePercent)
	migrate.AutoMigrate(&GuildIncomeSettlementLog{})
}

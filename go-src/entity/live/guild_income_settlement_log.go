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
	GuildIncomeSettlementLogSettlementShareAmount db.TbCol = "settlement_share_amount"
	GuildIncomeSettlementLogGuildSharePercent     db.TbCol = "guild_share_percent"
)

// GuildIncomeSettlementLog 工会周结算成功日志(每次结算一条,历史留存)
type GuildIncomeSettlementLog struct {
	migrate.OneModel
	GuildId uint64 `gorm:"index;default:0;comment:工会ID" json:"guildId"`
	LiveRoomIncomeAmounts
	SettlementSalary      float64 `gorm:"type:decimal(16,4);default:0;comment:本次结算薪资" json:"settlementSalary"`
	SettlementShareAmount float64 `gorm:"type:decimal(16,4);default:0;comment:本次结算分佣金额" json:"settlementShareAmount"`
	GuildSharePercent     float64 `gorm:"type:decimal(6,2);default:0;comment:本次结算工会分佣比例(%)" json:"guildSharePercent"`
}

// NewGuildIncomeSettlementLog 新建一条工会结算日志并入库
func NewGuildIncomeSettlementLog(guildId uint64, a *LiveRoomIncomeAmounts, shareAmount, guildSharePercent float64) *GuildIncomeSettlementLog {
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
	ret.SettlementShareAmount = shareAmount
	ret.GuildSharePercent = guildSharePercent

	syndb.AddData(TbGuildIncomeSettlementLog, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbGuildIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogGuildId, &syndb.ColData{IdVal: ret.ID, ColVal: guildId})
	writeGuildIncomeSettlementLogAmounts(ret.ID, a, shareAmount, guildSharePercent)
	return ret
}

func initGuildIncomeSettlementLog() {
	syndb.RegQuick(TbGuildIncomeSettlementLog, db.CreatedAtName)
	syndb.RegQuick(TbGuildIncomeSettlementLog, db.UpdatedAtName)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogGuildId)
	regLiveRoomIncomeCols(TbGuildIncomeSettlementLog)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogSettlementSalary)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogSettlementShareAmount)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogGuildSharePercent)
	migrate.AutoMigrate(&GuildIncomeSettlementLog{})
}

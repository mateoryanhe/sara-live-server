package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbGuildIncomeTotal db.TbName = "guild_income_totals"
)

// GuildIncomeTotal 工会生涯累计收益(主键ID=工会ID,只增不结算)
type GuildIncomeTotal struct {
	migrate.OneModel
	LiveRoomIncomeAmounts
	SettlementSalary float64 `gorm:"type:decimal(16,4);default:0;comment:结算薪资" json:"settlementSalary"`
}

func NewGuildIncomeTotal(guildId uint64) *GuildIncomeTotal {
	ret := &GuildIncomeTotal{}
	ret.ID = guildId
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	syndb.AddData(TbGuildIncomeTotal, db.CreatedAtName, &syndb.ColData{IdVal: guildId, ColVal: now})
	syndb.AddData(TbGuildIncomeTotal, db.UpdatedAtName, &syndb.ColData{IdVal: guildId, ColVal: now})
	return ret
}

func (r *GuildIncomeTotal) AddTotalLiveDuration(v float64) {
	addIncomeAmount(TbGuildIncomeTotal, LiveRoomIncomeTotalLiveDuration, r.ID, &r.TotalLiveDuration, v, true, &r.UpdatedAt)
}
func (r *GuildIncomeTotal) AddGiftEarn(v float64) {
	addIncomeEarn(TbGuildIncomeTotal, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalGiftIncome, &r.TotalGiftIncome)
}
func (r *GuildIncomeTotal) AddPaidDanmakuEarn(v float64) {
	addIncomeEarn(TbGuildIncomeTotal, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPaidDanmakuIncome, &r.TotalPaidDanmakuIncome)
}
func (r *GuildIncomeTotal) AddPrivateRoomTicketEarn(v float64) {
	addIncomeEarn(TbGuildIncomeTotal, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPrivateRoomTicketIncome, &r.TotalPrivateRoomTicketIncome)
}
func (r *GuildIncomeTotal) AddPrivateRoomWatchEarn(v float64) {
	addIncomeEarn(TbGuildIncomeTotal, r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPrivateRoomWatchIncome, &r.TotalPrivateRoomWatchIncome)
}

// AddSettlementSalary 累加结算薪资
func (r *GuildIncomeTotal) AddSettlementSalary(v float64) {
	if r == nil || v == 0 {
		return
	}
	addIncomeAmount(TbGuildIncomeTotal, LiveRoomIncomeSettlementSalary, r.ID, &r.SettlementSalary, v, false, &r.UpdatedAt)
}

func initGuildIncomeTotal() {
	regLiveRoomIncomeCols(TbGuildIncomeTotal)
	syndb.RegQuick(TbGuildIncomeTotal, LiveRoomIncomeSettlementSalary)
	migrate.AutoMigrate(&GuildIncomeTotal{})
}

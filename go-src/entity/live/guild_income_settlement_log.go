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
	GuildIncomeSettlementLogGuildId                  db.TbCol = "guild_id"
	GuildIncomeSettlementLogSettlementSalary         db.TbCol = "settlement_salary"
	GuildIncomeSettlementLogSettlementShareAmount    db.TbCol = "settlement_share_amount"
	GuildIncomeSettlementLogSettlementShareAmountUsd db.TbCol = "settlement_share_amount_usd"
	GuildIncomeSettlementLogGuildSharePercent        db.TbCol = "guild_share_percent"
	GuildIncomeSettlementLogStatus                   db.TbCol = "status"
	GuildIncomeSettlementLogTransferAt               db.TbCol = "transfer_at"
	GuildIncomeSettlementLogTransferOrderId          db.TbCol = "transfer_order_id"
	GuildIncomeSettlementLogTransferPlatformNo       db.TbCol = "transfer_platform_no"
	GuildIncomeSettlementLogTransferLocalAmount      db.TbCol = "transfer_local_amount"
	GuildIncomeSettlementLogTransferCurrency         db.TbCol = "transfer_currency"
	GuildIncomeSettlementLogTransferFailMsg          db.TbCol = "transfer_fail_msg"
)

const (
	GuildIncomeSettlementStatusPending      uint8 = 0 // 未审核
	GuildIncomeSettlementStatusApproved     uint8 = 1 // 审核通过
	GuildIncomeSettlementStatusTransferred  uint8 = 2 // 转账成功
	GuildIncomeSettlementStatusTransferring uint8 = 3 // 已提交代付,等待回调
)

// GuildIncomeSettlementLog 工会周结算成功日志(每次结算一条,历史留存)
type GuildIncomeSettlementLog struct {
	migrate.OneModel
	CreatedAt time.Time `gorm:"index:idx_gis_guild_status_created,priority:3" json:"-"`
	GuildId   uint64    `gorm:"index:idx_gis_guild_status_created,priority:1;default:0;comment:工会ID" json:"guildId"`
	LiveRoomIncomeAmounts
	SettlementSalary         float64    `gorm:"type:decimal(16,4);default:0;comment:结算薪资" json:"settlementSalary"`
	SettlementShareAmount    float64    `gorm:"type:decimal(16,4);default:0;comment:结算分佣金额" json:"settlementShareAmount"`
	SettlementShareAmountUsd float64    `gorm:"type:decimal(16,4);default:0;comment:结算分佣金额(USD)" json:"settlementShareAmountUsd"`
	SettlementReceivableUsd  float64    `gorm:"type:decimal(16,4);default:0;comment:结算可收金额(USD)=流水分佣+开播薪资" json:"settlementReceivableUsd"`
	GuildSharePercent        float64    `gorm:"type:decimal(6,2);default:0;comment:本次结算工会分佣比例(%)" json:"guildSharePercent"`
	Status                   uint8      `gorm:"index:idx_gis_guild_status_created,priority:2;default:0;comment:状态(0未审核1审核通过2转账成功3代付中)" json:"status"`
	TransferAt               *time.Time `gorm:"comment:转账时间" json:"transferAt"`
	TransferOrderId          string     `gorm:"size:64;default:'';index;comment:HaiPay代付商户单号" json:"transferOrderId"`
	TransferPlatformNo       string     `gorm:"size:64;default:'';comment:HaiPay平台单号" json:"transferPlatformNo"`
	TransferLocalAmount      float64    `gorm:"type:decimal(18,4);default:0;comment:代付当地币金额" json:"transferLocalAmount"`
	TransferCurrency         string     `gorm:"size:16;default:'';comment:代付币种" json:"transferCurrency"`
	TransferFailMsg          string     `gorm:"size:512;default:'';comment:代付失败原因" json:"transferFailMsg"`
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
	ret.Status = GuildIncomeSettlementStatusPending

	syndb.AddData(TbGuildIncomeSettlementLog, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbGuildIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogGuildId, &syndb.ColData{IdVal: ret.ID, ColVal: guildId})
	writeGuildIncomeSettlementLogAmounts(ret.ID, a, salary, shareAmount, shareAmountUsd, receivableUsd, guildSharePercent)
	ret.SetStatus(GuildIncomeSettlementStatusPending)
	return ret
}

func (r *GuildIncomeSettlementLog) SetStatus(v uint8) {
	if r == nil {
		return
	}
	r.Status = v
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogStatus, &syndb.ColData{IdVal: r.ID, ColVal: v})
	syndb.AddData(TbGuildIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
}

func (r *GuildIncomeSettlementLog) SetTransferAt(v *time.Time) {
	if r == nil {
		return
	}
	r.TransferAt = v
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferAt, &syndb.ColData{IdVal: r.ID, ColVal: v})
	syndb.AddData(TbGuildIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
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
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferOrderId, &syndb.ColData{IdVal: r.ID, ColVal: orderId})
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferPlatformNo, &syndb.ColData{IdVal: r.ID, ColVal: platformNo})
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferCurrency, &syndb.ColData{IdVal: r.ID, ColVal: currency})
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferLocalAmount, &syndb.ColData{IdVal: r.ID, ColVal: localAmount})
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferFailMsg, &syndb.ColData{IdVal: r.ID, ColVal: ""})
	syndb.AddData(TbGuildIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
}

func (r *GuildIncomeSettlementLog) SetTransferFailMsg(msg string) {
	if r == nil {
		return
	}
	r.TransferFailMsg = msg
	now := time.Now()
	r.UpdatedAt = now
	syndb.AddData(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferFailMsg, &syndb.ColData{IdVal: r.ID, ColVal: msg})
	syndb.AddData(TbGuildIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: now})
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
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogStatus)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferAt)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferOrderId)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferPlatformNo)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferLocalAmount)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferCurrency)
	syndb.RegQuick(TbGuildIncomeSettlementLog, GuildIncomeSettlementLogTransferFailMsg)
	migrate.AutoMigrate(&GuildIncomeSettlementLog{})
}

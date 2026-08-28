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
	AnchorIncomeSettlementLogRoomId                db.TbCol = "room_id"
	AnchorIncomeSettlementLogSettlementSalary      db.TbCol = "settlement_salary"
	AnchorIncomeSettlementLogSettlementShareAmount    db.TbCol = "settlement_share_amount"
	AnchorIncomeSettlementLogSettlementShareAmountUsd db.TbCol = "settlement_share_amount_usd"
	AnchorIncomeSettlementLogAnchorSharePercent       db.TbCol = "anchor_share_percent"
)

// AnchorIncomeSettlementLog 主播周结算成功日志(每次结算一条,历史留存)
type AnchorIncomeSettlementLog struct {
	migrate.OneModel
	RoomId uint64 `gorm:"index;default:0;comment:直播间ID(==主播用户ID)" json:"roomId"`
	LiveRoomIncomeAmounts
	SettlementSalary          float64 `gorm:"type:decimal(16,4);default:0;comment:本次结算薪资" json:"settlementSalary"`
	SettlementShareAmount     float64 `gorm:"type:decimal(16,4);default:0;comment:本次结算分佣金额" json:"settlementShareAmount"`
	SettlementShareAmountUsd  float64 `gorm:"type:decimal(16,4);default:0;comment:本次结算分佣金额(USD)" json:"settlementShareAmountUsd"`
	AnchorSharePercent        float64 `gorm:"type:decimal(6,2);default:0;comment:本次结算主播分佣比例(%)" json:"anchorSharePercent"`
}

// NewAnchorIncomeSettlementLog 新建一条主播结算日志并入库
func NewAnchorIncomeSettlementLog(roomId uint64, a *LiveRoomIncomeAmounts, salary, shareAmount, shareAmountUsd, anchorSharePercent float64) *AnchorIncomeSettlementLog {
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

	syndb.AddData(TbAnchorIncomeSettlementLog, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbAnchorIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogRoomId, &syndb.ColData{IdVal: ret.ID, ColVal: roomId})
	writeAnchorIncomeSettlementLogAmounts(ret.ID, a, salary, shareAmount, shareAmountUsd, anchorSharePercent)
	return ret
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
	migrate.AutoMigrate(&AnchorIncomeSettlementLog{})
}

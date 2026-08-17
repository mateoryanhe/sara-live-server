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
	AnchorIncomeSettlementLogRoomId           db.TbCol = "room_id"
	AnchorIncomeSettlementLogSettlementSalary db.TbCol = "settlement_salary"
)

// AnchorIncomeSettlementLog 主播周结算成功日志(每次结算一条,历史留存)
type AnchorIncomeSettlementLog struct {
	migrate.OneModel
	RoomId uint64 `gorm:"index;default:0;comment:直播间ID(==主播用户ID)" json:"roomId"`
	LiveRoomIncomeAmounts
	SettlementSalary float64 `gorm:"type:decimal(16,4);default:0;comment:本次结算薪资" json:"settlementSalary"`
}

// NewAnchorIncomeSettlementLog 新建一条主播结算日志并入库
func NewAnchorIncomeSettlementLog(roomId uint64, a *LiveRoomIncomeAmounts, salary float64) *AnchorIncomeSettlementLog {
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

	syndb.AddData(TbAnchorIncomeSettlementLog, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbAnchorIncomeSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogRoomId, &syndb.ColData{IdVal: ret.ID, ColVal: roomId})
	writeIncomeSettlementLogAmounts(TbAnchorIncomeSettlementLog, ret.ID, a, salary)
	return ret
}

func initAnchorIncomeSettlementLog() {
	syndb.RegQuick(TbAnchorIncomeSettlementLog, db.CreatedAtName)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, db.UpdatedAtName)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogRoomId)
	regLiveRoomIncomeCols(TbAnchorIncomeSettlementLog)
	syndb.RegQuick(TbAnchorIncomeSettlementLog, AnchorIncomeSettlementLogSettlementSalary)
	migrate.AutoMigrate(&AnchorIncomeSettlementLog{})
}

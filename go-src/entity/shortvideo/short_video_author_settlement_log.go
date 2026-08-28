package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbShortVideoAuthorSettlementLog db.TbName = "short_video_author_settlement_logs"
)

const (
	ShortVideoAuthorSettlementLogUserId              db.TbCol = "user_id"
	ShortVideoAuthorSettlementLogUnsettledIncome     db.TbCol = "unsettled_income"
	ShortVideoAuthorSettlementLogSettlementDiamond   db.TbCol = "settlement_diamond"
	ShortVideoAuthorSettlementLogAnchorSharePercent  db.TbCol = "anchor_share_percent"
)

// ShortVideoAuthorSettlementLog 非主播作者短视频周结算日志
type ShortVideoAuthorSettlementLog struct {
	migrate.OneModel
	UserId              uint64  `gorm:"index;default:0;comment:作者用户ID" json:"userId"`
	UnsettledIncome     float64 `gorm:"type:decimal(16,4);default:0;comment:本次结算未结算流水" json:"unsettledIncome"`
	SettlementDiamond   float64 `gorm:"type:decimal(16,4);default:0;comment:本次结算到账钻石" json:"settlementDiamond"`
	AnchorSharePercent  float64 `gorm:"type:decimal(6,2);default:0;comment:本次结算主播分佣比例(%)" json:"anchorSharePercent"`
}

// NewShortVideoAuthorSettlementLog 新建一条非主播作者短视频结算日志并入库
func NewShortVideoAuthorSettlementLog(userId uint64, unsettledIncome, settlementDiamond, anchorSharePercent float64) *ShortVideoAuthorSettlementLog {
	ret := &ShortVideoAuthorSettlementLog{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.CreatedAt = now
	ret.UpdatedAt = now
	ret.UserId = userId
	ret.UnsettledIncome = unsettledIncome
	ret.SettlementDiamond = settlementDiamond
	ret.AnchorSharePercent = anchorSharePercent

	syndb.AddData(TbShortVideoAuthorSettlementLog, db.CreatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbShortVideoAuthorSettlementLog, db.UpdatedAtName, &syndb.ColData{IdVal: ret.ID, ColVal: now})
	syndb.AddData(TbShortVideoAuthorSettlementLog, ShortVideoAuthorSettlementLogUserId, &syndb.ColData{IdVal: ret.ID, ColVal: userId})
	syndb.AddData(TbShortVideoAuthorSettlementLog, ShortVideoAuthorSettlementLogUnsettledIncome, &syndb.ColData{IdVal: ret.ID, ColVal: unsettledIncome})
	syndb.AddData(TbShortVideoAuthorSettlementLog, ShortVideoAuthorSettlementLogSettlementDiamond, &syndb.ColData{IdVal: ret.ID, ColVal: settlementDiamond})
	syndb.AddData(TbShortVideoAuthorSettlementLog, ShortVideoAuthorSettlementLogAnchorSharePercent, &syndb.ColData{IdVal: ret.ID, ColVal: anchorSharePercent})
	return ret
}

func initShortVideoAuthorSettlementLog() {
	syndb.RegQuick(TbShortVideoAuthorSettlementLog, db.CreatedAtName)
	syndb.RegQuick(TbShortVideoAuthorSettlementLog, db.UpdatedAtName)
	syndb.RegQuick(TbShortVideoAuthorSettlementLog, ShortVideoAuthorSettlementLogUserId)
	syndb.RegQuick(TbShortVideoAuthorSettlementLog, ShortVideoAuthorSettlementLogUnsettledIncome)
	syndb.RegQuick(TbShortVideoAuthorSettlementLog, ShortVideoAuthorSettlementLogSettlementDiamond)
	syndb.RegQuick(TbShortVideoAuthorSettlementLog, ShortVideoAuthorSettlementLogAnchorSharePercent)
	migrate.AutoMigrate(&ShortVideoAuthorSettlementLog{})
}

package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/constants/gameplatform"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const TbGameWinLog db.TbName = "game_win_logs"

const (
	GameWinLogUserId       db.TbCol = "user_id"
	GameWinLogGameCode     db.TbCol = "game_code"
	GameWinLogNameEn       db.TbCol = "name_en"
	GameWinLogCover        db.TbCol = "cover"
	GameWinLogAmount       db.TbCol = "amount"
	GameWinLogPlatformType db.TbCol = "platform_type"
	GameWinLogOrderId      db.TbCol = "order_id"
)

// GameWinLog 游戏派彩记录
type GameWinLog struct {
	migrate.OneModel
	UserId       uint64  `gorm:"index;default:0;comment:用户ID" json:"userId"`
	GameCode     string  `gorm:"size:64;default:'';comment:游戏编码" json:"gameCode"`
	NameEn       string  `gorm:"size:128;default:'';comment:英文名称" json:"nameEn"`
	Cover        string  `gorm:"size:512;default:'';comment:封面" json:"cover"`
	Amount       float64 `gorm:"default:0;comment:派彩金额" json:"amount"`
	PlatformType string  `gorm:"size:32;default:'';comment:平台类型" json:"platformType"`
	OrderId      string  `gorm:"size:64;default:'';index;comment:订单ID" json:"orderId"`
}

func NewGameWinLog(userId uint64, gameCode, nameEn, cover string, platformType gameplatform.Platform, orderId string, amount float64) *GameWinLog {
	now := time.Now()
	row := &GameWinLog{}
	row.ID = snowflake.GetId()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	row.SetGameCode(gameCode)
	row.SetNameEn(nameEn)
	row.SetCover(cover)
	row.SetAmount(amount)
	row.SetPlatformType(platformType.String())
	row.SetOrderId(orderId)
	return row
}

func (r *GameWinLog) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbGameWinLog, GameWinLogUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameWinLog) SetGameCode(v string) {
	r.GameCode = v
	syndb.AddData(TbGameWinLog, GameWinLogGameCode, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameWinLog) SetNameEn(v string) {
	r.NameEn = v
	syndb.AddData(TbGameWinLog, GameWinLogNameEn, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameWinLog) SetCover(v string) {
	r.Cover = v
	syndb.AddData(TbGameWinLog, GameWinLogCover, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameWinLog) SetAmount(v float64) {
	r.Amount = v
	syndb.AddData(TbGameWinLog, GameWinLogAmount, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameWinLog) SetPlatformType(v string) {
	r.PlatformType = v
	syndb.AddData(TbGameWinLog, GameWinLogPlatformType, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameWinLog) SetOrderId(v string) {
	r.OrderId = v
	syndb.AddData(TbGameWinLog, GameWinLogOrderId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameWinLog) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbGameWinLog, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameWinLog) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbGameWinLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initGameWinLog() {
	syndb.RegQuick(TbGameWinLog, db.CreatedAtName)
	syndb.RegQuick(TbGameWinLog, db.UpdatedAtName)
	syndb.RegQuick(TbGameWinLog, GameWinLogUserId)
	syndb.RegQuick(TbGameWinLog, GameWinLogGameCode)
	syndb.RegQuick(TbGameWinLog, GameWinLogNameEn)
	syndb.RegQuick(TbGameWinLog, GameWinLogCover)
	syndb.RegQuick(TbGameWinLog, GameWinLogAmount)
	syndb.RegQuick(TbGameWinLog, GameWinLogPlatformType)
	syndb.RegQuick(TbGameWinLog, GameWinLogOrderId)
	migrate.AutoMigrate(&GameWinLog{})
}

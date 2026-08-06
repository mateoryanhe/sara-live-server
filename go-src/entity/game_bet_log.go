package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/constants/gameplatform"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const TbGameBetLog db.TbName = "game_bet_logs"

const (
	GameBetLogUserId       db.TbCol = "user_id"
	GameBetLogGameCode     db.TbCol = "game_code"
	GameBetLogNameEn       db.TbCol = "name_en"
	GameBetLogCover        db.TbCol = "cover"
	GameBetLogAmount       db.TbCol = "amount"
	GameBetLogPlatformType db.TbCol = "platform_type"
	GameBetLogOrderId      db.TbCol = "order_id"
	GameBetLogLiveRoomId   db.TbCol = "live_room_id"
	GameBetLogLiveRecordId db.TbCol = "live_record_id"
)

// GameBetLog 游戏下注记录
type GameBetLog struct {
	migrate.OneModel
	UserId       uint64  `gorm:"index;default:0;comment:用户ID" json:"userId"`
	GameCode     string  `gorm:"size:64;default:'';comment:游戏编码" json:"gameCode"`
	NameEn       string  `gorm:"size:128;default:'';comment:英文名称" json:"nameEn"`
	Cover        string  `gorm:"size:512;default:'';comment:封面" json:"cover"`
	Amount       float64 `gorm:"default:0;comment:下注金额" json:"amount"`
	PlatformType string  `gorm:"size:32;default:'';comment:平台类型" json:"platformType"`
	OrderId      string  `gorm:"size:64;default:'';index;comment:订单ID" json:"orderId"`
	LiveRoomId   uint64  `gorm:"index;default:0;comment:直播间ID" json:"liveRoomId"`
	LiveRecordId uint64  `gorm:"index;default:0;comment:直播记录ID" json:"liveRecordId"`
}

func NewGameBetLog(userId uint64, gameCode, nameEn, cover string, platformType gameplatform.Platform, orderId string, amount float64, liveRoomId, liveRecordId uint64) *GameBetLog {
	now := time.Now()
	row := &GameBetLog{}
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
	row.SetLiveRoomId(liveRoomId)
	row.SetLiveRecordId(liveRecordId)
	return row
}

func (r *GameBetLog) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToQuickChan(TbGameBetLog, GameBetLogUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameBetLog) SetGameCode(v string) {
	r.GameCode = v
	syndb.AddDataToQuickChan(TbGameBetLog, GameBetLogGameCode, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameBetLog) SetNameEn(v string) {
	r.NameEn = v
	syndb.AddDataToQuickChan(TbGameBetLog, GameBetLogNameEn, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameBetLog) SetCover(v string) {
	r.Cover = v
	syndb.AddDataToQuickChan(TbGameBetLog, GameBetLogCover, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameBetLog) SetAmount(v float64) {
	r.Amount = v
	syndb.AddDataToQuickChan(TbGameBetLog, GameBetLogAmount, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameBetLog) SetPlatformType(v string) {
	r.PlatformType = v
	syndb.AddDataToQuickChan(TbGameBetLog, GameBetLogPlatformType, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameBetLog) SetOrderId(v string) {
	r.OrderId = v
	syndb.AddDataToQuickChan(TbGameBetLog, GameBetLogOrderId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameBetLog) SetLiveRoomId(v uint64) {
	r.LiveRoomId = v
	syndb.AddDataToQuickChan(TbGameBetLog, GameBetLogLiveRoomId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameBetLog) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	syndb.AddDataToQuickChan(TbGameBetLog, GameBetLogLiveRecordId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameBetLog) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddDataToQuickChan(TbGameBetLog, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *GameBetLog) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddDataToQuickChan(TbGameBetLog, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initGameBetLog() {
	syndb.RegQuick(TbGameBetLog, db.CreatedAtName)
	syndb.RegQuick(TbGameBetLog, db.UpdatedAtName)
	syndb.RegQuick(TbGameBetLog, GameBetLogUserId)
	syndb.RegQuick(TbGameBetLog, GameBetLogGameCode)
	syndb.RegQuick(TbGameBetLog, GameBetLogNameEn)
	syndb.RegQuick(TbGameBetLog, GameBetLogCover)
	syndb.RegQuick(TbGameBetLog, GameBetLogAmount)
	syndb.RegQuick(TbGameBetLog, GameBetLogPlatformType)
	syndb.RegQuick(TbGameBetLog, GameBetLogOrderId)
	syndb.RegQuick(TbGameBetLog, GameBetLogLiveRoomId)
	syndb.RegQuick(TbGameBetLog, GameBetLogLiveRecordId)
	migrate.AutoMigrate(&GameBetLog{})
}

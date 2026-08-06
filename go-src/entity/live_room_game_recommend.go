package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbLiveRoomGameRecommend db.TbName = "live_room_game_recommends"

const (
	LiveRoomGameRecommendLiveRoomId db.TbCol = "live_room_id"
	LiveRoomGameRecommendGameCode   db.TbCol = "game_code"
	LiveRoomGameRecommendStatus     db.TbCol = "status"
)

const (
	LiveRoomGameRecommendStatusDeleted uint8 = 0
	LiveRoomGameRecommendStatusActive  uint8 = 1
)

// LiveRoomGameRecommend 直播间游戏推荐
type LiveRoomGameRecommend struct {
	migrate.OneModel
	LiveRoomId uint64 `gorm:"index:idx_room_game,priority:1;default:0;comment:直播间ID" json:"liveRoomId"`
	GameCode   string `gorm:"index:idx_room_game,priority:2;size:64;comment:游戏编码" json:"gameCode"`
	Status     uint8  `gorm:"default:1;comment:状态(0删除,1有效)" json:"status"`
}

func (r *LiveRoomGameRecommend) SetLiveRoomId(v uint64) {
	r.LiveRoomId = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoomGameRecommend, LiveRoomGameRecommendLiveRoomId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomGameRecommend) SetGameCode(v string) {
	r.GameCode = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoomGameRecommend, LiveRoomGameRecommendGameCode, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomGameRecommend) SetStatus(v uint8) {
	r.Status = v
	r.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbLiveRoomGameRecommend, LiveRoomGameRecommendStatus, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomGameRecommend) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddDataToQuickChan(TbLiveRoomGameRecommend, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomGameRecommend) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddDataToQuickChan(TbLiveRoomGameRecommend, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoomGameRecommend) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	syndb.AddDataToQuickChan(TbLiveRoomGameRecommend, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.UpdatedAt,
	})
}

func initLiveRoomGameRecommend() {
	syndb.RegQuick(TbLiveRoomGameRecommend, db.CreatedAtName)
	syndb.RegQuick(TbLiveRoomGameRecommend, db.UpdatedAtName)
	syndb.RegQuick(TbLiveRoomGameRecommend, LiveRoomGameRecommendLiveRoomId)
	syndb.RegQuick(TbLiveRoomGameRecommend, LiveRoomGameRecommendGameCode)
	syndb.RegQuick(TbLiveRoomGameRecommend, LiveRoomGameRecommendStatus)
	migrate.AutoMigrate(&LiveRoomGameRecommend{})
}

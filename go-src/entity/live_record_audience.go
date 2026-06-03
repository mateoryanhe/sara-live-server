package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRecordAudience db.TbName = "live_record_audiences"
)

const (
	LiveRecordAudienceLiveRecordId db.TbCol = "live_record_id"
	LiveRecordAudienceUserId       db.TbCol = "user_id"
)

// LiveRecordAudience 单场直播观众去重记录(主键 ID = liveRecordId_userId)
type LiveRecordAudience struct {
	ID           string    `gorm:"primaryKey;size:64;comment:复合ID(liveRecordId_userId)" json:"id"`
	LiveRecordId uint64    `gorm:"index;default:0;comment:直播记录ID" json:"liveRecordId"`
	UserId       uint64    `gorm:"index;default:0;comment:观众用户ID" json:"userId"`
	CreatedAt    time.Time `json:"createdAt"`
}

func BuildLiveRecordAudienceId(liveRecordId, userId uint64) string {
	return fmt.Sprintf("%d_%d", liveRecordId, userId)
}

func NewLiveRecordAudience(liveRecordId, userId uint64) *LiveRecordAudience {
	r := &LiveRecordAudience{}
	r.ID = BuildLiveRecordAudienceId(liveRecordId, userId)
	r.SetLiveRecordId(liveRecordId)
	r.SetUserId(userId)
	return r
}

func (r *LiveRecordAudience) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	syndb.AddDataToLazyChan(TbLiveRecordAudience, LiveRecordAudienceLiveRecordId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRecordAudience) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbLiveRecordAudience, LiveRecordAudienceUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRecordAudience) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbLiveRecordAudience, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initLiveRecordAudience() {
	syndb.RegLazyWithMiddle(TbLiveRecordAudience, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbLiveRecordAudience, LiveRecordAudienceLiveRecordId)
	syndb.RegLazyWithMiddle(TbLiveRecordAudience, LiveRecordAudienceUserId)
	migrate.AutoMigrate(&LiveRecordAudience{})
}

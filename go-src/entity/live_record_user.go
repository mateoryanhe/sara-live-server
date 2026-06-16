package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRecordUser db.TbName = "live_record_users"
)

const (
	LiveRecordUserLiveRecordId db.TbCol = "live_record_id"
	LiveRecordUserUserId       db.TbCol = "user_id"
	LiveRecordUserAudienceAt   db.TbCol = "audience_at"
	LiveRecordUserGiftSenderAt db.TbCol = "gift_sender_at"
)

// LiveRecordUser 单场直播用户去重记录(主键 ID = liveRecordId_userId)
type LiveRecordUser struct {
	ID           string    `gorm:"primaryKey;size:64;comment:复合ID(liveRecordId_userId)" json:"id"`
	LiveRecordId uint64    `gorm:"index;default:0;comment:直播记录ID" json:"liveRecordId"`
	UserId       uint64    `gorm:"index;default:0;comment:用户ID" json:"userId"`
	AudienceAt   time.Time `gorm:"comment:观众统计时间(零值表示未统计)" json:"audienceAt"`
	GiftSenderAt time.Time `gorm:"comment:送礼人统计时间(零值表示未统计)" json:"giftSenderAt"`
}

func BuildLiveRecordUserId(liveRecordId, userId uint64) string {
	return fmt.Sprintf("%d_%d", liveRecordId, userId)
}

func NewLiveRecordUser(liveRecordId, userId uint64) *LiveRecordUser {
	r := &LiveRecordUser{}
	r.ID = BuildLiveRecordUserId(liveRecordId, userId)
	r.SetLiveRecordId(liveRecordId)
	r.SetUserId(userId)
	return r
}

func (r *LiveRecordUser) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	syndb.AddDataToLazyChan(TbLiveRecordUser, LiveRecordUserLiveRecordId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRecordUser) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToLazyChan(TbLiveRecordUser, LiveRecordUserUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRecordUser) SetAudienceAt(v time.Time) {
	r.AudienceAt = v
	syndb.AddDataToLazyChan(TbLiveRecordUser, LiveRecordUserAudienceAt, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRecordUser) SetGiftSenderAt(v time.Time) {
	r.GiftSenderAt = v
	syndb.AddDataToLazyChan(TbLiveRecordUser, LiveRecordUserGiftSenderAt, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initLiveRecordUser() {
	syndb.RegLazyWithMiddle(TbLiveRecordUser, LiveRecordUserLiveRecordId)
	syndb.RegLazyWithMiddle(TbLiveRecordUser, LiveRecordUserUserId)
	syndb.RegLazyWithMiddle(TbLiveRecordUser, LiveRecordUserAudienceAt)
	syndb.RegLazyWithMiddle(TbLiveRecordUser, LiveRecordUserGiftSenderAt)
	migrate.AutoMigrate(&LiveRecordUser{})
}

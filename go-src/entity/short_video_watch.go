package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbShortVideoWatch db.TbName = "short_video_watches"
)

const (
	ShortVideoWatchUserId        db.TbCol = "user_id"
	ShortVideoWatchVideoId       db.TbCol = "video_id"
	ShortVideoWatchBilledSeconds db.TbCol = "billed_seconds"
)

// ShortVideoWatch 短视频观看扣费进度(按用户+视频维度累计)
type ShortVideoWatch struct {
	ID            string    `gorm:"primaryKey;size:64;comment:复合ID(userId_videoId)" json:"id"`
	UserId        uint64    `gorm:"index:idx_user_video,priority:1;default:0;comment:用户ID" json:"userId"`
	VideoId       uint64    `gorm:"index:idx_user_video,priority:2;default:0;comment:短视频ID" json:"videoId"`
	BilledSeconds uint32    `gorm:"default:0;comment:已累计计费观看秒数" json:"billedSeconds"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func BuildShortVideoWatchId(userId, videoId uint64) string {
	return fmt.Sprintf("%d_%d", userId, videoId)
}

func NewShortVideoWatch(userId, videoId uint64) *ShortVideoWatch {
	watch := &ShortVideoWatch{}
	watch.ID = BuildShortVideoWatchId(userId, videoId)
	now := time.Now()
	watch.SetCreatedAt(now)
	watch.SetUpdatedAt(now)
	watch.SetUserId(userId)
	watch.SetVideoId(videoId)
	watch.SetBilledSeconds(0)
	return watch
}

func (watch *ShortVideoWatch) SetUserId(v uint64) {
	watch.UserId = v
	syndb.AddDataToQuickChan(TbShortVideoWatch, ShortVideoWatchUserId, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})
}

func (watch *ShortVideoWatch) SetVideoId(v uint64) {
	watch.VideoId = v
	syndb.AddDataToQuickChan(TbShortVideoWatch, ShortVideoWatchVideoId, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})
}

func (watch *ShortVideoWatch) SetBilledSeconds(v uint32) {
	watch.BilledSeconds = v
	watch.UpdatedAt = time.Now()
	syndb.AddDataToQuickChan(TbShortVideoWatch, ShortVideoWatchBilledSeconds, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})
	syndb.AddDataToQuickChan(TbShortVideoWatch, db.UpdatedAtName, &syndb.ColData{
		IdVal: watch.ID, ColVal: watch.UpdatedAt,
	})
}

func (watch *ShortVideoWatch) SetCreatedAt(v time.Time) {
	watch.CreatedAt = v
	syndb.AddDataToQuickChan(TbShortVideoWatch, db.CreatedAtName, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})
}

func (watch *ShortVideoWatch) SetUpdatedAt(v time.Time) {
	watch.UpdatedAt = v
	syndb.AddDataToQuickChan(TbShortVideoWatch, db.UpdatedAtName, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})
}

func initShortVideoWatch() {
	syndb.RegQuickWithLarge(TbShortVideoWatch, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbShortVideoWatch, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbShortVideoWatch, ShortVideoWatchUserId)
	syndb.RegQuickWithLarge(TbShortVideoWatch, ShortVideoWatchVideoId)
	syndb.RegQuickWithLarge(TbShortVideoWatch, ShortVideoWatchBilledSeconds)
	migrate.AutoMigrate(&ShortVideoWatch{})
}

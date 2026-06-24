package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
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
	ShortVideoWatchWatchSeconds  db.TbCol = "watch_seconds"
	ShortVideoWatchViewCounted   db.TbCol = "view_counted"
	ShortVideoLikeStatus         db.TbCol = "status"
)

const (
	ShortVideoWatchViewCountedNo  uint8 = 0 // 未计入观看人数
	ShortVideoWatchViewCountedYes uint8 = 1 // 已计入观看人数

	ShortVideoLikeStatusCancelled uint8 = 0 // 已取消点赞
	ShortVideoLikeStatusLiked     uint8 = 1 // 已点赞
)

// ShortVideoWatch 短视频观看扣费进度(按用户+视频维度累计)
type ShortVideoWatch struct {
	ID            string    `gorm:"primaryKey;comment:复合ID(userId_videoId)" json:"id"`
	UserId        uint64    `gorm:"index:idx_user_video,priority:1;default:0;comment:用户ID" json:"userId"`
	VideoId       uint64    `gorm:"index:idx_user_video,priority:2;default:0;comment:短视频ID" json:"videoId"`
	BilledSeconds uint64    `gorm:"default:0;comment:已累计计费观看秒数" json:"billedSeconds"`
	WatchSeconds  uint64    `gorm:"default:0;comment:累计观看时长(秒,含免费与付费)" json:"watchSeconds"`
	ViewCounted   uint8     `gorm:"default:0;comment:是否已计入观看人数(0否,1是)" json:"viewCounted"`
	Status        uint8     `gorm:"index:idx_user_video,priority:3;default:0;comment:状态(0已取消,1已点赞)" json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewShortVideoWatch(userId, videoId uint64) *ShortVideoWatch {
	watch := &ShortVideoWatch{}
	watch.ID = fmt.Sprintf("%v_%v", userId, videoId)
	now := time.Now()
	watch.SetCreatedAt(now)
	watch.SetUpdatedAt(now)
	watch.SetUserId(userId)
	watch.SetVideoId(videoId)

	watch.SetViewCounted(ShortVideoWatchViewCountedNo)
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

func (watch *ShortVideoWatch) SetViewCounted(v uint8) {
	watch.ViewCounted = v

	syndb.AddDataToLazyChan(TbShortVideoWatch, ShortVideoWatchViewCounted, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})

}

func (watch *ShortVideoWatch) AddBilledSeconds(v uint64) {
	watch.BilledSeconds = math.Add(watch.BilledSeconds, v)
	syndb.AddDataToLazyChan(TbShortVideoWatch, ShortVideoWatchBilledSeconds, &syndb.ColData{
		IdVal: watch.ID, ColVal: watch.BilledSeconds,
	})
}

func (watch *ShortVideoWatch) AddWatchSeconds(v uint64) {
	watch.WatchSeconds = math.Add(watch.WatchSeconds, v)
	syndb.AddDataToLazyChan(TbShortVideoWatch, ShortVideoWatchWatchSeconds, &syndb.ColData{
		IdVal: watch.ID, ColVal: watch.WatchSeconds,
	})
}

func (watch *ShortVideoWatch) SubWatchSeconds(v uint64) {
	watch.WatchSeconds = math.Sub(watch.WatchSeconds, v)
	syndb.AddDataToLazyChan(TbShortVideoWatch, ShortVideoWatchWatchSeconds, &syndb.ColData{
		IdVal: watch.ID, ColVal: watch.WatchSeconds,
	})
}

func (watch *ShortVideoWatch) ResetWatchSeconds() {
	watch.WatchSeconds = 0
	syndb.AddDataToLazyChan(TbShortVideoWatch, ShortVideoWatchWatchSeconds, &syndb.ColData{
		IdVal: watch.ID, ColVal: 0,
	})
}

func (like *ShortVideoWatch) SetStatus(v uint8) {
	like.Status = v

	syndb.AddDataToQuickChan(TbShortVideoWatch, ShortVideoLikeStatus, &syndb.ColData{
		IdVal: like.ID, ColVal: v,
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
	syndb.AddDataToLazyChan(TbShortVideoWatch, db.UpdatedAtName, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})
}

func initShortVideoWatch() {
	syndb.RegQuickWithLarge(TbShortVideoWatch, db.CreatedAtName)
	syndb.RegLazyWithLarge(TbShortVideoWatch, db.UpdatedAtName)

	syndb.RegQuickWithLarge(TbShortVideoWatch, ShortVideoWatchUserId)
	syndb.RegQuickWithLarge(TbShortVideoWatch, ShortVideoWatchVideoId)
	syndb.RegLazyWithLarge(TbShortVideoWatch, ShortVideoWatchBilledSeconds)
	syndb.RegLazyWithLarge(TbShortVideoWatch, ShortVideoWatchWatchSeconds)
	syndb.RegLazyWithLarge(TbShortVideoWatch, ShortVideoWatchViewCounted)

	syndb.RegQuickWithLarge(TbShortVideoWatch, ShortVideoLikeStatus)

	migrate.AutoMigrate(&ShortVideoWatch{})
}

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
	ShortVideoWatchUserId      db.TbCol = "user_id"
	ShortVideoWatchVideoId     db.TbCol = "video_id"
	ShortVideoWatchPaidTime    db.TbCol = "paid_time"
	ShortVideoWatchViewCounted db.TbCol = "view_counted"
	ShortVideoLikeStatus       db.TbCol = "status"
)

const (
	ShortVideoWatchViewCountedNo  uint8 = 0 // 未计入观看人数
	ShortVideoWatchViewCountedYes uint8 = 1 // 已计入观看人数

	ShortVideoLikeStatusCancelled uint8 = 0 // 已取消点赞
	ShortVideoLikeStatusLiked     uint8 = 1 // 已点赞
)

// ShortVideoWatch 短视频观看扣费进度(按用户+视频维度累计)
type ShortVideoWatch struct {
	ID          string     `gorm:"primaryKey;comment:复合ID(userId_videoId)" json:"id"`
	UserId      uint64     `gorm:"index:idx_user_video,priority:1;default:0;comment:用户ID" json:"userId"`
	VideoId     uint64     `gorm:"index:idx_user_video,priority:2;index;default:0;comment:短视频ID" json:"videoId"`
	PaidTime    *time.Time `gorm:"index;comment:付费时间" json:"paidTime"`
	ViewCounted uint8      `gorm:"default:0;comment:是否已计入观看人数(0否,1是)" json:"viewCounted"`
	Status      uint8      `gorm:"index:idx_user_video,priority:3;default:0;comment:状态(0已取消,1已点赞)" json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
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
	syndb.AddData(TbShortVideoWatch, ShortVideoWatchUserId, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})
}

func (watch *ShortVideoWatch) SetVideoId(v uint64) {
	watch.VideoId = v
	syndb.AddData(TbShortVideoWatch, ShortVideoWatchVideoId, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})
}

func (watch *ShortVideoWatch) SetPaidTime(v *time.Time) {
	watch.PaidTime = v
	syndb.AddData(TbShortVideoWatch, ShortVideoWatchPaidTime, &syndb.ColData{
		IdVal: watch.ID, ColVal: watch.PaidTime,
	})
}

func (watch *ShortVideoWatch) SetViewCounted(v uint8) {
	watch.ViewCounted = v

	syndb.AddData(TbShortVideoWatch, ShortVideoWatchViewCounted, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})

}

func (like *ShortVideoWatch) SetStatus(v uint8) {
	like.Status = v

	syndb.AddData(TbShortVideoWatch, ShortVideoLikeStatus, &syndb.ColData{
		IdVal: like.ID, ColVal: v,
	})
}

func (watch *ShortVideoWatch) SetCreatedAt(v time.Time) {
	watch.CreatedAt = v
	syndb.AddData(TbShortVideoWatch, db.CreatedAtName, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})
}

func (watch *ShortVideoWatch) SetUpdatedAt(v time.Time) {
	watch.UpdatedAt = v
	syndb.AddData(TbShortVideoWatch, db.UpdatedAtName, &syndb.ColData{
		IdVal: watch.ID, ColVal: v,
	})
}

func initShortVideoWatch() {
	syndb.RegQuick(TbShortVideoWatch, db.CreatedAtName)
	syndb.RegLazy(TbShortVideoWatch, db.UpdatedAtName)

	syndb.RegQuick(TbShortVideoWatch, ShortVideoWatchUserId)
	syndb.RegQuick(TbShortVideoWatch, ShortVideoWatchVideoId)
	syndb.RegLazy(TbShortVideoWatch, ShortVideoWatchPaidTime)
	syndb.RegLazy(TbShortVideoWatch, ShortVideoWatchViewCounted)

	syndb.RegQuick(TbShortVideoWatch, ShortVideoLikeStatus)

	migrate.AutoMigrate(&ShortVideoWatch{})
}

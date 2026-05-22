package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbShortVideoLike db.TbName = "short_video_likes"
)

const (
	ShortVideoLikeUserId  db.TbCol = "user_id"
	ShortVideoLikeVideoId db.TbCol = "video_id"
	ShortVideoLikeStatus  db.TbCol = "status"
)

const (
	ShortVideoLikeStatusCancelled uint8 = 0 // 已取消点赞
	ShortVideoLikeStatusLiked     uint8 = 1 // 已点赞
)

// ShortVideoLike 短视频点赞记录
// 主键 ID 使用复合字符串: "{userId}_{videoId}",同一用户/视频仅一行,通过 Status 区分
type ShortVideoLike struct {
	ID        string    `gorm:"primaryKey;size:64;comment:复合ID(userId_videoId)" json:"id"`
	UserId    uint64    `gorm:"index:idx_user_video,priority:1;default:0;comment:用户ID" json:"userId"`
	VideoId   uint64    `gorm:"index:idx_user_video,priority:2;default:0;comment:短视频ID" json:"videoId"`
	Status    uint8     `gorm:"index:idx_user_video,priority:3;default:0;comment:状态(0已取消,1已点赞)" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func BuildShortVideoLikeId(userId, videoId uint64) string {
	return fmt.Sprintf("%d_%d", userId, videoId)
}

func NewShortVideoLike(userId, videoId uint64) *ShortVideoLike {
	like := &ShortVideoLike{}
	like.ID = BuildShortVideoLikeId(userId, videoId)
	now := time.Now()
	like.SetCreatedAt(now)
	like.SetUpdatedAt(now)
	like.SetUserId(userId)
	like.SetVideoId(videoId)
	like.SetStatus(ShortVideoLikeStatusLiked)
	return like
}

func (like *ShortVideoLike) SetUserId(v uint64) {
	like.UserId = v
	syndb.AddDataToQuickChan(TbShortVideoLike, ShortVideoLikeUserId, &syndb.ColData{
		IdVal: like.ID, ColVal: v,
	})
}

func (like *ShortVideoLike) SetVideoId(v uint64) {
	like.VideoId = v
	syndb.AddDataToQuickChan(TbShortVideoLike, ShortVideoLikeVideoId, &syndb.ColData{
		IdVal: like.ID, ColVal: v,
	})
}

func (like *ShortVideoLike) SetStatus(v uint8) {
	like.Status = v
	like.UpdatedAt = time.Now()
	syndb.AddDataToQuickChan(TbShortVideoLike, ShortVideoLikeStatus, &syndb.ColData{
		IdVal: like.ID, ColVal: v,
	})
	syndb.AddDataToQuickChan(TbShortVideoLike, db.UpdatedAtName, &syndb.ColData{
		IdVal: like.ID, ColVal: like.UpdatedAt,
	})
}

func (like *ShortVideoLike) SetCreatedAt(v time.Time) {
	like.CreatedAt = v
	syndb.AddDataToQuickChan(TbShortVideoLike, db.CreatedAtName, &syndb.ColData{
		IdVal: like.ID, ColVal: v,
	})
}

func (like *ShortVideoLike) SetUpdatedAt(v time.Time) {
	like.UpdatedAt = v
	syndb.AddDataToQuickChan(TbShortVideoLike, db.UpdatedAtName, &syndb.ColData{
		IdVal: like.ID, ColVal: v,
	})
}

func initShortVideoLike() {
	syndb.RegQuickWithLarge(TbShortVideoLike, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbShortVideoLike, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbShortVideoLike, ShortVideoLikeUserId)
	syndb.RegQuickWithLarge(TbShortVideoLike, ShortVideoLikeVideoId)
	syndb.RegQuickWithLarge(TbShortVideoLike, ShortVideoLikeStatus)
	migrate.AutoMigrate(&ShortVideoLike{})
}

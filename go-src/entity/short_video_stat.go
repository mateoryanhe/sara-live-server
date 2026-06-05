package entity

import (
	"time"
	"xr-game-server/constants/db"
	xrmath "xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbShortVideoStat db.TbName = "short_video_stats"
)

const (
	ShortVideoStatLikeCount db.TbCol = "like_count"
	ShortVideoStatViewCount db.TbCol = "view_count"
)

// ShortVideoStat 短视频统计数据(与短视频一一对应,主键ID即视频ID)
type ShortVideoStat struct {
	migrate.OneModel
	LikeCount uint64 `gorm:"default:0;comment:点赞累计数量" json:"likeCount"`
	ViewCount uint64 `gorm:"default:0;comment:观看人数" json:"viewCount"`
}

func NewShortVideoStat(videoId uint64) *ShortVideoStat {
	ret := &ShortVideoStat{}
	ret.ID = videoId
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

func (s *ShortVideoStat) AddLikeCount(val uint64) {
	s.LikeCount = xrmath.Add(s.LikeCount, val)

	syndb.AddDataToLazyChan(TbShortVideoStat, ShortVideoStatLikeCount, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.LikeCount,
	})
}

func (s *ShortVideoStat) AddViewCount(val uint64) {
	s.ViewCount = xrmath.Add(s.ViewCount, val)

	syndb.AddDataToLazyChan(TbShortVideoStat, ShortVideoStatViewCount, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.ViewCount,
	})
}

func (s *ShortVideoStat) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
	syndb.AddDataToLazyChan(TbShortVideoStat, db.CreatedAtName, &syndb.ColData{
		ColVal: val,
		IdVal:  s.ID,
	})
}

func (s *ShortVideoStat) SetUpdatedAt(val time.Time) {
	s.UpdatedAt = val
	syndb.AddDataToLazyChan(TbShortVideoStat, db.UpdatedAtName, &syndb.ColData{
		ColVal: val,
		IdVal:  s.ID,
	})
}

func initShortVideoStat() {
	syndb.RegLazyWithMiddle(TbShortVideoStat, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbShortVideoStat, db.UpdatedAtName)
	syndb.RegLazyWithMiddle(TbShortVideoStat, ShortVideoStatLikeCount)
	syndb.RegLazyWithMiddle(TbShortVideoStat, ShortVideoStatViewCount)
	migrate.AutoMigrate(&ShortVideoStat{})
}

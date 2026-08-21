package entity

import (
	"time"
	"xr-game-server/constants/db"
	xrmath "xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbShortVideoAuthorStat db.TbName = "short_video_author_stats"
)

const (
	ShortVideoAuthorStatLikeCount          db.TbCol = "like_count"
	ShortVideoAuthorStatViewCount          db.TbCol = "view_count"
	ShortVideoAuthorStatTotalDiamondIncome db.TbCol = "total_diamond_income"
)

// ShortVideoAuthorStat 短视频作者统计数据(与作者一一对应,主键ID即作者用户ID)
type ShortVideoAuthorStat struct {
	migrate.OneModel
	LikeCount          uint64  `gorm:"default:0;comment:累计点赞总数" json:"likeCount"`
	ViewCount          uint64  `gorm:"default:0;comment:累计短视频观看人数(去重)" json:"viewCount"`
	TotalDiamondIncome float64 `gorm:"type:decimal(16,4);default:0;comment:累计短视频收入总额" json:"totalDiamondIncome"`
}

func NewShortVideoAuthorStat(authorId uint64) *ShortVideoAuthorStat {
	ret := &ShortVideoAuthorStat{}
	ret.ID = authorId
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

func (s *ShortVideoAuthorStat) AddLikeCount(val uint64) {
	s.LikeCount = xrmath.Add(s.LikeCount, val)

	syndb.AddData(TbShortVideoAuthorStat, ShortVideoAuthorStatLikeCount, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.LikeCount,
	})
}

func (s *ShortVideoAuthorStat) AddViewCount(val uint64) {
	s.ViewCount = xrmath.Add(s.ViewCount, val)

	syndb.AddData(TbShortVideoAuthorStat, ShortVideoAuthorStatViewCount, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.ViewCount,
	})
}

func (s *ShortVideoAuthorStat) AddTotalDiamondIncome(val float64) {
	if val <= 0 {
		return
	}
	s.TotalDiamondIncome = xrmath.AddFloat64(s.TotalDiamondIncome, val)

	syndb.AddData(TbShortVideoAuthorStat, ShortVideoAuthorStatTotalDiamondIncome, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalDiamondIncome,
	})
}

func (s *ShortVideoAuthorStat) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
	syndb.AddData(TbShortVideoAuthorStat, db.CreatedAtName, &syndb.ColData{
		ColVal: val,
		IdVal:  s.ID,
	})
}

func (s *ShortVideoAuthorStat) SetUpdatedAt(val time.Time) {
	s.UpdatedAt = val
	syndb.AddData(TbShortVideoAuthorStat, db.UpdatedAtName, &syndb.ColData{
		ColVal: val,
		IdVal:  s.ID,
	})
}

func initShortVideoAuthorStat() {
	syndb.RegLazy(TbShortVideoAuthorStat, db.CreatedAtName)
	syndb.RegLazy(TbShortVideoAuthorStat, db.UpdatedAtName)
	syndb.RegLazy(TbShortVideoAuthorStat, ShortVideoAuthorStatLikeCount)
	syndb.RegLazy(TbShortVideoAuthorStat, ShortVideoAuthorStatViewCount)
	syndb.RegLazy(TbShortVideoAuthorStat, ShortVideoAuthorStatTotalDiamondIncome)
	migrate.AutoMigrate(&ShortVideoAuthorStat{})
}

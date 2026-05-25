package entity

import (
	"time"
	"xr-game-server/constants/db"
	xrmath "xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbShortVideo db.TbName = "short_videos"
)

const (
	ShortVideoLikeCount db.TbCol = "like_count"
)

const (
	ShortVideoStatusOffShelf uint8 = 0
	ShortVideoStatusOnShelf  uint8 = 1
)

const (
	ShortVideoPaidNo  uint8 = 0 // 免费
	ShortVideoPaidYes uint8 = 1 // 付费
)

// ShortVideo 短视频(CMS 管理)
type ShortVideo struct {
	migrate.OneModel
	Title            string `gorm:"size:64;comment:标题" json:"title"`
	Video            string `gorm:"size:255;default:'';comment:视频资源名" json:"video"`
	Cover            string `gorm:"size:255;default:'';comment:封面资源名" json:"cover"`
	Sort             int    `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status           uint8  `gorm:"default:0;comment:状态(0-下架,1-上架)" json:"status"`
	IsPaid           uint8  `gorm:"default:0;comment:是否付费(0免费,1付费)" json:"isPaid"`
	DiamondPerSecond uint64 `gorm:"default:0;comment:每秒钻石数(付费时有效)" json:"diamondPerSecond"`
	Description      string `gorm:"size:255;default:'';comment:描述" json:"description"`
	LikeCount        uint64 `gorm:"default:0;comment:点赞累计数量" json:"likeCount"`
}

func (receiver *ShortVideo) AddLikeCount(val uint64) {
	receiver.LikeCount = xrmath.Add(receiver.LikeCount, val)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbShortVideo, ShortVideoLikeCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.LikeCount,
	})
}

func (receiver *ShortVideo) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToLazyChan(TbShortVideo, db.UpdatedAtName, &syndb.ColData{
		ColVal: val,
		IdVal:  receiver.ID,
	})
}

func initShortVideo() {
	syndb.RegLazyWithMiddle(TbShortVideo, ShortVideoLikeCount)
	syndb.RegLazyWithMiddle(TbShortVideo, db.UpdatedAtName)
	migrate.AutoMigrate(&ShortVideo{})
}

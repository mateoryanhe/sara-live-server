package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbShortVideo db.TbName = "short_videos"
)

const (
	ShortVideoTitle            db.TbCol = "title"
	ShortVideoVideo            db.TbCol = "video"
	ShortVideoCover            db.TbCol = "cover"
	ShortVideoSort             db.TbCol = "sort"
	ShortVideoStatusCol        db.TbCol = "status"
	ShortVideoIsPaid           db.TbCol = "is_paid"
	ShortVideoDiamondPerMinute db.TbCol = "diamond_per_minute"
	ShortVideoCategoryId       db.TbCol = "category_id"
	ShortVideoSource           db.TbCol = "source"
	ShortVideoAuthorId         db.TbCol = "author_id"
)

const (
	ShortVideoStatusOffShelf uint8 = 0
	ShortVideoStatusOnShelf  uint8 = 1
)

const (
	ShortVideoPaidNo  uint8 = 0 // 免费
	ShortVideoPaidYes uint8 = 1 // 付费
)

// 视频来源
const (
	ShortVideoSourceOriginal uint8 = 1 // 原创
	ShortVideoSourceRepost   uint8 = 2 // 转发
	ShortVideoSourceAIGen    uint8 = 3 // AI生成
)

// ShortVideo 短视频(CMS 管理)
type ShortVideo struct {
	migrate.OneModel
	Title            string  `gorm:"size:64;comment:标题" json:"title"`
	Video            string  `gorm:"size:255;default:'';comment:视频资源名" json:"video"`
	Cover            string  `gorm:"size:255;default:'';comment:封面资源名" json:"cover"`
	Sort             int     `gorm:"default:0;comment:排序值(越大越靠前)" json:"sort"`
	Status           uint8   `gorm:"default:0;comment:状态(0-下架,1-上架)" json:"status"`
	IsPaid           uint8   `gorm:"default:0;comment:是否付费(0免费,1付费)" json:"isPaid"`
	DiamondPerMinute float64 `gorm:"type:decimal(10,4);default:0;comment:每分钟钻石数(付费时有效)" json:"diamondPerMinute"`
	CategoryId       int     `gorm:"default:0;comment:视频分类ID" json:"categoryId"`
	Source           uint8   `gorm:"default:1;comment:视频来源(1原创,2转发,3AI生成)" json:"source"`
	AuthorId         uint64  `gorm:"default:0;comment:作者用户ID" json:"authorId"`
}

// NewShortVideo 构造内存对象,字段通过 syndb 异步入库
func NewShortVideo(id uint64, title, video, cover string, sort int, isPaid uint8, diamondPerMinute float64, categoryId int, source uint8, authorId uint64) *ShortVideo {
	v := &ShortVideo{}
	v.ID = id
	now := time.Now()
	v.SetCreatedAt(now)
	v.SetUpdatedAt(now)
	v.SetTitle(title)
	v.SetVideo(video)
	v.SetCover(cover)
	v.SetSort(sort)
	v.SetStatus(ShortVideoStatusOffShelf)
	v.SetIsPaid(isPaid)
	v.SetDiamondPerMinute(diamondPerMinute)
	v.SetCategoryId(categoryId)
	v.SetSource(source)
	v.SetAuthorId(authorId)
	return v
}

func (v *ShortVideo) SetTitle(val string) {
	v.Title = val
	v.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbShortVideo, ShortVideoTitle, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetVideo(val string) {
	v.Video = val
	v.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbShortVideo, ShortVideoVideo, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetCover(val string) {
	v.Cover = val
	v.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbShortVideo, ShortVideoCover, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetSort(val int) {
	v.Sort = val
	v.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbShortVideo, ShortVideoSort, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetStatus(val uint8) {
	v.Status = val
	v.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbShortVideo, ShortVideoStatusCol, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetIsPaid(val uint8) {
	v.IsPaid = val
	v.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbShortVideo, ShortVideoIsPaid, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetDiamondPerMinute(val float64) {
	v.DiamondPerMinute = val
	v.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbShortVideo, ShortVideoDiamondPerMinute, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetCategoryId(val int) {
	v.CategoryId = val
	v.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbShortVideo, ShortVideoCategoryId, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetSource(val uint8) {
	v.Source = val
	v.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbShortVideo, ShortVideoSource, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetAuthorId(val uint64) {
	v.AuthorId = val
	v.touchUpdatedAt()
	syndb.AddDataToQuickChan(TbShortVideo, ShortVideoAuthorId, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetCreatedAt(val time.Time) {
	v.CreatedAt = val
	syndb.AddDataToQuickChan(TbShortVideo, db.CreatedAtName, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetUpdatedAt(val time.Time) {
	v.UpdatedAt = val
	syndb.AddDataToQuickChan(TbShortVideo, db.UpdatedAtName, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) touchUpdatedAt() {
	v.UpdatedAt = time.Now()
	syndb.AddDataToQuickChan(TbShortVideo, db.UpdatedAtName, &syndb.ColData{
		IdVal: v.ID, ColVal: v.UpdatedAt,
	})
}

func initShortVideo() {
	syndb.RegQuickWithMiddle(TbShortVideo, db.CreatedAtName)
	syndb.RegQuickWithMiddle(TbShortVideo, db.UpdatedAtName)
	syndb.RegQuickWithMiddle(TbShortVideo, ShortVideoTitle)
	syndb.RegQuickWithMiddle(TbShortVideo, ShortVideoVideo)
	syndb.RegQuickWithMiddle(TbShortVideo, ShortVideoCover)
	syndb.RegQuickWithMiddle(TbShortVideo, ShortVideoSort)
	syndb.RegQuickWithMiddle(TbShortVideo, ShortVideoStatusCol)
	syndb.RegQuickWithMiddle(TbShortVideo, ShortVideoIsPaid)
	syndb.RegQuickWithMiddle(TbShortVideo, ShortVideoDiamondPerMinute)
	syndb.RegQuickWithMiddle(TbShortVideo, ShortVideoCategoryId)
	syndb.RegQuickWithMiddle(TbShortVideo, ShortVideoSource)
	syndb.RegQuickWithMiddle(TbShortVideo, ShortVideoAuthorId)
	migrate.AutoMigrate(&ShortVideo{})
}

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
	ShortVideoPayDiamond       db.TbCol = "pay_diamond"
	ShortVideoCategoryId       db.TbCol = "category_id"
	ShortVideoSource           db.TbCol = "source"
	ShortVideoAuthorId         db.TbCol = "author_id"
	ShortVideoAuthorType       db.TbCol = "author_type"
	ShortVideoDuration         db.TbCol = "duration"
	ShortVideoFreeWatchSeconds db.TbCol = "free_watch_seconds"
)

const ShortVideoDefaultFreeWatchSeconds uint32 = 15

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

// 作者类型
const (
	ShortVideoAuthorTypeApp uint8 = 0 // App用户
	ShortVideoAuthorTypeCMS uint8 = 1 // CMS用户
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
	PayDiamond       float64 `gorm:"type:decimal(10,4);default:0;comment:付费钻石(一次性,付费时有效)" json:"payDiamond"`
	CategoryId       int     `gorm:"default:0;comment:视频分类ID" json:"categoryId"`
	Source           uint8   `gorm:"default:1;comment:视频来源(1原创,2转发,3AI生成)" json:"source"`
	AuthorId         uint64  `gorm:"default:0;comment:作者用户ID" json:"authorId"`
	AuthorType       uint8   `gorm:"default:0;comment:作者类型(0App用户,1CMS用户)" json:"authorType"`
	Duration         uint32  `gorm:"default:0;comment:视频时长(秒)" json:"duration"`
	FreeWatchSeconds uint32  `gorm:"default:15;comment:免费观看时长(秒)" json:"freeWatchSeconds"`
}

// NewShortVideo 构造内存对象,字段通过 syndb 异步入库
func NewShortVideo(id uint64, title, video, cover string, sort int, isPaid uint8, payDiamond float64, categoryId int, source uint8, authorId uint64, authorType uint8, duration, freeWatchSeconds uint32) *ShortVideo {
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
	v.SetPayDiamond(payDiamond)
	v.SetCategoryId(categoryId)
	v.SetSource(source)
	v.SetAuthorId(authorId)
	v.initAuthorType(authorType)
	v.SetDuration(duration)
	v.SetFreeWatchSeconds(freeWatchSeconds)
	return v
}

func (v *ShortVideo) SetTitle(val string) {
	v.Title = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoTitle, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetVideo(val string) {
	v.Video = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoVideo, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetCover(val string) {
	v.Cover = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoCover, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetSort(val int) {
	v.Sort = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoSort, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetStatus(val uint8) {
	v.Status = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoStatusCol, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetIsPaid(val uint8) {
	v.IsPaid = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoIsPaid, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetPayDiamond(val float64) {
	v.PayDiamond = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoPayDiamond, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetCategoryId(val int) {
	v.CategoryId = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoCategoryId, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetSource(val uint8) {
	v.Source = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoSource, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetAuthorId(val uint64) {
	v.AuthorId = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoAuthorId, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

// initAuthorType 仅在创建时写入,创建后不可修改
func (v *ShortVideo) initAuthorType(val uint8) {
	v.AuthorType = val
	syndb.AddData(TbShortVideo, ShortVideoAuthorType, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetDuration(val uint32) {
	v.Duration = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoDuration, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetFreeWatchSeconds(val uint32) {
	v.FreeWatchSeconds = val
	v.touchUpdatedAt()
	syndb.AddData(TbShortVideo, ShortVideoFreeWatchSeconds, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetCreatedAt(val time.Time) {
	v.CreatedAt = val
	syndb.AddData(TbShortVideo, db.CreatedAtName, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) SetUpdatedAt(val time.Time) {
	v.UpdatedAt = val
	syndb.AddData(TbShortVideo, db.UpdatedAtName, &syndb.ColData{
		IdVal: v.ID, ColVal: val,
	})
}

func (v *ShortVideo) touchUpdatedAt() {
	v.UpdatedAt = time.Now()
	syndb.AddData(TbShortVideo, db.UpdatedAtName, &syndb.ColData{
		IdVal: v.ID, ColVal: v.UpdatedAt,
	})
}

func initShortVideo() {
	syndb.RegQuick(TbShortVideo, db.CreatedAtName)
	syndb.RegQuick(TbShortVideo, db.UpdatedAtName)
	syndb.RegQuick(TbShortVideo, ShortVideoTitle)
	syndb.RegQuick(TbShortVideo, ShortVideoVideo)
	syndb.RegQuick(TbShortVideo, ShortVideoCover)
	syndb.RegQuick(TbShortVideo, ShortVideoSort)
	syndb.RegQuick(TbShortVideo, ShortVideoStatusCol)
	syndb.RegQuick(TbShortVideo, ShortVideoIsPaid)
	syndb.RegQuick(TbShortVideo, ShortVideoPayDiamond)
	syndb.RegQuick(TbShortVideo, ShortVideoCategoryId)
	syndb.RegQuick(TbShortVideo, ShortVideoSource)
	syndb.RegQuick(TbShortVideo, ShortVideoAuthorId)
	syndb.RegQuick(TbShortVideo, ShortVideoAuthorType)
	syndb.RegQuick(TbShortVideo, ShortVideoDuration)
	syndb.RegQuick(TbShortVideo, ShortVideoFreeWatchSeconds)
	migrate.AutoMigrate(&ShortVideo{})
}

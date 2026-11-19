package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbAnchorVideo db.TbName = "anchor_videos"
)

const (
	AnchorVideoAnchorId    db.TbCol = "anchor_id"
	AnchorVideoVideo       db.TbCol = "video"
	AnchorVideoCover       db.TbCol = "cover"
	AnchorVideoFileType    db.TbCol = "file_type"
	AnchorVideoDuration    db.TbCol = "duration"
	AnchorVideoPrice       db.TbCol = "price"
	AnchorVideoDescription db.TbCol = "description"
)

// AnchorVideo 主播相册视频(软删走 syndb)
type AnchorVideo struct {
	migrate.MoreModel
	AnchorId    uint64  `gorm:"index:idx_av_anchor_created,priority:1;default:0;comment:主播用户ID" json:"anchorId"`
	Video       string  `gorm:"size:255;default:'';comment:视频资源名" json:"video"`
	Cover       string  `gorm:"size:512;default:'';comment:封面(原样存储)" json:"cover"`
	FileType    string  `gorm:"size:128;default:'';comment:文件类型(原样存储)" json:"fileType"`
	Duration    uint32  `gorm:"default:0;comment:视频秒数" json:"duration"`
	Price       float64 `gorm:"type:decimal(10,4);default:0;comment:钻石价格" json:"price"`
	Description string  `gorm:"size:200;default:'';comment:描述(可空)" json:"description"`
	// CreatedAt 覆盖嵌入字段以建复合索引
	CreatedAt time.Time `gorm:"index:idx_av_anchor_created,priority:2" json:"createdAt"`
}

func NewAnchorVideo(anchorId uint64, video string, price float64, description, cover, fileType string, duration uint32) *AnchorVideo {
	ret := &AnchorVideo{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetAnchorId(anchorId)
	ret.SetVideo(video)
	ret.SetCover(cover)
	ret.SetFileType(fileType)
	ret.SetDuration(duration)
	ret.SetPrice(price)
	ret.SetDescription(description)
	return ret
}

func (r *AnchorVideo) SetAnchorId(v uint64) {
	r.AnchorId = v
	syndb.AddData(TbAnchorVideo, AnchorVideoAnchorId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) SetVideo(v string) {
	r.Video = v
	syndb.AddData(TbAnchorVideo, AnchorVideoVideo, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) SetCover(v string) {
	r.Cover = v
	syndb.AddData(TbAnchorVideo, AnchorVideoCover, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) SetFileType(v string) {
	r.FileType = v
	syndb.AddData(TbAnchorVideo, AnchorVideoFileType, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) SetDuration(v uint32) {
	r.Duration = v
	syndb.AddData(TbAnchorVideo, AnchorVideoDuration, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) SetPrice(v float64) {
	r.Price = v
	syndb.AddData(TbAnchorVideo, AnchorVideoPrice, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) SetDescription(v string) {
	r.Description = v
	syndb.AddData(TbAnchorVideo, AnchorVideoDescription, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbAnchorVideo, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbAnchorVideo, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) SetIsDeleted(v bool) {
	r.IsDeleted = v
	syndb.AddData(TbAnchorVideo, db.IsDeletedName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) SetDeletedAt(v time.Time) {
	r.DeletedAt = v
	syndb.AddData(TbAnchorVideo, db.DeletedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideo) MarkDeleted() {
	now := time.Now()
	r.SetIsDeleted(true)
	r.SetDeletedAt(now)
	r.SetUpdatedAt(now)
}

func initAnchorVideo() {
	syndb.RegQuick(TbAnchorVideo, db.CreatedAtName)
	syndb.RegQuick(TbAnchorVideo, db.UpdatedAtName)
	syndb.RegQuick(TbAnchorVideo, db.IsDeletedName)
	syndb.RegQuick(TbAnchorVideo, db.DeletedAtName)
	syndb.RegQuick(TbAnchorVideo, AnchorVideoAnchorId)
	syndb.RegQuick(TbAnchorVideo, AnchorVideoVideo)
	syndb.RegQuick(TbAnchorVideo, AnchorVideoCover)
	syndb.RegQuick(TbAnchorVideo, AnchorVideoFileType)
	syndb.RegQuick(TbAnchorVideo, AnchorVideoDuration)
	syndb.RegQuick(TbAnchorVideo, AnchorVideoPrice)
	syndb.RegQuick(TbAnchorVideo, AnchorVideoDescription)
	migrate.AutoMigrate(&AnchorVideo{})
}

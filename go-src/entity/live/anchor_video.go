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
	AnchorVideoAnchorId db.TbCol = "anchor_id"
	AnchorVideoVideo    db.TbCol = "video"
)

// AnchorVideo 主播相册视频(软删走 syndb)
type AnchorVideo struct {
	migrate.MoreModel
	AnchorId uint64 `gorm:"index:idx_av_anchor_created,priority:1;default:0;comment:主播用户ID" json:"anchorId"`
	Video    string `gorm:"size:255;default:'';comment:视频资源名" json:"video"`
	// CreatedAt 覆盖嵌入字段以建复合索引
	CreatedAt time.Time `gorm:"index:idx_av_anchor_created,priority:2" json:"createdAt"`
}

func NewAnchorVideo(anchorId uint64, video string) *AnchorVideo {
	ret := &AnchorVideo{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetAnchorId(anchorId)
	ret.SetVideo(video)
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
	migrate.AutoMigrate(&AnchorVideo{})
}

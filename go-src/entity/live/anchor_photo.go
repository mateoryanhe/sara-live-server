package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbAnchorPhoto db.TbName = "anchor_photos"
)

const (
	AnchorPhotoAnchorId db.TbCol = "anchor_id"
	AnchorPhotoPhoto    db.TbCol = "photo"
)

// AnchorPhoto 主播相册照片(最多100张/人;软删走 syndb)
type AnchorPhoto struct {
	migrate.MoreModel
	AnchorId uint64 `gorm:"index:idx_ap_anchor_created,priority:1;default:0;comment:主播用户ID" json:"anchorId"`
	Photo    string `gorm:"size:255;default:'';comment:照片资源名" json:"photo"`
	// CreatedAt 覆盖嵌入字段以建复合索引
	CreatedAt time.Time `gorm:"index:idx_ap_anchor_created,priority:2" json:"createdAt"`
}

func NewAnchorPhoto(anchorId uint64, photo string) *AnchorPhoto {
	ret := &AnchorPhoto{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetAnchorId(anchorId)
	ret.SetPhoto(photo)
	return ret
}

func (r *AnchorPhoto) SetAnchorId(v uint64) {
	r.AnchorId = v
	syndb.AddData(TbAnchorPhoto, AnchorPhotoAnchorId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorPhoto) SetPhoto(v string) {
	r.Photo = v
	syndb.AddData(TbAnchorPhoto, AnchorPhotoPhoto, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorPhoto) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbAnchorPhoto, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorPhoto) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbAnchorPhoto, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorPhoto) SetIsDeleted(v bool) {
	r.IsDeleted = v
	syndb.AddData(TbAnchorPhoto, db.IsDeletedName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorPhoto) SetDeletedAt(v time.Time) {
	r.DeletedAt = v
	syndb.AddData(TbAnchorPhoto, db.DeletedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorPhoto) MarkDeleted() {
	now := time.Now()
	r.SetIsDeleted(true)
	r.SetDeletedAt(now)
	r.SetUpdatedAt(now)
}

func initAnchorPhoto() {
	syndb.RegQuick(TbAnchorPhoto, db.CreatedAtName)
	syndb.RegQuick(TbAnchorPhoto, db.UpdatedAtName)
	syndb.RegQuick(TbAnchorPhoto, db.IsDeletedName)
	syndb.RegQuick(TbAnchorPhoto, db.DeletedAtName)
	syndb.RegQuick(TbAnchorPhoto, AnchorPhotoAnchorId)
	syndb.RegQuick(TbAnchorPhoto, AnchorPhotoPhoto)
	migrate.AutoMigrate(&AnchorPhoto{})
}

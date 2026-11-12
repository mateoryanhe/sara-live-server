package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbAnchorVideoPurchase db.TbName = "anchor_video_purchases"
)

const (
	AnchorVideoPurchaseUserId  db.TbCol = "user_id"
	AnchorVideoPurchaseVideoId db.TbCol = "video_id"
	AnchorVideoPurchaseAnchorId db.TbCol = "anchor_id"
	AnchorVideoPurchasePrice   db.TbCol = "price"
	AnchorVideoPurchasePaidAt  db.TbCol = "paid_at"
)

// AnchorVideoPurchase 主播相册视频付费解锁记录(用户+视频维度,一次付费长期有效)
type AnchorVideoPurchase struct {
	ID        string    `gorm:"primaryKey;size:64;comment:复合ID(userId_videoId)" json:"id"`
	UserId    uint64    `gorm:"index:idx_avp_user_paid,priority:1;default:0;comment:付费用户ID" json:"userId"`
	VideoId   uint64    `gorm:"index;default:0;comment:相册视频ID" json:"videoId"`
	AnchorId  uint64    `gorm:"index;default:0;comment:主播用户ID" json:"anchorId"`
	Price     float64   `gorm:"type:decimal(10,4);default:0;comment:实付钻石" json:"price"`
	PaidAt    time.Time `gorm:"index:idx_avp_user_paid,priority:2;comment:付费时间" json:"paidAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewAnchorVideoPurchase(userId, videoId, anchorId uint64, price float64, paidAt time.Time) *AnchorVideoPurchase {
	row := &AnchorVideoPurchase{}
	row.ID = fmt.Sprintf("%d_%d", userId, videoId)
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	row.SetVideoId(videoId)
	row.SetAnchorId(anchorId)
	row.SetPrice(price)
	row.SetPaidAt(paidAt)
	return row
}

func (r *AnchorVideoPurchase) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbAnchorVideoPurchase, AnchorVideoPurchaseUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideoPurchase) SetVideoId(v uint64) {
	r.VideoId = v
	syndb.AddData(TbAnchorVideoPurchase, AnchorVideoPurchaseVideoId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideoPurchase) SetAnchorId(v uint64) {
	r.AnchorId = v
	syndb.AddData(TbAnchorVideoPurchase, AnchorVideoPurchaseAnchorId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideoPurchase) SetPrice(v float64) {
	r.Price = v
	syndb.AddData(TbAnchorVideoPurchase, AnchorVideoPurchasePrice, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideoPurchase) SetPaidAt(v time.Time) {
	r.PaidAt = v
	syndb.AddData(TbAnchorVideoPurchase, AnchorVideoPurchasePaidAt, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideoPurchase) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbAnchorVideoPurchase, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *AnchorVideoPurchase) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbAnchorVideoPurchase, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initAnchorVideoPurchase() {
	syndb.RegQuick(TbAnchorVideoPurchase, db.CreatedAtName)
	syndb.RegQuick(TbAnchorVideoPurchase, db.UpdatedAtName)
	syndb.RegQuick(TbAnchorVideoPurchase, AnchorVideoPurchaseUserId)
	syndb.RegQuick(TbAnchorVideoPurchase, AnchorVideoPurchaseVideoId)
	syndb.RegQuick(TbAnchorVideoPurchase, AnchorVideoPurchaseAnchorId)
	syndb.RegQuick(TbAnchorVideoPurchase, AnchorVideoPurchasePrice)
	syndb.RegQuick(TbAnchorVideoPurchase, AnchorVideoPurchasePaidAt)
	migrate.AutoMigrate(&AnchorVideoPurchase{})
}

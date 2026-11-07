package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbAnchorReceivedGift db.TbName = "anchor_received_gifts"
)

const (
	AnchorReceivedGiftAnchorId db.TbCol = "anchor_id"
	AnchorReceivedGiftGiftId   db.TbCol = "gift_id"
	AnchorReceivedGiftCount    db.TbCol = "count"
)

// AnchorReceivedGift 主播收到过的礼物(按礼物ID去重一行,累计数量;主键 ID = anchorId_giftId)
type AnchorReceivedGift struct {
	ID        string    `gorm:"primaryKey;size:64;comment:复合ID(anchorId_giftId)" json:"id"`
	AnchorId  uint64    `gorm:"uniqueIndex:uk_anchor_gift,priority:1;index:idx_arg_anchor_created,priority:1;default:0;comment:主播用户ID" json:"anchorId"`
	GiftId    uint64    `gorm:"uniqueIndex:uk_anchor_gift,priority:2;default:0;comment:礼物ID" json:"giftId"`
	Count     uint64    `gorm:"default:0;comment:累计收到数量" json:"count"`
	CreatedAt time.Time `gorm:"index:idx_arg_anchor_created,priority:2" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func BuildAnchorReceivedGiftId(anchorId, giftId uint64) string {
	return fmt.Sprintf("%d_%d", anchorId, giftId)
}

// NewAnchorReceivedGift 构造主播收礼记录,字段通过 syndb 异步入库
func NewAnchorReceivedGift(anchorId, giftId, count uint64) *AnchorReceivedGift {
	if count == 0 {
		count = 1
	}
	r := &AnchorReceivedGift{}
	r.ID = BuildAnchorReceivedGiftId(anchorId, giftId)
	now := time.Now()
	r.SetCreatedAt(now)
	r.SetUpdatedAt(now)
	r.SetAnchorId(anchorId)
	r.SetGiftId(giftId)
	r.SetCount(count)
	return r
}

func (r *AnchorReceivedGift) SetAnchorId(v uint64) {
	r.AnchorId = v
	syndb.AddData(TbAnchorReceivedGift, AnchorReceivedGiftAnchorId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *AnchorReceivedGift) SetGiftId(v uint64) {
	r.GiftId = v
	syndb.AddData(TbAnchorReceivedGift, AnchorReceivedGiftGiftId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *AnchorReceivedGift) SetCount(v uint64) {
	r.Count = v
	r.UpdatedAt = time.Now()
	syndb.AddData(TbAnchorReceivedGift, AnchorReceivedGiftCount, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
	syndb.AddData(TbAnchorReceivedGift, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.UpdatedAt,
	})
}

func (r *AnchorReceivedGift) AddCount(delta uint64) {
	if delta == 0 {
		return
	}
	r.SetCount(r.Count + delta)
}

func (r *AnchorReceivedGift) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbAnchorReceivedGift, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *AnchorReceivedGift) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbAnchorReceivedGift, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initAnchorReceivedGift() {
	syndb.RegQuick(TbAnchorReceivedGift, db.CreatedAtName)
	syndb.RegQuick(TbAnchorReceivedGift, db.UpdatedAtName)
	syndb.RegQuick(TbAnchorReceivedGift, AnchorReceivedGiftAnchorId)
	syndb.RegQuick(TbAnchorReceivedGift, AnchorReceivedGiftGiftId)
	syndb.RegQuick(TbAnchorReceivedGift, AnchorReceivedGiftCount)
	migrate.AutoMigrate(&AnchorReceivedGift{})
}

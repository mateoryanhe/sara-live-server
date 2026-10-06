package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbOneToOneRoom db.TbName = "one_to_one_rooms"
)

const (
	OneToOneRoomStatus   db.TbCol = "status"
	OneToOneRoomBilling  db.TbCol = "billing"
	OneToOneRoomTitle db.TbCol = "title"
	OneToOneRoomCover db.TbCol = "cover"
	OneToOneRoomTagId db.TbCol = "tag_id"
)

const (
	OneToOneRoomStatusOffShelf uint8 = 0
	OneToOneRoomStatusOnShelf  uint8 = 1
)

// OneToOneRoom 1v1 房间展示与计费(主键=主播用户ID);与直播间资料独立维护。
type OneToOneRoom struct {
	migrate.OneModel
	UpdatedAt time.Time `json:"-"`
	Status    uint8     `gorm:"index;default:1;comment:状态(0-下架,1-上架)" json:"status"`
	Billing   float64   `gorm:"type:decimal(10,4);default:0;comment:1v1通话每分钟钻石" json:"billing"`
	Title     string    `gorm:"size:128;default:'';comment:1v1房间标题" json:"title"`
	Cover string `gorm:"size:255;default:'';comment:1v1封面(对象名)" json:"cover"`
	TagId uint64 `gorm:"default:0;comment:标签ID" json:"tagId"`
}

func NewOneToOneRoom(anchorId uint64) *OneToOneRoom {
	r := &OneToOneRoom{}
	r.ID = anchorId
	now := time.Now()
	r.SetCreatedAt(now)
	r.SetUpdatedAt(now)
	r.SetStatus(OneToOneRoomStatusOnShelf)
	return r
}

func NormalizeOneToOneBilling(v float64) float64 {
	if v < 0 {
		return 0
	}
	return v
}

func (r *OneToOneRoom) SetBilling(v float64) {
	v = NormalizeOneToOneBilling(v)
	r.Billing = v
	r.touchUpdatedAt()
	syndb.AddData(TbOneToOneRoom, OneToOneRoomBilling, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *OneToOneRoom) SetStatus(v uint8) {
	if v != OneToOneRoomStatusOffShelf && v != OneToOneRoomStatusOnShelf {
		v = OneToOneRoomStatusOnShelf
	}
	r.Status = v
	r.touchUpdatedAt()
	syndb.AddData(TbOneToOneRoom, OneToOneRoomStatus, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *OneToOneRoom) SetTitle(v string) {
	r.Title = v
	r.touchUpdatedAt()
	syndb.AddData(TbOneToOneRoom, OneToOneRoomTitle, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *OneToOneRoom) SetCover(v string) {
	r.Cover = v
	r.touchUpdatedAt()
	syndb.AddData(TbOneToOneRoom, OneToOneRoomCover, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *OneToOneRoom) SetTagId(v uint64) {
	r.TagId = v
	r.touchUpdatedAt()
	syndb.AddData(TbOneToOneRoom, OneToOneRoomTagId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *OneToOneRoom) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbOneToOneRoom, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *OneToOneRoom) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbOneToOneRoom, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *OneToOneRoom) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	syndb.AddData(TbOneToOneRoom, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.UpdatedAt,
	})
}

func (r *OneToOneRoom) IsOnShelf() bool {
	return r != nil && r.Status == OneToOneRoomStatusOnShelf
}

func initOneToOneRoom() {
	syndb.RegQuick(TbOneToOneRoom, db.CreatedAtName)
	syndb.RegQuick(TbOneToOneRoom, db.UpdatedAtName)
	syndb.RegQuick(TbOneToOneRoom, OneToOneRoomStatus)
	syndb.RegQuick(TbOneToOneRoom, OneToOneRoomBilling)
	syndb.RegQuick(TbOneToOneRoom, OneToOneRoomTitle)
	syndb.RegQuick(TbOneToOneRoom, OneToOneRoomCover)
	syndb.RegQuick(TbOneToOneRoom, OneToOneRoomTagId)
	migrate.AutoMigrate(&OneToOneRoom{})
}

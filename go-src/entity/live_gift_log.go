package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbLiveGiftLog db.TbName = "live_gift_logs"
)

const (
	LiveGiftLogRoomId       db.TbCol = "room_id"
	LiveGiftLogLiveRecordId db.TbCol = "live_record_id"
	LiveGiftLogSenderId     db.TbCol = "sender_id"
	LiveGiftLogReceiverId   db.TbCol = "receiver_id"
	LiveGiftLogGiftId       db.TbCol = "gift_id"
	LiveGiftLogCount        db.TbCol = "count"
	LiveGiftLogUnitPrice    db.TbCol = "unit_price"
	LiveGiftLogTotalCost    db.TbCol = "total_cost"
)

// LiveGiftLog 直播礼物流水
type LiveGiftLog struct {
	migrate.OneModel
	RoomId       uint64 `gorm:"index;default:0;comment:直播间ID" json:"roomId"`
	LiveRecordId uint64 `gorm:"index;default:0;comment:直播记录ID" json:"liveRecordId"`
	SenderId     uint64 `gorm:"index;default:0;comment:送礼用户ID" json:"senderId"`
	ReceiverId   uint64 `gorm:"index;default:0;comment:接收者ID(主播)" json:"receiverId"`
	GiftId       uint64 `gorm:"index;default:0;comment:礼物ID" json:"giftId"`
	Count        int    `gorm:"default:0;comment:赠送数量" json:"count"`
	UnitPrice    uint64 `gorm:"default:0;comment:礼物单价(钻石)" json:"unitPrice"`
	TotalCost    uint64 `gorm:"default:0;comment:总消耗钻石" json:"totalCost"`
}

func NewLiveGiftLog(id uint64) *LiveGiftLog {
	ret := &LiveGiftLog{}
	ret.ID = id
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

func NewLiveGiftLogRecord(roomId, liveRecordId, senderId, receiverId, giftId uint64, count int, unitPrice, totalCost uint64) *LiveGiftLog {
	ret := NewLiveGiftLog(snowflake.GetId())
	ret.SetRoomId(roomId)
	ret.SetLiveRecordId(liveRecordId)
	ret.SetSenderId(senderId)
	ret.SetReceiverId(receiverId)
	ret.SetGiftId(giftId)
	ret.SetCount(count)
	ret.SetUnitPrice(unitPrice)
	ret.SetTotalCost(totalCost)
	return ret
}

func (r *LiveGiftLog) SetRoomId(v uint64) {
	r.RoomId = v
	syndb.AddDataToQuickChan(TbLiveGiftLog, LiveGiftLogRoomId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveGiftLog) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	syndb.AddDataToQuickChan(TbLiveGiftLog, LiveGiftLogLiveRecordId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveGiftLog) SetSenderId(v uint64) {
	r.SenderId = v
	syndb.AddDataToQuickChan(TbLiveGiftLog, LiveGiftLogSenderId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveGiftLog) SetReceiverId(v uint64) {
	r.ReceiverId = v
	syndb.AddDataToQuickChan(TbLiveGiftLog, LiveGiftLogReceiverId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveGiftLog) SetGiftId(v uint64) {
	r.GiftId = v
	syndb.AddDataToQuickChan(TbLiveGiftLog, LiveGiftLogGiftId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveGiftLog) SetCount(v int) {
	r.Count = v
	syndb.AddDataToQuickChan(TbLiveGiftLog, LiveGiftLogCount, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveGiftLog) SetUnitPrice(v uint64) {
	r.UnitPrice = v
	syndb.AddDataToQuickChan(TbLiveGiftLog, LiveGiftLogUnitPrice, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveGiftLog) SetTotalCost(v uint64) {
	r.TotalCost = v
	syndb.AddDataToQuickChan(TbLiveGiftLog, LiveGiftLogTotalCost, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveGiftLog) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddDataToQuickChan(TbLiveGiftLog, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveGiftLog) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddDataToQuickChan(TbLiveGiftLog, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initLiveGiftLog() {
	syndb.RegQuickWithLarge(TbLiveGiftLog, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbLiveGiftLog, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbLiveGiftLog, LiveGiftLogRoomId)
	syndb.RegQuickWithLarge(TbLiveGiftLog, LiveGiftLogLiveRecordId)
	syndb.RegQuickWithLarge(TbLiveGiftLog, LiveGiftLogSenderId)
	syndb.RegQuickWithLarge(TbLiveGiftLog, LiveGiftLogReceiverId)
	syndb.RegQuickWithLarge(TbLiveGiftLog, LiveGiftLogGiftId)
	syndb.RegQuickWithLarge(TbLiveGiftLog, LiveGiftLogCount)
	syndb.RegQuickWithLarge(TbLiveGiftLog, LiveGiftLogUnitPrice)
	syndb.RegQuickWithLarge(TbLiveGiftLog, LiveGiftLogTotalCost)
	migrate.AutoMigrate(&LiveGiftLog{})
}

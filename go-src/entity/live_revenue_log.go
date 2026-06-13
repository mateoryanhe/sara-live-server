package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/constants/liverevenue"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRevenueLog db.TbName = "live_revenue_logs"
)

const (
	LiveRevenueLogRevenueType  db.TbCol = "revenue_type"
	LiveRevenueLogRoomId       db.TbCol = "room_id"
	LiveRevenueLogLiveRecordId db.TbCol = "live_record_id"
	LiveRevenueLogSenderId     db.TbCol = "sender_id"
	LiveRevenueLogReceiverId   db.TbCol = "receiver_id"
	LiveRevenueLogBizId        db.TbCol = "biz_id"
	LiveRevenueLogCount        db.TbCol = "count"
	LiveRevenueLogUnitPrice    db.TbCol = "unit_price"
	LiveRevenueLogTotalAmount  db.TbCol = "total_amount"
)

// LiveRevenueLog 直播收益流水(礼物/付费弹幕/游戏下注等)
type LiveRevenueLog struct {
	migrate.OneModel
	RevenueType  uint8   `gorm:"index;default:1;comment:收益类型(1礼物,2付费弹幕,3游戏下注)" json:"revenueType"`
	RoomId       uint64  `gorm:"index;default:0;comment:直播间ID" json:"roomId"`
	LiveRecordId uint64  `gorm:"index;default:0;comment:直播记录ID" json:"liveRecordId"`
	SenderId     uint64  `gorm:"index;default:0;comment:付款用户ID" json:"senderId"`
	ReceiverId   uint64  `gorm:"index;default:0;comment:收益用户ID(主播)" json:"receiverId"`
	BizId        uint64  `gorm:"index;default:0;comment:业务关联ID(礼物ID/弹幕ID/游戏ID等)" json:"bizId"`
	Count        int     `gorm:"default:0;comment:数量" json:"count"`
	UnitPrice    float64 `gorm:"default:0;comment:单价(钻石)" json:"unitPrice"`
	TotalAmount  float64 `gorm:"default:0;comment:流水金额(钻石)" json:"totalAmount"`
}

func NewLiveRevenueLog(id uint64) *LiveRevenueLog {
	ret := &LiveRevenueLog{}
	ret.ID = id
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

func NewLiveRevenueLogRecord(roomId, liveRecordId, senderId, receiverId, bizId uint64, count int, unitPrice, totalAmount float64, revenueType ...uint8) *LiveRevenueLog {
	ret := NewLiveRevenueLog(snowflake.GetId())
	rt := uint8(liverevenue.Gift)
	if len(revenueType) > 0 && liverevenue.IsValid(liverevenue.Type(revenueType[0])) {
		rt = revenueType[0]
	}
	ret.SetRevenueType(rt)
	ret.SetRoomId(roomId)
	ret.SetLiveRecordId(liveRecordId)
	ret.SetSenderId(senderId)
	ret.SetReceiverId(receiverId)
	ret.SetBizId(bizId)
	ret.SetCount(count)
	ret.SetUnitPrice(unitPrice)
	ret.SetTotalAmount(totalAmount)
	return ret
}

func (r *LiveRevenueLog) SetRevenueType(v uint8) {
	r.RevenueType = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, LiveRevenueLogRevenueType, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetRoomId(v uint64) {
	r.RoomId = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, LiveRevenueLogRoomId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, LiveRevenueLogLiveRecordId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetSenderId(v uint64) {
	r.SenderId = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, LiveRevenueLogSenderId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetReceiverId(v uint64) {
	r.ReceiverId = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, LiveRevenueLogReceiverId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetBizId(v uint64) {
	r.BizId = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, LiveRevenueLogBizId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetCount(v int) {
	r.Count = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, LiveRevenueLogCount, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetUnitPrice(v float64) {
	r.UnitPrice = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, LiveRevenueLogUnitPrice, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetTotalAmount(v float64) {
	r.TotalAmount = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, LiveRevenueLogTotalAmount, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddDataToQuickChan(TbLiveRevenueLog, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initLiveRevenueLog() {
	syndb.RegQuickWithLarge(TbLiveRevenueLog, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbLiveRevenueLog, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbLiveRevenueLog, LiveRevenueLogRevenueType)
	syndb.RegQuickWithLarge(TbLiveRevenueLog, LiveRevenueLogRoomId)
	syndb.RegQuickWithLarge(TbLiveRevenueLog, LiveRevenueLogLiveRecordId)
	syndb.RegQuickWithLarge(TbLiveRevenueLog, LiveRevenueLogSenderId)
	syndb.RegQuickWithLarge(TbLiveRevenueLog, LiveRevenueLogReceiverId)
	syndb.RegQuickWithLarge(TbLiveRevenueLog, LiveRevenueLogBizId)
	syndb.RegQuickWithLarge(TbLiveRevenueLog, LiveRevenueLogCount)
	syndb.RegQuickWithLarge(TbLiveRevenueLog, LiveRevenueLogUnitPrice)
	syndb.RegQuickWithLarge(TbLiveRevenueLog, LiveRevenueLogTotalAmount)
	migrate.AutoMigrate(&LiveRevenueLog{})
}

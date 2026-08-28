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
	LiveRevenueLogStatusNormal   uint8 = 0 // 正常
	LiveRevenueLogStatusRefunded uint8 = 1 // 已退款
)

const (
	LiveRevenueLogRevenueType  db.TbCol = "revenue_type"
	LiveRevenueLogRoomId       db.TbCol = "room_id"
	LiveRevenueLogLiveRecordId db.TbCol = "live_record_id"
	LiveRevenueLogSenderId     db.TbCol = "sender_id"
	LiveRevenueLogBizId        db.TbCol = "biz_id"
	LiveRevenueLogCount        db.TbCol = "count"
	LiveRevenueLogUnitPrice    db.TbCol = "unit_price"
	LiveRevenueLogTotalAmount  db.TbCol = "total_amount"
	LiveRevenueLogStatus       db.TbCol = "status"
)

// LiveRevenueLog 直播间社交流水(礼物/付费弹幕/私密房/门票/视频通话等钻石消费)
type LiveRevenueLog struct {
	migrate.OneModel
	RevenueType  uint8   `gorm:"index;default:1;comment:流水类型(1礼物,2付费弹幕,4私密房计费,5门票,6视频通话门票,7视频通话计费)" json:"revenueType"`
	RoomId       uint64  `gorm:"index;default:0;comment:直播间ID(主播用户ID)" json:"roomId"`
	LiveRecordId uint64  `gorm:"index;default:0;comment:直播记录ID" json:"liveRecordId"`
	SenderId     uint64  `gorm:"index;default:0;comment:付款用户ID" json:"senderId"`
	BizId        uint64  `gorm:"index;default:0;comment:业务关联ID(礼物ID/通话订单ID等)" json:"bizId"`
	Count        int     `gorm:"default:0;comment:数量(礼物件数,其余类型多为0或1)" json:"count"`
	UnitPrice    float64 `gorm:"default:0;comment:单价(钻石,礼物场景有效)" json:"unitPrice"`
	TotalAmount  float64 `gorm:"default:0;comment:流水金额(钻石)" json:"totalAmount"`
	Status       uint8   `gorm:"index;default:0;comment:状态(0正常,1已退款)" json:"status"`
}

func NewLiveRevenueLog(id uint64) *LiveRevenueLog {
	ret := &LiveRevenueLog{}
	ret.ID = id
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

func NewLiveRevenueLogRecord(roomId, liveRecordId, senderId, bizId uint64, count int, unitPrice, totalAmount float64, revenueType ...uint8) *LiveRevenueLog {
	ret := NewLiveRevenueLog(snowflake.GetId())
	rt := uint8(liverevenue.Gift)
	if len(revenueType) > 0 && liverevenue.IsValid(liverevenue.Type(revenueType[0])) {
		rt = revenueType[0]
	}
	ret.SetRevenueType(rt)
	ret.SetRoomId(roomId)
	ret.SetLiveRecordId(liveRecordId)
	ret.SetSenderId(senderId)
	ret.SetBizId(bizId)
	ret.SetCount(count)
	ret.SetUnitPrice(unitPrice)
	ret.SetTotalAmount(totalAmount)
	ret.SetStatus(LiveRevenueLogStatusNormal)
	return ret
}

func (r *LiveRevenueLog) IsRefunded() bool {
	return r != nil && r.Status == LiveRevenueLogStatusRefunded
}

func (r *LiveRevenueLog) SetRevenueType(v uint8) {
	r.RevenueType = v
	syndb.AddData(TbLiveRevenueLog, LiveRevenueLogRevenueType, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetRoomId(v uint64) {
	r.RoomId = v
	syndb.AddData(TbLiveRevenueLog, LiveRevenueLogRoomId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	syndb.AddData(TbLiveRevenueLog, LiveRevenueLogLiveRecordId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetSenderId(v uint64) {
	r.SenderId = v
	syndb.AddData(TbLiveRevenueLog, LiveRevenueLogSenderId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetBizId(v uint64) {
	r.BizId = v
	syndb.AddData(TbLiveRevenueLog, LiveRevenueLogBizId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetCount(v int) {
	r.Count = v
	syndb.AddData(TbLiveRevenueLog, LiveRevenueLogCount, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetUnitPrice(v float64) {
	r.UnitPrice = v
	syndb.AddData(TbLiveRevenueLog, LiveRevenueLogUnitPrice, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetTotalAmount(v float64) {
	r.TotalAmount = v
	syndb.AddData(TbLiveRevenueLog, LiveRevenueLogTotalAmount, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetStatus(v uint8) {
	r.Status = v
	syndb.AddData(TbLiveRevenueLog, LiveRevenueLogStatus, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbLiveRevenueLog, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRevenueLog) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbLiveRevenueLog, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initLiveRevenueLog() {
	syndb.RegQuick(TbLiveRevenueLog, db.CreatedAtName)
	syndb.RegQuick(TbLiveRevenueLog, db.UpdatedAtName)
	syndb.RegQuick(TbLiveRevenueLog, LiveRevenueLogRevenueType)
	syndb.RegQuick(TbLiveRevenueLog, LiveRevenueLogRoomId)
	syndb.RegQuick(TbLiveRevenueLog, LiveRevenueLogLiveRecordId)
	syndb.RegQuick(TbLiveRevenueLog, LiveRevenueLogSenderId)
	syndb.RegQuick(TbLiveRevenueLog, LiveRevenueLogBizId)
	syndb.RegQuick(TbLiveRevenueLog, LiveRevenueLogCount)
	syndb.RegQuick(TbLiveRevenueLog, LiveRevenueLogUnitPrice)
	syndb.RegQuick(TbLiveRevenueLog, LiveRevenueLogTotalAmount)
	syndb.RegQuick(TbLiveRevenueLog, LiveRevenueLogStatus)
	migrate.AutoMigrate(&LiveRevenueLog{})
}

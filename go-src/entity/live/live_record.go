package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRecord db.TbName = "live_records"
)

const (
	LiveRecordAnchorId                     db.TbCol = "anchor_id"
	LiveRecordStartTime                    db.TbCol = "start_time"
	LiveRecordEndTime                      db.TbCol = "end_time"
	LiveRecordTotalAudience                db.TbCol = "total_audience"
	LiveRecordTotalLiveDuration            db.TbCol = "total_live_duration"
	LiveRecordTotalIncome                  db.TbCol = "total_income"
	LiveRecordTotalGiftIncome              db.TbCol = "total_gift_income"
	LiveRecordTotalPaidDanmakuIncome       db.TbCol = "total_paid_danmaku_income"
	LiveRecordTotalPrivateRoomIncome       db.TbCol = "total_private_room_income"
	LiveRecordTotalPrivateRoomTicketIncome db.TbCol = "total_private_room_ticket_income"
	LiveRecordTotalPrivateRoomWatchIncome  db.TbCol = "total_private_room_watch_income"
	LiveRecordTotalVideoCallIncome         db.TbCol = "total_video_call_income"
	LiveRecordTotalVideoCallTicketIncome   db.TbCol = "total_video_call_ticket_income"
	LiveRecordTotalVideoCallBillingIncome  db.TbCol = "total_video_call_billing_income"
	LiveRecordTotalGameBet                 db.TbCol = "total_game_bet"
	LiveRecordTotalGiftSender              db.TbCol = "total_gift_sender"
	LiveRecordTotalNewFollower             db.TbCol = "total_new_follower"
)

// LiveRecord 单场直播数据记录
type LiveRecord struct {
	migrate.OneModel
	AnchorId                     uint64     `gorm:"index;default:0;comment:主播ID" json:"anchorId"`
	StartTime                    time.Time  `gorm:"comment:直播开始时间" json:"startTime"`
	EndTime                      *time.Time `gorm:"comment:直播结束时间" json:"endTime"`
	TotalAudience                uint64     `gorm:"default:0;comment:累计观众人数" json:"totalAudience"`
	TotalLiveDuration            float64    `gorm:"default:0;comment:累计直播时长(秒)" json:"totalLiveDuration"`
	TotalIncome                  float64    `gorm:"default:0;comment:总收益" json:"totalIncome"`
	TotalGiftIncome              float64    `gorm:"default:0;comment:礼物收入" json:"totalGiftIncome"`
	TotalPaidDanmakuIncome       float64    `gorm:"default:0;comment:付费弹幕收入" json:"totalPaidDanmakuIncome"`
	TotalPrivateRoomIncome       float64    `gorm:"default:0;comment:私密直播间收入" json:"totalPrivateRoomIncome"`
	TotalPrivateRoomTicketIncome float64    `gorm:"type:decimal(10,4);default:0;comment:私密直播间门票收入" json:"totalPrivateRoomTicketIncome"`
	TotalPrivateRoomWatchIncome  float64    `gorm:"type:decimal(10,4);default:0;comment:私密房观看收入" json:"totalPrivateRoomWatchIncome"`
	TotalVideoCallIncome         float64    `gorm:"type:decimal(10,4);default:0;comment:直播间视频通话收入" json:"totalVideoCallIncome"`
	TotalVideoCallTicketIncome   float64    `gorm:"type:decimal(10,4);default:0;comment:直播间视频通话门票收入" json:"totalVideoCallTicketIncome"`
	TotalVideoCallBillingIncome  float64    `gorm:"type:decimal(10,4);default:0;comment:直播间视频通话计费收入" json:"totalVideoCallBillingIncome"`
	TotalGameBet                 float64    `gorm:"default:0;comment:游戏下注总金额" json:"totalGameBet"`
	TotalGiftSender              uint64     `gorm:"default:0;comment:送礼人数(去重)" json:"totalGiftSender"`
	TotalNewFollower             uint64     `gorm:"default:0;comment:新加粉丝数(去重)" json:"totalNewFollower"`
}

// NewLiveRecord 构造一条直播记录,字段写入通过 syndb 异步入库
func NewLiveRecord(id uint64) *LiveRecord {
	ret := &LiveRecord{}
	ret.ID = id
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetStartTime(now)
	return ret
}

func (r *LiveRecord) SetAnchorId(v uint64) {
	r.AnchorId = v
	syndb.AddData(TbLiveRecord, LiveRecordAnchorId, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: v,
	})
}

func (r *LiveRecord) SetStartTime(v time.Time) {
	r.StartTime = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRecord, LiveRecordStartTime, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: v,
	})
}

func (r *LiveRecord) SetEndTime(v *time.Time) {
	r.EndTime = v
	syndb.AddData(TbLiveRecord, LiveRecordEndTime, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.EndTime,
	})
}

func (r *LiveRecord) AddTotalAudience(val uint64) {
	r.TotalAudience = math.Add(r.TotalAudience, val)

	syndb.AddData(TbLiveRecord, LiveRecordTotalAudience, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalAudience,
	})
}

func (r *LiveRecord) AddTotalLiveDuration(v float64) {
	r.TotalLiveDuration = math.AddFloat64(r.TotalLiveDuration, v)
	syndb.AddData(TbLiveRecord, LiveRecordTotalLiveDuration, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalLiveDuration,
	})
}

func (r *LiveRecord) AddTotalIncome(v float64) {
	r.TotalIncome = math.AddFloat64(r.TotalIncome, v)

	syndb.AddData(TbLiveRecord, LiveRecordTotalIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalIncome,
	})
}

func (r *LiveRecord) AddTotalGiftIncome(v float64) {
	r.TotalGiftIncome = math.AddFloat64(r.TotalGiftIncome, v)

	syndb.AddData(TbLiveRecord, LiveRecordTotalGiftIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalGiftIncome,
	})
}

func (r *LiveRecord) AddTotalPaidDanmakuIncome(v float64) {
	r.TotalPaidDanmakuIncome = math.AddFloat64(r.TotalPaidDanmakuIncome, v)

	syndb.AddData(TbLiveRecord, LiveRecordTotalPaidDanmakuIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalPaidDanmakuIncome,
	})
}

func (r *LiveRecord) AddTotalPrivateRoomIncome(v float64) {
	r.TotalPrivateRoomIncome = math.AddFloat64(r.TotalPrivateRoomIncome, v)

	syndb.AddData(TbLiveRecord, LiveRecordTotalPrivateRoomIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalPrivateRoomIncome,
	})
}

func (r *LiveRecord) AddTotalPrivateRoomTicketIncome(v float64) {
	r.TotalPrivateRoomTicketIncome = math.AddFloat64(r.TotalPrivateRoomTicketIncome, v)

	syndb.AddData(TbLiveRecord, LiveRecordTotalPrivateRoomTicketIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalPrivateRoomTicketIncome,
	})
}

func (r *LiveRecord) AddTotalPrivateRoomWatchIncome(v float64) {
	r.TotalPrivateRoomWatchIncome = math.AddFloat64(r.TotalPrivateRoomWatchIncome, v)

	syndb.AddData(TbLiveRecord, LiveRecordTotalPrivateRoomWatchIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalPrivateRoomWatchIncome,
	})
}

func (r *LiveRecord) AddTotalVideoCallIncome(v float64) {
	r.TotalVideoCallIncome = math.AddFloat64(r.TotalVideoCallIncome, v)

	syndb.AddData(TbLiveRecord, LiveRecordTotalVideoCallIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalVideoCallIncome,
	})
}

func (r *LiveRecord) AddTotalVideoCallTicketIncome(v float64) {
	r.TotalVideoCallTicketIncome = math.AddFloat64(r.TotalVideoCallTicketIncome, v)

	syndb.AddData(TbLiveRecord, LiveRecordTotalVideoCallTicketIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalVideoCallTicketIncome,
	})
}

func (r *LiveRecord) AddTotalVideoCallBillingIncome(v float64) {
	r.TotalVideoCallBillingIncome = math.AddFloat64(r.TotalVideoCallBillingIncome, v)

	syndb.AddData(TbLiveRecord, LiveRecordTotalVideoCallBillingIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalVideoCallBillingIncome,
	})
}

func (r *LiveRecord) AddTotalGameBet(v float64) {
	r.TotalGameBet = math.AddFloat64(r.TotalGameBet, v)

	syndb.AddData(TbLiveRecord, LiveRecordTotalGameBet, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalGameBet,
	})
}

func (r *LiveRecord) AddTotalGiftSender(val uint64) {
	r.TotalGiftSender = math.Add(r.TotalGiftSender, val)

	syndb.AddData(TbLiveRecord, LiveRecordTotalGiftSender, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalGiftSender,
	})
}

func (r *LiveRecord) AddTotalNewFollower(val uint64) {
	r.TotalNewFollower = math.Add(r.TotalNewFollower, val)

	syndb.AddData(TbLiveRecord, LiveRecordTotalNewFollower, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalNewFollower,
	})
}

func (r *LiveRecord) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbLiveRecord, db.CreatedAtName, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: v,
	})
}

func (r *LiveRecord) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbLiveRecord, db.UpdatedAtName, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: v,
	})
}

func (r *LiveRecord) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	syndb.AddData(TbLiveRecord, db.UpdatedAtName, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.UpdatedAt,
	})
}

func initLiveRecord() {
	syndb.RegLazy(TbLiveRecord, db.CreatedAtName)
	syndb.RegLazy(TbLiveRecord, db.UpdatedAtName)
	syndb.RegLazy(TbLiveRecord, LiveRecordAnchorId)
	syndb.RegLazy(TbLiveRecord, LiveRecordStartTime)
	syndb.RegLazy(TbLiveRecord, LiveRecordEndTime)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalAudience)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalLiveDuration)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalIncome)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalGiftIncome)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalPaidDanmakuIncome)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalPrivateRoomIncome)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalPrivateRoomTicketIncome)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalPrivateRoomWatchIncome)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalVideoCallIncome)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalVideoCallTicketIncome)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalVideoCallBillingIncome)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalGameBet)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalGiftSender)
	syndb.RegLazy(TbLiveRecord, LiveRecordTotalNewFollower)
	migrate.AutoMigrate(&LiveRecord{})
}

package entity

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gmlock"
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

func liveRecordLockKey(id uint64) string {
	return fmt.Sprintf("live_record:%d", id)
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

func (r *LiveRecord) addFloatLocked(col db.TbCol, cur *float64, v float64) {
	*cur = math.AddFloat64(*cur, v)
	syndb.AddData(TbLiveRecord, col, &syndb.ColData{IdVal: r.ID, ColVal: *cur})
}

func (r *LiveRecord) addUintLocked(col db.TbCol, cur *uint64, v uint64) {
	*cur = math.Add(*cur, v)
	syndb.AddData(TbLiveRecord, col, &syndb.ColData{IdVal: r.ID, ColVal: *cur})
}

func (r *LiveRecord) withLock(fn func()) {
	key := liveRecordLockKey(r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	fn()
}

func (r *LiveRecord) AddTotalAudience(val uint64) {
	r.withLock(func() {
		r.addUintLocked(LiveRecordTotalAudience, &r.TotalAudience, val)
	})
}

func (r *LiveRecord) AddTotalLiveDuration(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalLiveDuration, &r.TotalLiveDuration, v)
	})
}

func (r *LiveRecord) AddTotalIncome(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalIncome, &r.TotalIncome, v)
	})
}

func (r *LiveRecord) AddTotalGiftIncome(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalGiftIncome, &r.TotalGiftIncome, v)
	})
}

func (r *LiveRecord) AddTotalPaidDanmakuIncome(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalPaidDanmakuIncome, &r.TotalPaidDanmakuIncome, v)
	})
}

func (r *LiveRecord) AddTotalPrivateRoomIncome(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalPrivateRoomIncome, &r.TotalPrivateRoomIncome, v)
	})
}

func (r *LiveRecord) AddTotalPrivateRoomTicketIncome(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalPrivateRoomTicketIncome, &r.TotalPrivateRoomTicketIncome, v)
	})
}

func (r *LiveRecord) AddTotalPrivateRoomWatchIncome(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalPrivateRoomWatchIncome, &r.TotalPrivateRoomWatchIncome, v)
	})
}

func (r *LiveRecord) AddTotalVideoCallIncome(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalVideoCallIncome, &r.TotalVideoCallIncome, v)
	})
}

func (r *LiveRecord) AddTotalVideoCallTicketIncome(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalVideoCallTicketIncome, &r.TotalVideoCallTicketIncome, v)
	})
}

func (r *LiveRecord) AddTotalVideoCallBillingIncome(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalVideoCallBillingIncome, &r.TotalVideoCallBillingIncome, v)
	})
}

func (r *LiveRecord) AddTotalGameBet(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalGameBet, &r.TotalGameBet, v)
	})
}

func (r *LiveRecord) AddTotalGiftSender(val uint64) {
	r.withLock(func() {
		r.addUintLocked(LiveRecordTotalGiftSender, &r.TotalGiftSender, val)
	})
}

func (r *LiveRecord) AddTotalNewFollower(val uint64) {
	r.withLock(func() {
		r.addUintLocked(LiveRecordTotalNewFollower, &r.TotalNewFollower, val)
	})
}

// AddGiftEarn 礼物收益(总收益+礼物,内部加锁)
func (r *LiveRecord) AddGiftEarn(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalIncome, &r.TotalIncome, v)
		if v > 0 {
			r.addFloatLocked(LiveRecordTotalGiftIncome, &r.TotalGiftIncome, v)
		}
	})
}

// AddPaidDanmakuEarn 付费弹幕收益(总收益+弹幕,内部加锁)
func (r *LiveRecord) AddPaidDanmakuEarn(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalIncome, &r.TotalIncome, v)
		if v > 0 {
			r.addFloatLocked(LiveRecordTotalPaidDanmakuIncome, &r.TotalPaidDanmakuIncome, v)
		}
	})
}

// AddPrivateRoomTicketEarn 私密房门票收益(总/私密房/门票,内部加锁)
func (r *LiveRecord) AddPrivateRoomTicketEarn(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalIncome, &r.TotalIncome, v)
		if v > 0 {
			r.addFloatLocked(LiveRecordTotalPrivateRoomIncome, &r.TotalPrivateRoomIncome, v)
			r.addFloatLocked(LiveRecordTotalPrivateRoomTicketIncome, &r.TotalPrivateRoomTicketIncome, v)
		}
	})
}

// AddPrivateRoomWatchEarn 私密房观看收益(总/私密房/观看,内部加锁)
func (r *LiveRecord) AddPrivateRoomWatchEarn(v float64) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalIncome, &r.TotalIncome, v)
		if v > 0 {
			r.addFloatLocked(LiveRecordTotalPrivateRoomIncome, &r.TotalPrivateRoomIncome, v)
			r.addFloatLocked(LiveRecordTotalPrivateRoomWatchIncome, &r.TotalPrivateRoomWatchIncome, v)
		}
	})
}

// ApplyVideoCallIncomeDelta 通话收益增减(支持负数退款),内部加锁
func (r *LiveRecord) ApplyVideoCallIncomeDelta(amount float64, ticket, billing bool) {
	r.withLock(func() {
		r.addFloatLocked(LiveRecordTotalIncome, &r.TotalIncome, amount)
		r.addFloatLocked(LiveRecordTotalVideoCallIncome, &r.TotalVideoCallIncome, amount)
		if ticket {
			r.addFloatLocked(LiveRecordTotalVideoCallTicketIncome, &r.TotalVideoCallTicketIncome, amount)
		}
		if billing {
			r.addFloatLocked(LiveRecordTotalVideoCallBillingIncome, &r.TotalVideoCallBillingIncome, amount)
		}
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

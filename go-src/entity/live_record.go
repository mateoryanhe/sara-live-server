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
	LiveRecordAnchorId               db.TbCol = "anchor_id"
	LiveRecordStartTime              db.TbCol = "start_time"
	LiveRecordEndTime                db.TbCol = "end_time"
	LiveRecordTotalAudience          db.TbCol = "total_audience"
	LiveRecordTotalLiveDuration      db.TbCol = "total_live_duration"
	LiveRecordTotalIncome            db.TbCol = "total_income"
	LiveRecordTotalGiftIncome        db.TbCol = "total_gift_income"
	LiveRecordTotalPrivateRoomIncome db.TbCol = "total_private_room_income"
	LiveRecordTotalGameBet           db.TbCol = "total_game_bet"
	LiveRecordTotalGiftSender        db.TbCol = "total_gift_sender"
)

// LiveRecord 单场直播数据记录
type LiveRecord struct {
	migrate.OneModel
	AnchorId               uint64     `gorm:"index;default:0;comment:主播ID" json:"anchorId"`
	StartTime              time.Time  `gorm:"comment:直播开始时间" json:"startTime"`
	EndTime                *time.Time `gorm:"comment:直播结束时间" json:"endTime"`
	TotalAudience          uint64     `gorm:"default:0;comment:累计观众人数" json:"totalAudience"`
	TotalLiveDuration      float64    `gorm:"default:0;comment:累计直播时长(秒)" json:"totalLiveDuration"`
	TotalIncome            float64    `gorm:"default:0;comment:总收益" json:"totalIncome"`
	TotalGiftIncome        float64    `gorm:"default:0;comment:礼物收入" json:"totalGiftIncome"`
	TotalPrivateRoomIncome float64    `gorm:"default:0;comment:私密直播间收入" json:"totalPrivateRoomIncome"`
	TotalGameBet           float64    `gorm:"default:0;comment:游戏下注总金额" json:"totalGameBet"`
	TotalGiftSender        uint64     `gorm:"default:0;comment:送礼人数(去重)" json:"totalGiftSender"`
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
	syndb.AddDataToLazyChan(TbLiveRecord, LiveRecordAnchorId, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: v,
	})
}

func (r *LiveRecord) SetStartTime(v time.Time) {
	r.StartTime = v
	r.touchUpdatedAt()
	syndb.AddDataToLazyChan(TbLiveRecord, LiveRecordStartTime, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: v,
	})
}

func (r *LiveRecord) SetEndTime(v *time.Time) {
	r.EndTime = v
	syndb.AddDataToLazyChan(TbLiveRecord, LiveRecordEndTime, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.EndTime,
	})
}

func (r *LiveRecord) AddTotalAudience(val uint64) {
	r.TotalAudience = math.Add(r.TotalAudience, val)

	syndb.AddDataToLazyChan(TbLiveRecord, LiveRecordTotalAudience, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalAudience,
	})
}

func (r *LiveRecord) AddTotalLiveDuration(v float64) {
	r.TotalLiveDuration = math.AddFloat64(r.TotalLiveDuration, v)
	syndb.AddDataToLazyChan(TbLiveRecord, LiveRecordTotalLiveDuration, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalLiveDuration,
	})
}

func (r *LiveRecord) AddTotalIncome(v float64) {
	r.TotalIncome = math.AddFloat64(r.TotalIncome, v)

	syndb.AddDataToLazyChan(TbLiveRecord, LiveRecordTotalIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalIncome,
	})
}

func (r *LiveRecord) AddTotalGiftIncome(v float64) {
	r.TotalGiftIncome = math.AddFloat64(r.TotalGiftIncome, v)

	syndb.AddDataToLazyChan(TbLiveRecord, LiveRecordTotalGiftIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalGiftIncome,
	})
}

func (r *LiveRecord) AddTotalPrivateRoomIncome(v float64) {
	r.TotalPrivateRoomIncome = math.AddFloat64(r.TotalPrivateRoomIncome, v)

	syndb.AddDataToLazyChan(TbLiveRecord, LiveRecordTotalPrivateRoomIncome, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalPrivateRoomIncome,
	})
}

func (r *LiveRecord) AddTotalGameBet(v float64) {
	r.TotalGameBet = math.AddFloat64(r.TotalGameBet, v)

	syndb.AddDataToLazyChan(TbLiveRecord, LiveRecordTotalGameBet, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalGameBet,
	})
}

func (r *LiveRecord) AddTotalGiftSender(val uint64) {
	r.TotalGiftSender = math.Add(r.TotalGiftSender, val)

	syndb.AddDataToLazyChan(TbLiveRecord, LiveRecordTotalGiftSender, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.TotalGiftSender,
	})
}

func (r *LiveRecord) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddDataToLazyChan(TbLiveRecord, db.CreatedAtName, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: v,
	})
}

func (r *LiveRecord) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddDataToLazyChan(TbLiveRecord, db.UpdatedAtName, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: v,
	})
}

func (r *LiveRecord) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	syndb.AddDataToLazyChan(TbLiveRecord, db.UpdatedAtName, &syndb.ColData{
		IdVal:  r.ID,
		ColVal: r.UpdatedAt,
	})
}

func initLiveRecord() {
	syndb.RegLazyWithMiddle(TbLiveRecord, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbLiveRecord, db.UpdatedAtName)
	syndb.RegLazyWithMiddle(TbLiveRecord, LiveRecordAnchorId)
	syndb.RegLazyWithMiddle(TbLiveRecord, LiveRecordStartTime)
	syndb.RegLazyWithMiddle(TbLiveRecord, LiveRecordEndTime)
	syndb.RegLazyWithMiddle(TbLiveRecord, LiveRecordTotalAudience)
	syndb.RegLazyWithMiddle(TbLiveRecord, LiveRecordTotalLiveDuration)
	syndb.RegLazyWithMiddle(TbLiveRecord, LiveRecordTotalIncome)
	syndb.RegLazyWithMiddle(TbLiveRecord, LiveRecordTotalGiftIncome)
	syndb.RegLazyWithMiddle(TbLiveRecord, LiveRecordTotalPrivateRoomIncome)
	syndb.RegLazyWithMiddle(TbLiveRecord, LiveRecordTotalGameBet)
	syndb.RegLazyWithMiddle(TbLiveRecord, LiveRecordTotalGiftSender)
	migrate.AutoMigrate(&LiveRecord{})
}

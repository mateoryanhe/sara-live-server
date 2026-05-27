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
	LiveRecordAnchorId          db.TbCol = "anchor_id"
	LiveRecordStartTime         db.TbCol = "start_time"
	LiveRecordEndTime           db.TbCol = "end_time"
	LiveRecordTotalAudience     db.TbCol = "total_audience"
	LiveRecordTotalLiveDuration db.TbCol = "total_live_duration"
)

// LiveRecord 单场直播数据记录
type LiveRecord struct {
	migrate.OneModel
	AnchorId          uint64     `gorm:"index;default:0;comment:主播ID" json:"anchorId"`
	StartTime         time.Time  `gorm:"comment:直播开始时间" json:"startTime"`
	EndTime           *time.Time `gorm:"comment:直播结束时间" json:"endTime"`
	TotalAudience     uint64     `gorm:"default:0;comment:累计观众人数" json:"totalAudience"`
	TotalLiveDuration float64    `gorm:"default:0;comment:累计直播时长(秒)" json:"totalLiveDuration"`
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
	r.touchUpdatedAt()
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
	migrate.AutoMigrate(&LiveRecord{})
}

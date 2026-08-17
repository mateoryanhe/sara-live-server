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
	TbDailyAnchorEffectiveLive db.TbName = "daily_anchor_effective_lives"
)

// MinAccumulateLiveSessionSec 单场直播时长大于此值(秒)才计入日表累计时长
const MinAccumulateLiveSessionSec = 30 * 60

const (
	DailyAnchorEffectiveLiveRoomId       db.TbCol = "room_id"
	DailyAnchorEffectiveLiveLiveDate     db.TbCol = "live_date"
	DailyAnchorEffectiveLiveLiveDuration db.TbCol = "live_duration"
	DailyAnchorEffectiveLiveSettled      db.TbCol = "settled"
)

// DailyAnchorEffectiveLive 主播每日直播时长统计
// 主键 ID = "{date}_{roomId}"(roomId==主播用户ID),字段经 syndb 缓冲落库
type DailyAnchorEffectiveLive struct {
	ID           string    `gorm:"primaryKey;size:64;comment:复合ID(date_roomId)" json:"id"`
	RoomId       uint64    `gorm:"index;default:0;comment:直播间ID(==主播用户ID)" json:"roomId"`
	LiveDate     string    `gorm:"size:10;index;default:'';comment:日期(YYYY-MM-DD)" json:"liveDate"`
	LiveDuration float64   `gorm:"default:0;comment:当日累计直播时长(秒,仅统计单场>30分钟)" json:"liveDuration"`
	Settled      bool      `gorm:"default:0;comment:结算标记(0未结算,1已结算)" json:"settled"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// FormatDailyAnchorEffectiveLiveDate 格式化统计日期
func FormatDailyAnchorEffectiveLiveDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// BuildDailyAnchorEffectiveLiveId 拼接复合主键
func BuildDailyAnchorEffectiveLiveId(date string, roomId uint64) string {
	return fmt.Sprintf("%s_%d", date, roomId)
}

// NewDailyAnchorEffectiveLive 构造内存对象,字段写入通过 syndb 异步入库
func NewDailyAnchorEffectiveLive(date string, roomId uint64) *DailyAnchorEffectiveLive {
	r := &DailyAnchorEffectiveLive{}
	r.ID = BuildDailyAnchorEffectiveLiveId(date, roomId)
	now := time.Now()
	r.SetCreatedAt(now)
	r.SetUpdatedAt(now)
	r.SetRoomId(roomId)
	r.SetLiveDate(date)
	return r
}

func (r *DailyAnchorEffectiveLive) SetRoomId(v uint64) {
	r.RoomId = v
	syndb.AddData(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveRoomId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyAnchorEffectiveLive) SetLiveDate(v string) {
	r.LiveDate = v
	syndb.AddData(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveLiveDate, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyAnchorEffectiveLive) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbDailyAnchorEffectiveLive, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyAnchorEffectiveLive) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbDailyAnchorEffectiveLive, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyAnchorEffectiveLive) SetSettled(v bool) {
	r.Settled = v
	syndb.AddData(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveSettled, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

// SetLiveDuration 设置当日累计直播时长(可置0,syndb 缓冲)
func (r *DailyAnchorEffectiveLive) SetLiveDuration(v float64) {
	if r == nil {
		return
	}
	key := fmt.Sprintf("daily_anchor_effective_live:%s", r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	r.LiveDuration = v
	r.UpdatedAt = time.Now()
	syndb.AddData(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveLiveDuration, &syndb.ColData{
		IdVal: r.ID, ColVal: r.LiveDuration,
	})
	syndb.AddData(TbDailyAnchorEffectiveLive, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.UpdatedAt,
	})
}

// AddLiveDuration 累加当日直播时长(syndb 缓冲)
func (r *DailyAnchorEffectiveLive) AddLiveDuration(v float64) {
	if r == nil || v <= 0 {
		return
	}
	key := fmt.Sprintf("daily_anchor_effective_live:%s", r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	r.LiveDuration = math.AddFloat64(r.LiveDuration, v)
	r.UpdatedAt = time.Now()
	syndb.AddData(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveLiveDuration, &syndb.ColData{
		IdVal: r.ID, ColVal: r.LiveDuration,
	})
	syndb.AddData(TbDailyAnchorEffectiveLive, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.UpdatedAt,
	})
}

func initDailyAnchorEffectiveLive() {
	syndb.RegQuick(TbDailyAnchorEffectiveLive, db.CreatedAtName)
	syndb.RegLazy(TbDailyAnchorEffectiveLive, db.UpdatedAtName)
	syndb.RegQuick(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveRoomId)
	syndb.RegQuick(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveLiveDate)
	syndb.RegLazy(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveLiveDuration)
	syndb.RegQuick(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveSettled)
	migrate.AutoMigrate(&DailyAnchorEffectiveLive{})
}

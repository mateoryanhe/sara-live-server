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

const (
	DailyAnchorEffectiveLiveRoomId             db.TbCol = "room_id"
	DailyAnchorEffectiveLiveLiveDate           db.TbCol = "live_date"
	DailyAnchorEffectiveLiveEffectiveLiveCount db.TbCol = "effective_live_count"
)

// DailyAnchorEffectiveLive 主播每日有效直播次数
// 主键 ID = "{date}_{roomId}"(roomId==主播用户ID),字段经 syndb 缓冲落库
type DailyAnchorEffectiveLive struct {
	ID                 string    `gorm:"primaryKey;size:64;comment:复合ID(date_roomId)" json:"id"`
	RoomId             uint64    `gorm:"index;default:0;comment:直播间ID(==主播用户ID)" json:"roomId"`
	LiveDate           string    `gorm:"size:10;index;default:'';comment:日期(YYYY-MM-DD)" json:"liveDate"`
	EffectiveLiveCount uint64    `gorm:"default:0;comment:当日有效直播次数" json:"effectiveLiveCount"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
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

// AddEffectiveLiveCount 累加当日有效直播次数(syndb 缓冲)
func (r *DailyAnchorEffectiveLive) AddEffectiveLiveCount(v uint64) {
	if r == nil || v == 0 {
		return
	}
	key := fmt.Sprintf("daily_anchor_effective_live:%s", r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	r.EffectiveLiveCount = math.Add(r.EffectiveLiveCount, v)
	r.UpdatedAt = time.Now()
	syndb.AddData(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveEffectiveLiveCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.EffectiveLiveCount,
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
	syndb.RegLazy(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveEffectiveLiveCount)
	migrate.AutoMigrate(&DailyAnchorEffectiveLive{})
}

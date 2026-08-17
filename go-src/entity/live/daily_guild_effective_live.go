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
	TbDailyGuildEffectiveLive db.TbName = "daily_guild_effective_lives"
)

const (
	DailyGuildEffectiveLiveGuildId      db.TbCol = "guild_id"
	DailyGuildEffectiveLiveLiveDate     db.TbCol = "live_date"
	DailyGuildEffectiveLiveLiveDuration db.TbCol = "live_duration"
	DailyGuildEffectiveLiveSettled      db.TbCol = "settled"
)

// DailyGuildEffectiveLive 工会每日直播时长统计
// 主键 ID = "{date}_{guildId}",字段经 syndb 缓冲落库
type DailyGuildEffectiveLive struct {
	ID           string    `gorm:"primaryKey;size:64;comment:复合ID(date_guildId)" json:"id"`
	GuildId      uint64    `gorm:"index;default:0;comment:工会ID" json:"guildId"`
	LiveDate     string    `gorm:"size:10;index;default:'';comment:日期(YYYY-MM-DD)" json:"liveDate"`
	LiveDuration float64   `gorm:"default:0;comment:当日累计直播时长(秒,仅统计单场>30分钟)" json:"liveDuration"`
	Settled      bool      `gorm:"default:0;comment:结算标记(0未结算,1已结算)" json:"settled"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func BuildDailyGuildEffectiveLiveId(date string, guildId uint64) string {
	return fmt.Sprintf("%s_%d", date, guildId)
}

func NewDailyGuildEffectiveLive(date string, guildId uint64) *DailyGuildEffectiveLive {
	r := &DailyGuildEffectiveLive{}
	r.ID = BuildDailyGuildEffectiveLiveId(date, guildId)
	now := time.Now()
	r.SetCreatedAt(now)
	r.SetUpdatedAt(now)
	r.SetGuildId(guildId)
	r.SetLiveDate(date)
	return r
}

func (r *DailyGuildEffectiveLive) SetGuildId(v uint64) {
	r.GuildId = v
	syndb.AddData(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveGuildId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyGuildEffectiveLive) SetLiveDate(v string) {
	r.LiveDate = v
	syndb.AddData(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveLiveDate, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyGuildEffectiveLive) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbDailyGuildEffectiveLive, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyGuildEffectiveLive) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbDailyGuildEffectiveLive, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *DailyGuildEffectiveLive) SetSettled(v bool) {
	r.Settled = v
	syndb.AddData(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveSettled, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

// SetLiveDuration 设置当日累计直播时长(可置0,syndb 缓冲)
func (r *DailyGuildEffectiveLive) SetLiveDuration(v float64) {
	if r == nil {
		return
	}
	key := fmt.Sprintf("daily_guild_effective_live:%s", r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	r.LiveDuration = v
	r.UpdatedAt = time.Now()
	syndb.AddData(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveLiveDuration, &syndb.ColData{
		IdVal: r.ID, ColVal: r.LiveDuration,
	})
	syndb.AddData(TbDailyGuildEffectiveLive, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.UpdatedAt,
	})
}

// AddLiveDuration 累加当日直播时长(syndb 缓冲)
func (r *DailyGuildEffectiveLive) AddLiveDuration(v float64) {
	if r == nil || v <= 0 {
		return
	}
	key := fmt.Sprintf("daily_guild_effective_live:%s", r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	r.LiveDuration = math.AddFloat64(r.LiveDuration, v)
	r.UpdatedAt = time.Now()
	syndb.AddData(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveLiveDuration, &syndb.ColData{
		IdVal: r.ID, ColVal: r.LiveDuration,
	})
	syndb.AddData(TbDailyGuildEffectiveLive, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.UpdatedAt,
	})
}

func initDailyGuildEffectiveLive() {
	syndb.RegQuick(TbDailyGuildEffectiveLive, db.CreatedAtName)
	syndb.RegLazy(TbDailyGuildEffectiveLive, db.UpdatedAtName)
	syndb.RegQuick(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveGuildId)
	syndb.RegQuick(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveLiveDate)
	syndb.RegLazy(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveLiveDuration)
	syndb.RegQuick(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveSettled)
	migrate.AutoMigrate(&DailyGuildEffectiveLive{})
}

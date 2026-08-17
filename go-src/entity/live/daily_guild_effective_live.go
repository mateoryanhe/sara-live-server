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
	DailyGuildEffectiveLiveGuildId            db.TbCol = "guild_id"
	DailyGuildEffectiveLiveLiveDate           db.TbCol = "live_date"
	DailyGuildEffectiveLiveEffectiveLiveCount db.TbCol = "effective_live_count"
	DailyGuildEffectiveLiveSettled            db.TbCol = "settled"
)

// DailyGuildEffectiveLive 工会每日有效直播次数
// 主键 ID = "{date}_{guildId}",字段经 syndb 缓冲落库
type DailyGuildEffectiveLive struct {
	ID                 string    `gorm:"primaryKey;size:64;comment:复合ID(date_guildId)" json:"id"`
	GuildId            uint64    `gorm:"index;default:0;comment:工会ID" json:"guildId"`
	LiveDate           string    `gorm:"size:10;index;default:'';comment:日期(YYYY-MM-DD)" json:"liveDate"`
	EffectiveLiveCount uint64    `gorm:"default:0;comment:当日有效直播次数" json:"effectiveLiveCount"`
	Settled            bool      `gorm:"default:0;comment:结算标记(0未结算,1已结算)" json:"settled"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
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

// SetEffectiveLiveCount 设置当日有效直播次数(可置0,syndb 缓冲)
func (r *DailyGuildEffectiveLive) SetEffectiveLiveCount(v uint64) {
	if r == nil {
		return
	}
	key := fmt.Sprintf("daily_guild_effective_live:%s", r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	r.EffectiveLiveCount = v
	r.UpdatedAt = time.Now()
	syndb.AddData(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveEffectiveLiveCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.EffectiveLiveCount,
	})
	syndb.AddData(TbDailyGuildEffectiveLive, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.UpdatedAt,
	})
}

func (r *DailyGuildEffectiveLive) AddEffectiveLiveCount(v uint64) {
	if r == nil || v == 0 {
		return
	}
	key := fmt.Sprintf("daily_guild_effective_live:%s", r.ID)
	gmlock.Lock(key)
	defer gmlock.Unlock(key)
	r.EffectiveLiveCount = math.Add(r.EffectiveLiveCount, v)
	r.UpdatedAt = time.Now()
	syndb.AddData(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveEffectiveLiveCount, &syndb.ColData{
		IdVal: r.ID, ColVal: r.EffectiveLiveCount,
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
	syndb.RegLazy(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveEffectiveLiveCount)
	syndb.RegQuick(TbDailyGuildEffectiveLive, DailyGuildEffectiveLiveSettled)
	migrate.AutoMigrate(&DailyGuildEffectiveLive{})
}

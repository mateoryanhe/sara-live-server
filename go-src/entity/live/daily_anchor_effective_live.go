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

// DailyAnchorEffectiveLive 主播每日直播时长与收益流水
// 主键 ID = "{date}_{roomId}"(roomId==主播用户ID),字段经 syndb 缓冲落库
type DailyAnchorEffectiveLive struct {
	ID           string    `gorm:"primaryKey;size:64;comment:复合ID(date_roomId)" json:"id"`
	RoomId       uint64    `gorm:"index;default:0;comment:直播间ID(==主播用户ID)" json:"roomId"`
	LiveDate     string    `gorm:"size:10;index;default:'';comment:日期(YYYY-MM-DD)" json:"liveDate"`
	LiveDuration float64   `gorm:"default:0;comment:当日累计直播时长(秒,仅统计单场>30分钟)" json:"liveDuration"`
	Settled      bool      `gorm:"default:0;comment:结算标记(0未结算,1已结算)" json:"settled"`
	LiveRoomIncomeAmounts
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

func (r *DailyAnchorEffectiveLive) lockKey() string {
	return fmt.Sprintf("daily_anchor_effective_live:%s", r.ID)
}

// AddGiftEarn 礼物收益(总收益+礼物细分,内部加锁)
func (r *DailyAnchorEffectiveLive) AddGiftEarn(v float64) {
	if r == nil || r.ID == "" {
		return
	}
	addIncomeEarnWithLockKey(TbDailyAnchorEffectiveLive, r.lockKey(), r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalGiftIncome, &r.TotalGiftIncome)
}

// AddPaidDanmakuEarn 付费弹幕收益(总收益+弹幕细分,内部加锁)
func (r *DailyAnchorEffectiveLive) AddPaidDanmakuEarn(v float64) {
	if r == nil || r.ID == "" {
		return
	}
	addIncomeEarnWithLockKey(TbDailyAnchorEffectiveLive, r.lockKey(), r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPaidDanmakuIncome, &r.TotalPaidDanmakuIncome)
}

// AddPrivateRoomTicketEarn 私密房门票收益(总收益+门票细分,内部加锁)
func (r *DailyAnchorEffectiveLive) AddPrivateRoomTicketEarn(v float64) {
	if r == nil || r.ID == "" {
		return
	}
	addIncomeEarnWithLockKey(TbDailyAnchorEffectiveLive, r.lockKey(), r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPrivateRoomTicketIncome, &r.TotalPrivateRoomTicketIncome)
}

// AddPrivateRoomWatchEarn 私密房观看收益(总收益+观看细分,内部加锁)
func (r *DailyAnchorEffectiveLive) AddPrivateRoomWatchEarn(v float64) {
	if r == nil || r.ID == "" {
		return
	}
	addIncomeEarnWithLockKey(TbDailyAnchorEffectiveLive, r.lockKey(), r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalPrivateRoomWatchIncome, &r.TotalPrivateRoomWatchIncome)
}

// AddShortVideoEarn 短视频付费观看收益(总收益+短视频细分,内部加锁)
func (r *DailyAnchorEffectiveLive) AddShortVideoEarn(v float64) {
	if r == nil || r.ID == "" {
		return
	}
	addIncomeEarnWithLockKey(TbDailyAnchorEffectiveLive, r.lockKey(), r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, v, LiveRoomIncomeTotalShortVideoIncome, &r.TotalShortVideoIncome)
}

// AddGameEarn 游戏收益(总收益按钻石折算 + 游戏细分金币下注额,内部加锁)
func (r *DailyAnchorEffectiveLive) AddGameEarn(goldAmount, incomeDelta float64) {
	if r == nil || r.ID == "" {
		return
	}
	addGameEarnWithLockKey(TbDailyAnchorEffectiveLive, r.lockKey(), r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, goldAmount, incomeDelta)
}

// ApplyVideoCallIncomeDelta 通话收益增减(支持负数退款,内部加锁)
func (r *DailyAnchorEffectiveLive) ApplyVideoCallIncomeDelta(amount float64, ticket, billing bool) {
	if r == nil || r.ID == "" {
		return
	}
	ApplyVideoCallIncomeDeltaWithLockKey(TbDailyAnchorEffectiveLive, r.lockKey(), r.ID, &r.LiveRoomIncomeAmounts, &r.UpdatedAt, amount, ticket, billing)
}

// AddTotalLiveDuration 累加心跳上报直播时长(秒,syndb 缓冲)
func (r *DailyAnchorEffectiveLive) AddTotalLiveDuration(v float64) {
	if r == nil || r.ID == "" || v <= 0 {
		return
	}
	gmlock.Lock(r.lockKey())
	defer gmlock.Unlock(r.lockKey())
	addIncomeAmountLocked(TbDailyAnchorEffectiveLive, LiveRoomIncomeTotalLiveDuration, r.ID, &r.TotalLiveDuration, v)
	touchIncomeUpdatedAt(TbDailyAnchorEffectiveLive, r.ID, &r.UpdatedAt)
}

func initDailyAnchorEffectiveLive() {
	syndb.RegQuick(TbDailyAnchorEffectiveLive, db.CreatedAtName)
	syndb.RegLazy(TbDailyAnchorEffectiveLive, db.UpdatedAtName)
	syndb.RegQuick(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveRoomId)
	syndb.RegQuick(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveLiveDate)
	syndb.RegLazy(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveLiveDuration)
	syndb.RegQuick(TbDailyAnchorEffectiveLive, DailyAnchorEffectiveLiveSettled)
	regLiveRoomIncomeCols(TbDailyAnchorEffectiveLive)
	migrate.AutoMigrate(&DailyAnchorEffectiveLive{})
}

package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbMinuteLiveActivityStat db.TbName = "minute_live_activity_stats"

const (
	MinuteLiveActivityStatBucketAt     db.TbCol = "bucket_at"
	MinuteLiveActivityStatJoinCount    db.TbCol = "join_count"
	MinuteLiveActivityStatLeaveCount   db.TbCol = "leave_count"
	MinuteLiveActivityStatGiftCount    db.TbCol = "gift_count"
	MinuteLiveActivityStatChatCount    db.TbCol = "chat_count"
	MinuteLiveActivityStatBetCount           db.TbCol = "bet_count"
	MinuteLiveActivityStatCall1v1Count       db.TbCol = "call1v1_count"
	MinuteLiveActivityStatRechargeOrderCount db.TbCol = "recharge_order_count"
)

// MinuteLiveActivityStat 直播间分钟级活跃计数(UTC 分钟桶,保留约 3 天)
type MinuteLiveActivityStat struct {
	migrate.OneModel
	BucketAt     time.Time `gorm:"uniqueIndex;comment:分钟桶起点(UTC)" json:"bucketAt"`
	JoinCount    uint64    `gorm:"default:0;comment:进房次数" json:"joinCount"`
	LeaveCount   uint64    `gorm:"default:0;comment:离房次数" json:"leaveCount"`
	GiftCount    uint64    `gorm:"default:0;comment:送礼次数" json:"giftCount"`
	ChatCount    uint64    `gorm:"default:0;comment:公屏聊天次数(含付费弹幕)" json:"chatCount"`
	BetCount           uint64    `gorm:"default:0;comment:游戏下注次数" json:"betCount"`
	Call1v1Count       uint64    `gorm:"default:0;comment:1v1发起次数" json:"call1v1Count"`
	RechargeOrderCount uint64    `gorm:"default:0;comment:充值下单次数" json:"rechargeOrderCount"`
}

func TruncateUTCMinute(at time.Time) time.Time {
	if at.IsZero() {
		at = time.Now()
	}
	return at.UTC().Truncate(time.Minute)
}

func MinuteLiveActivityStatID(at time.Time) uint64 {
	return uint64(TruncateUTCMinute(at).Unix())
}

func NewMinuteLiveActivityStat(at time.Time) *MinuteLiveActivityStat {
	bucket := TruncateUTCMinute(at)
	row := &MinuteLiveActivityStat{BucketAt: bucket}
	row.ID = MinuteLiveActivityStatID(at)
	row.touchBucketAt()
	return row
}

func (r *MinuteLiveActivityStat) touchBucketAt() {
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatBucketAt, &syndb.ColData{IdVal: r.ID, ColVal: r.BucketAt})
}

func (r *MinuteLiveActivityStat) AddJoinCount(n uint64) {
	r.JoinCount = math.Add(r.JoinCount, n)
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatJoinCount, &syndb.ColData{IdVal: r.ID, ColVal: r.JoinCount})
}

func (r *MinuteLiveActivityStat) AddLeaveCount(n uint64) {
	r.LeaveCount = math.Add(r.LeaveCount, n)
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatLeaveCount, &syndb.ColData{IdVal: r.ID, ColVal: r.LeaveCount})
}

func (r *MinuteLiveActivityStat) AddGiftCount(n uint64) {
	r.GiftCount = math.Add(r.GiftCount, n)
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatGiftCount, &syndb.ColData{IdVal: r.ID, ColVal: r.GiftCount})
}

func (r *MinuteLiveActivityStat) AddChatCount(n uint64) {
	r.ChatCount = math.Add(r.ChatCount, n)
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatChatCount, &syndb.ColData{IdVal: r.ID, ColVal: r.ChatCount})
}

func (r *MinuteLiveActivityStat) AddBetCount(n uint64) {
	r.BetCount = math.Add(r.BetCount, n)
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatBetCount, &syndb.ColData{IdVal: r.ID, ColVal: r.BetCount})
}

func (r *MinuteLiveActivityStat) AddCall1v1Count(n uint64) {
	r.Call1v1Count = math.Add(r.Call1v1Count, n)
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatCall1v1Count, &syndb.ColData{IdVal: r.ID, ColVal: r.Call1v1Count})
}

func (r *MinuteLiveActivityStat) AddRechargeOrderCount(n uint64) {
	r.RechargeOrderCount = math.Add(r.RechargeOrderCount, n)
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatRechargeOrderCount, &syndb.ColData{IdVal: r.ID, ColVal: r.RechargeOrderCount})
}

func initMinuteLiveActivityStat() {
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatBucketAt)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatJoinCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatLeaveCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatGiftCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatChatCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatBetCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatCall1v1Count)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatRechargeOrderCount)
	migrate.AutoMigrate(&MinuteLiveActivityStat{})
}

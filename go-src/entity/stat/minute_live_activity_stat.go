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
	MinuteLiveActivityStatPlatform     db.TbCol = "platform"
	MinuteLiveActivityStatBucketAt     db.TbCol = "bucket_at"
	MinuteLiveActivityStatJoinCount    db.TbCol = "join_count"
	MinuteLiveActivityStatLeaveCount   db.TbCol = "leave_count"
	MinuteLiveActivityStatGiftCount    db.TbCol = "gift_count"
	MinuteLiveActivityStatChatCount    db.TbCol = "chat_count"
	MinuteLiveActivityStatBetCount           db.TbCol = "bet_count"
	MinuteLiveActivityStatCall1v1Count              db.TbCol = "call1v1_count"
	MinuteLiveActivityStatCall1v1ConnectSuccessCount db.TbCol = "call1v1_connect_success_count"
	MinuteLiveActivityStatCall1v1ConsumeCount        db.TbCol = "call1v1_consume_count"
	MinuteLiveActivityStatRechargeOrderCount         db.TbCol = "recharge_order_count"
)

// MinuteLiveActivityStat 直播间分钟级活跃计数(UTC 分钟桶,保留约 3 天)
type MinuteLiveActivityStat struct {
	migrate.OneModel
	BucketAt time.Time `gorm:"uniqueIndex:idx_minute_live_activity_bucket_platform,priority:1;comment:分钟桶起点(UTC)" json:"bucketAt"`
	Platform string    `gorm:"size:16;default:'';uniqueIndex:idx_minute_live_activity_bucket_platform,priority:2;comment:客户端平台 android/ios/web 空=未知" json:"platform"`
	JoinCount    uint64    `gorm:"default:0;comment:进房次数" json:"joinCount"`
	LeaveCount   uint64    `gorm:"default:0;comment:离房次数" json:"leaveCount"`
	GiftCount    uint64    `gorm:"default:0;comment:送礼次数" json:"giftCount"`
	ChatCount    uint64    `gorm:"default:0;comment:公屏聊天次数(含付费弹幕)" json:"chatCount"`
	BetCount                   uint64    `gorm:"default:0;comment:游戏下注次数" json:"betCount"`
	Call1v1Count               uint64    `gorm:"default:0;comment:1v1发起次数" json:"call1v1Count"`
	Call1v1ConnectSuccessCount uint64    `gorm:"default:0;comment:1v1成功连线次数" json:"call1v1ConnectSuccessCount"`
	Call1v1ConsumeCount        uint64    `gorm:"default:0;comment:1v1成功消费次数(分钟计费扣款)" json:"call1v1ConsumeCount"`
	RechargeOrderCount         uint64    `gorm:"default:0;comment:充值下单次数" json:"rechargeOrderCount"`
}

func TruncateUTCMinute(at time.Time) time.Time {
	if at.IsZero() {
		at = time.Now()
	}
	return at.UTC().Truncate(time.Minute)
}

func MinuteLiveActivityStatID(at time.Time, platform string) uint64 {
	base := uint64(TruncateUTCMinute(at).Unix())
	switch platform {
	case "android":
		return base*4 + 1
	case "ios":
		return base*4 + 2
	case "web":
		return base*4 + 3
	default:
		return base
	}
}

func NewMinuteLiveActivityStat(at time.Time, platform string) *MinuteLiveActivityStat {
	bucket := TruncateUTCMinute(at)
	row := &MinuteLiveActivityStat{BucketAt: bucket, Platform: platform}
	row.ID = MinuteLiveActivityStatID(at, platform)
	row.touchBucketAt()
	row.touchPlatform()
	return row
}

func (r *MinuteLiveActivityStat) touchPlatform() {
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatPlatform, &syndb.ColData{IdVal: r.ID, ColVal: r.Platform})
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

func (r *MinuteLiveActivityStat) AddCall1v1ConnectSuccessCount(n uint64) {
	r.Call1v1ConnectSuccessCount = math.Add(r.Call1v1ConnectSuccessCount, n)
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatCall1v1ConnectSuccessCount, &syndb.ColData{IdVal: r.ID, ColVal: r.Call1v1ConnectSuccessCount})
}

func (r *MinuteLiveActivityStat) AddCall1v1ConsumeCount(n uint64) {
	r.Call1v1ConsumeCount = math.Add(r.Call1v1ConsumeCount, n)
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatCall1v1ConsumeCount, &syndb.ColData{IdVal: r.ID, ColVal: r.Call1v1ConsumeCount})
}

func (r *MinuteLiveActivityStat) AddRechargeOrderCount(n uint64) {
	r.RechargeOrderCount = math.Add(r.RechargeOrderCount, n)
	syndb.AddData(TbMinuteLiveActivityStat, MinuteLiveActivityStatRechargeOrderCount, &syndb.ColData{IdVal: r.ID, ColVal: r.RechargeOrderCount})
}

func initMinuteLiveActivityStat() {
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatPlatform)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatBucketAt)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatJoinCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatLeaveCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatGiftCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatChatCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatBetCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatCall1v1Count)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatCall1v1ConnectSuccessCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatCall1v1ConsumeCount)
	syndb.RegQuick(TbMinuteLiveActivityStat, MinuteLiveActivityStatRechargeOrderCount)
	migrate.AutoMigrate(&MinuteLiveActivityStat{})
}

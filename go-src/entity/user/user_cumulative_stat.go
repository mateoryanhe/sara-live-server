package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbUserCumulativeStat db.TbName = "user_cumulative_stats"
)

const (
	UserCumulativeStatTotalRecharge       db.TbCol = "total_recharge"
	UserCumulativeStatTotalWithdraw       db.TbCol = "total_withdraw"
	UserCumulativeStatTotalFans           db.TbCol = "total_fans"
	UserCumulativeStatTotalFollow         db.TbCol = "total_follow"
	UserCumulativeStatTotalPayCount       db.TbCol = "total_pay_count"
	UserCumulativeStatTotalDiamondConsume db.TbCol = "total_diamond_consume"
	UserCumulativeStatTotalGoldConsume    db.TbCol = "total_gold_consume"
	UserCumulativeStatTotalLiveDuration   db.TbCol = "total_live_duration"
)

// UserCumulativeStat 玩家累计数值(与用户一一对应,主键ID即用户ID)
type UserCumulativeStat struct {
	migrate.OneModel
	TotalRecharge       float64 `gorm:"default:0;comment:累计充值" json:"totalRecharge"`
	TotalWithdraw       float64 `gorm:"default:0;comment:累计提现" json:"totalWithdraw"`
	TotalFans           uint64  `gorm:"default:0;comment:累计粉丝数量" json:"totalFans"`
	TotalFollow         uint64  `gorm:"default:0;comment:累计关注数量" json:"totalFollow"`
	TotalPayCount       uint64  `gorm:"default:0;comment:累计付费次数" json:"totalPayCount"`
	TotalDiamondConsume float64 `gorm:"default:0;comment:累计钻石消费" json:"totalDiamondConsume"`
	TotalGoldConsume    float64 `gorm:"default:0;comment:累计金币消费" json:"totalGoldConsume"`
	TotalLiveDuration   float64 `gorm:"default:0;comment:累计直播时长(秒)" json:"totalLiveDuration"`
}

func NewUserCumulativeStat(userId uint64) *UserCumulativeStat {
	ret := &UserCumulativeStat{}
	ret.ID = userId
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

func (receiver *UserCumulativeStat) AddTotalRecharge(val float64) bool {
	receiver.TotalRecharge = math.AddFloat64(receiver.TotalRecharge, val)

	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalRecharge, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalRecharge,
	})
	return true
}

func (receiver *UserCumulativeStat) AddTotalWithdraw(val float64) bool {
	receiver.TotalWithdraw = math.AddFloat64(receiver.TotalWithdraw, val)

	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalWithdraw, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalWithdraw,
	})
	return true
}

func (receiver *UserCumulativeStat) AddTotalFans(val uint64) bool {
	receiver.TotalFans = math.Add(receiver.TotalFans, val)

	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalFans, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalFans,
	})
	return true
}

func (receiver *UserCumulativeStat) AddTotalFollow(val uint64) bool {
	receiver.TotalFollow = math.Add(receiver.TotalFollow, val)

	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalFollow, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalFollow,
	})
	return true
}

func (receiver *UserCumulativeStat) AddTotalPayCount(val uint64) bool {
	receiver.TotalPayCount = math.Add(receiver.TotalPayCount, val)

	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalPayCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalPayCount,
	})
	return true
}

func (receiver *UserCumulativeStat) AddTotalDiamondConsume(val float64) bool {

	receiver.TotalDiamondConsume = math.AddFloat64(val, receiver.TotalDiamondConsume)

	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalDiamondConsume, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalDiamondConsume,
	})
	return true
}

func (receiver *UserCumulativeStat) AddTotalGoldConsume(val float64) bool {
	receiver.TotalGoldConsume = math.AddFloat64(val, receiver.TotalGoldConsume)

	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalGoldConsume, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalGoldConsume,
	})
	return true
}

func (receiver *UserCumulativeStat) AddTotalLiveDuration(val float64) bool {
	receiver.TotalLiveDuration = math.AddFloat64(receiver.TotalLiveDuration, val)
	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalLiveDuration, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalLiveDuration,
	})
	return true
}

func (receiver *UserCumulativeStat) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddData(TbUserCumulativeStat, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserCumulativeStat) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddData(TbUserCumulativeStat, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initUserCumulativeStat() {
	syndb.RegLazy(TbUserCumulativeStat, db.CreatedAtName)
	syndb.RegLazy(TbUserCumulativeStat, db.UpdatedAtName)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalRecharge)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalWithdraw)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalFans)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalFollow)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalPayCount)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalDiamondConsume)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalGoldConsume)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalLiveDuration)

	migrate.AutoMigrate(&UserCumulativeStat{})
}

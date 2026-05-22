package entity

import (
	"math"
	"time"
	"xr-game-server/constants/db"
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
)

const (
	cumulativeStatFloatScale        = 1000
	cumulativeStatFloatDecimalScale = 10000 // 保留4位小数
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
}

func NewUserCumulativeStat(userId uint64) *UserCumulativeStat {
	ret := &UserCumulativeStat{}
	ret.ID = userId
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

// normalizeCumulativeFloat 截断到4位小数后按1000倍缩放校验,防止浮点溢出
func normalizeCumulativeFloat(val float64) (float64, bool) {
	if math.IsNaN(val) || math.IsInf(val, 0) || val < 0 {
		return 0, false
	}
	truncated := math.Floor(val*cumulativeStatFloatDecimalScale) / cumulativeStatFloatDecimalScale
	scaled := truncated * cumulativeStatFloatScale
	if scaled > float64(math.MaxInt64) {
		return 0, false
	}
	return truncated, true
}

func (receiver *UserCumulativeStat) SetTotalRecharge(val float64) bool {
	normalized, ok := normalizeCumulativeFloat(val)
	if !ok {
		return false
	}
	receiver.TotalRecharge = normalized
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserCumulativeStat, UserCumulativeStatTotalRecharge, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: normalized,
	})
	return true
}

func (receiver *UserCumulativeStat) SetTotalWithdraw(val float64) bool {
	normalized, ok := normalizeCumulativeFloat(val)
	if !ok {
		return false
	}
	receiver.TotalWithdraw = normalized
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserCumulativeStat, UserCumulativeStatTotalWithdraw, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: normalized,
	})
	return true
}

func (receiver *UserCumulativeStat) SetTotalFans(val uint64) bool {
	receiver.TotalFans = val
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserCumulativeStat, UserCumulativeStatTotalFans, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
	return true
}

func (receiver *UserCumulativeStat) SetTotalFollow(val uint64) bool {
	receiver.TotalFollow = val
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserCumulativeStat, UserCumulativeStatTotalFollow, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
	return true
}

func (receiver *UserCumulativeStat) SetTotalPayCount(val uint64) bool {
	receiver.TotalPayCount = val
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserCumulativeStat, UserCumulativeStatTotalPayCount, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
	return true
}

func (receiver *UserCumulativeStat) SetTotalDiamondConsume(val float64) bool {
	normalized, ok := normalizeCumulativeFloat(val)
	if !ok {
		return false
	}
	receiver.TotalDiamondConsume = normalized
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserCumulativeStat, UserCumulativeStatTotalDiamondConsume, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: normalized,
	})
	return true
}

func (receiver *UserCumulativeStat) SetTotalGoldConsume(val float64) bool {
	normalized, ok := normalizeCumulativeFloat(val)
	if !ok {
		return false
	}
	receiver.TotalGoldConsume = normalized
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserCumulativeStat, UserCumulativeStatTotalGoldConsume, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: normalized,
	})
	return true
}

func (receiver *UserCumulativeStat) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToLazyChan(TbUserCumulativeStat, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserCumulativeStat) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToLazyChan(TbUserCumulativeStat, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initUserCumulativeStat() {
	syndb.RegLazyWithMiddle(TbUserCumulativeStat, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbUserCumulativeStat, db.UpdatedAtName)
	syndb.RegLazyWithMiddle(TbUserCumulativeStat, UserCumulativeStatTotalRecharge)
	syndb.RegLazyWithMiddle(TbUserCumulativeStat, UserCumulativeStatTotalWithdraw)
	syndb.RegLazyWithMiddle(TbUserCumulativeStat, UserCumulativeStatTotalFans)
	syndb.RegLazyWithMiddle(TbUserCumulativeStat, UserCumulativeStatTotalFollow)
	syndb.RegLazyWithMiddle(TbUserCumulativeStat, UserCumulativeStatTotalPayCount)
	syndb.RegLazyWithMiddle(TbUserCumulativeStat, UserCumulativeStatTotalDiamondConsume)
	syndb.RegLazyWithMiddle(TbUserCumulativeStat, UserCumulativeStatTotalGoldConsume)

	migrate.AutoMigrate(&UserCumulativeStat{})
}

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
	UserCumulativeStatTotalRechargeGold   db.TbCol = "total_recharge_gold"
	UserCumulativeStatTotalWithdraw       db.TbCol = "total_withdraw"
	UserCumulativeStatTotalPayCount       db.TbCol = "total_pay_count"
	UserCumulativeStatTotalDiamondConsume db.TbCol = "total_diamond_consume"
	UserCumulativeStatTotalGoldConsume    db.TbCol = "total_gold_consume"
)

// UserCumulativeStat 玩家累计数值(与用户一一对应,主键ID即用户ID)
type UserCumulativeStat struct {
	migrate.OneModel
	TotalRecharge       float64 `gorm:"default:0;comment:累计充值(USD)" json:"totalRecharge"`
	TotalRechargeGold   float64 `gorm:"default:0;comment:累计充值到账金币" json:"totalRechargeGold"`
	TotalWithdraw       float64 `gorm:"default:0;comment:累计提现" json:"totalWithdraw"`
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

func (receiver *UserCumulativeStat) AddTotalRecharge(val float64) bool {
	receiver.TotalRecharge = math.AddFloat64(receiver.TotalRecharge, val)

	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalRecharge, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalRecharge,
	})
	return true
}

func (receiver *UserCumulativeStat) AddTotalRechargeGold(val float64) bool {
	if val <= 0 {
		return false
	}
	receiver.TotalRechargeGold = math.AddFloat64(receiver.TotalRechargeGold, val)

	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalRechargeGold, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalRechargeGold,
	})
	return true
}

func (receiver *UserCumulativeStat) SetTotalRechargeGold(val float64) {
	receiver.TotalRechargeGold = val
	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalRechargeGold, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalRechargeGold,
	})
}

func (receiver *UserCumulativeStat) AddTotalWithdraw(val float64) bool {
	receiver.TotalWithdraw = math.AddFloat64(receiver.TotalWithdraw, val)

	syndb.AddData(TbUserCumulativeStat, UserCumulativeStatTotalWithdraw, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.TotalWithdraw,
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
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalRechargeGold)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalWithdraw)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalPayCount)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalDiamondConsume)
	syndb.RegLazy(TbUserCumulativeStat, UserCumulativeStatTotalGoldConsume)

	migrate.AutoMigrate(&UserCumulativeStat{})
}

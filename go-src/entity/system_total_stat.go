package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbSystemTotalStat db.TbName = "system_total_stats"
)

const (
	SystemTotalStatDefaultID uint64 = 1
)

const (
	SystemTotalStatTotalGold           db.TbCol = "total_gold"
	SystemTotalStatTotalGoldConsume    db.TbCol = "total_gold_consume"
	SystemTotalStatTotalDiamondConsume db.TbCol = "total_diamond_consume"
	SystemTotalStatTotalRecharge       db.TbCol = "total_recharge"
	SystemTotalStatTotalWithdraw       db.TbCol = "total_withdraw"
	SystemTotalStatTotalRegisterUser   db.TbCol = "total_register_user"
)

// SystemTotalStat 系统总数据(全局单条记录,默认ID=1)
type SystemTotalStat struct {
	migrate.OneModel
	TotalGold           float64 `gorm:"default:0;comment:金币总额" json:"totalGold"`
	TotalGoldConsume    float64 `gorm:"default:0;comment:金币总消费" json:"totalGoldConsume"`
	TotalDiamondConsume float64 `gorm:"default:0;comment:钻石总消费" json:"totalDiamondConsume"`
	TotalRecharge       float64 `gorm:"default:0;comment:总充值金额" json:"totalRecharge"`
	TotalWithdraw       float64 `gorm:"default:0;comment:总提现金额" json:"totalWithdraw"`
	TotalRegisterUser   uint64  `gorm:"default:0;comment:总注册用户数" json:"totalRegisterUser"`
}

// NewSystemTotalStat 构造系统总数据记录,字段写入通过 syndb lazy 异步入库
func NewSystemTotalStat(id uint64) *SystemTotalStat {
	if id == 0 {
		id = SystemTotalStatDefaultID
	}
	ret := &SystemTotalStat{}
	ret.ID = id
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

func (s *SystemTotalStat) AddTotalGold(val float64) {
	s.TotalGold = math.AddFloat64(s.TotalGold, val)

	syndb.AddDataToLazyChan(TbSystemTotalStat, SystemTotalStatTotalGold, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalGold,
	})
}

func (s *SystemTotalStat) AddTotalGoldConsume(val float64) {
	s.TotalGoldConsume = math.AddFloat64(s.TotalGoldConsume, val)

	syndb.AddDataToLazyChan(TbSystemTotalStat, SystemTotalStatTotalGoldConsume, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalGoldConsume,
	})
}

func (s *SystemTotalStat) AddTotalDiamondConsume(val float64) {
	s.TotalDiamondConsume = math.AddFloat64(s.TotalDiamondConsume, val)

	syndb.AddDataToLazyChan(TbSystemTotalStat, SystemTotalStatTotalDiamondConsume, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalDiamondConsume,
	})
}

func (s *SystemTotalStat) AddTotalRecharge(val float64) {
	s.TotalRecharge = math.AddFloat64(s.TotalRecharge, val)

	syndb.AddDataToLazyChan(TbSystemTotalStat, SystemTotalStatTotalRecharge, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalRecharge,
	})
}

func (s *SystemTotalStat) AddTotalWithdraw(val float64) {
	s.TotalWithdraw = math.AddFloat64(s.TotalWithdraw, val)

	syndb.AddDataToLazyChan(TbSystemTotalStat, SystemTotalStatTotalWithdraw, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalWithdraw,
	})
}

func (s *SystemTotalStat) AddTotalRegisterUser(val uint64) {
	s.TotalRegisterUser = math.Add(s.TotalRegisterUser, val)

	syndb.AddDataToLazyChan(TbSystemTotalStat, SystemTotalStatTotalRegisterUser, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalRegisterUser,
	})
}

func (s *SystemTotalStat) SetCreatedAt(v time.Time) {
	s.CreatedAt = v
	syndb.AddDataToLazyChan(TbSystemTotalStat, db.CreatedAtName, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: v,
	})
}

func (s *SystemTotalStat) SetUpdatedAt(v time.Time) {
	s.UpdatedAt = v
	syndb.AddDataToLazyChan(TbSystemTotalStat, db.UpdatedAtName, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: v,
	})
}

func (s *SystemTotalStat) touchUpdatedAt() {
	s.UpdatedAt = time.Now()
	syndb.AddDataToLazyChan(TbSystemTotalStat, db.UpdatedAtName, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.UpdatedAt,
	})
}

func initSystemTotalStat() {
	syndb.RegLazyWithMiddle(TbSystemTotalStat, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbSystemTotalStat, db.UpdatedAtName)
	syndb.RegLazyWithMiddle(TbSystemTotalStat, SystemTotalStatTotalGold)
	syndb.RegLazyWithMiddle(TbSystemTotalStat, SystemTotalStatTotalGoldConsume)
	syndb.RegLazyWithMiddle(TbSystemTotalStat, SystemTotalStatTotalDiamondConsume)
	syndb.RegLazyWithMiddle(TbSystemTotalStat, SystemTotalStatTotalRecharge)
	syndb.RegLazyWithMiddle(TbSystemTotalStat, SystemTotalStatTotalWithdraw)
	syndb.RegLazyWithMiddle(TbSystemTotalStat, SystemTotalStatTotalRegisterUser)
	migrate.AutoMigrate(&SystemTotalStat{})
}

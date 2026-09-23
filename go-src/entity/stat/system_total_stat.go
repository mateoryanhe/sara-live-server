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
	SystemTotalStatTotalGold                 db.TbCol = "total_gold"
	SystemTotalStatTotalGoldConsume          db.TbCol = "total_gold_consume"
	SystemTotalStatTotalDiamondConsume       db.TbCol = "total_diamond_consume"
	SystemTotalStatTotalRecharge             db.TbCol = "total_recharge"
	SystemTotalStatTotalNormalUserRecharge   db.TbCol = "total_normal_user_recharge"
	SystemTotalStatTotalCoinMerchantRecharge db.TbCol = "total_coin_merchant_recharge"
	SystemTotalStatTotalAnchorPayout         db.TbCol = "total_anchor_payout"
	SystemTotalStatTotalVirtualRecharge      db.TbCol = "total_virtual_recharge"
	SystemTotalStatTotalWithdraw             db.TbCol = "total_withdraw"
	SystemTotalStatTotalRegisterUser         db.TbCol = "total_register_user"
)

// SystemTotalStat 系统总数据(全局单条记录,默认ID=1)
type SystemTotalStat struct {
	migrate.OneModel
	TotalGold                 float64 `gorm:"default:0;comment:金币总额" json:"totalGold"`
	TotalGoldConsume          float64 `gorm:"default:0;comment:金币总消费" json:"totalGoldConsume"`
	TotalDiamondConsume       float64 `gorm:"default:0;comment:钻石总消费" json:"totalDiamondConsume"`
	TotalRecharge             float64 `gorm:"default:0;comment:全部美金入账累计(真实USD)" json:"totalRecharge"`
	TotalNormalUserRecharge   float64 `gorm:"default:0;comment:普通用户美金入账累计(真实USD)" json:"totalNormalUserRecharge"`
	TotalCoinMerchantRecharge float64 `gorm:"default:0;comment:币商美金入账累计(真实USD)" json:"totalCoinMerchantRecharge"`
	TotalAnchorPayout         float64 `gorm:"type:decimal(16,4);default:0;comment:主播代付累计(USD)" json:"totalAnchorPayout"`
	TotalVirtualRecharge      float64 `gorm:"default:0;comment:虚拟美金累计(充值白名单)" json:"totalVirtualRecharge"`
	TotalWithdraw             float64 `gorm:"default:0;comment:总提现金额" json:"totalWithdraw"`
	TotalRegisterUser         uint64  `gorm:"default:0;comment:总注册用户数" json:"totalRegisterUser"`
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

	syndb.AddData(TbSystemTotalStat, SystemTotalStatTotalGold, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalGold,
	})
}

func (s *SystemTotalStat) AddTotalGoldConsume(val float64) {
	s.TotalGoldConsume = math.AddFloat64(s.TotalGoldConsume, val)

	syndb.AddData(TbSystemTotalStat, SystemTotalStatTotalGoldConsume, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalGoldConsume,
	})
}

func (s *SystemTotalStat) AddTotalDiamondConsume(val float64) {
	s.TotalDiamondConsume = math.AddFloat64(s.TotalDiamondConsume, val)

	syndb.AddData(TbSystemTotalStat, SystemTotalStatTotalDiamondConsume, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalDiamondConsume,
	})
}

func (s *SystemTotalStat) AddTotalRecharge(val float64) {
	s.TotalRecharge = math.AddFloat64(s.TotalRecharge, val)

	syndb.AddData(TbSystemTotalStat, SystemTotalStatTotalRecharge, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalRecharge,
	})
}

func (s *SystemTotalStat) AddTotalNormalUserRecharge(val float64) {
	s.TotalNormalUserRecharge = math.AddFloat64(s.TotalNormalUserRecharge, val)

	syndb.AddData(TbSystemTotalStat, SystemTotalStatTotalNormalUserRecharge, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalNormalUserRecharge,
	})
}

func (s *SystemTotalStat) AddTotalCoinMerchantRecharge(val float64) {
	s.TotalCoinMerchantRecharge = math.AddFloat64(s.TotalCoinMerchantRecharge, val)

	syndb.AddData(TbSystemTotalStat, SystemTotalStatTotalCoinMerchantRecharge, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalCoinMerchantRecharge,
	})
}

func (s *SystemTotalStat) AddTotalAnchorPayout(val float64) {
	s.TotalAnchorPayout = math.AddFloat64(s.TotalAnchorPayout, val)

	syndb.AddData(TbSystemTotalStat, SystemTotalStatTotalAnchorPayout, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalAnchorPayout,
	})
}

func (s *SystemTotalStat) AddTotalVirtualRecharge(val float64) {
	s.TotalVirtualRecharge = math.AddFloat64(s.TotalVirtualRecharge, val)

	syndb.AddData(TbSystemTotalStat, SystemTotalStatTotalVirtualRecharge, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalVirtualRecharge,
	})
}

func (s *SystemTotalStat) AddTotalWithdraw(val float64) {
	s.TotalWithdraw = math.AddFloat64(s.TotalWithdraw, val)

	syndb.AddData(TbSystemTotalStat, SystemTotalStatTotalWithdraw, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalWithdraw,
	})
}

func (s *SystemTotalStat) AddTotalRegisterUser(val uint64) {
	s.TotalRegisterUser = math.Add(s.TotalRegisterUser, val)

	syndb.AddData(TbSystemTotalStat, SystemTotalStatTotalRegisterUser, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.TotalRegisterUser,
	})
}

func (s *SystemTotalStat) SetCreatedAt(v time.Time) {
	s.CreatedAt = v
	syndb.AddData(TbSystemTotalStat, db.CreatedAtName, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: v,
	})
}

func (s *SystemTotalStat) SetUpdatedAt(v time.Time) {
	s.UpdatedAt = v
	syndb.AddData(TbSystemTotalStat, db.UpdatedAtName, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: v,
	})
}

func (s *SystemTotalStat) touchUpdatedAt() {
	s.UpdatedAt = time.Now()
	syndb.AddData(TbSystemTotalStat, db.UpdatedAtName, &syndb.ColData{
		IdVal:  s.ID,
		ColVal: s.UpdatedAt,
	})
}

func initSystemTotalStat() {
	syndb.RegLazy(TbSystemTotalStat, db.CreatedAtName)
	syndb.RegLazy(TbSystemTotalStat, db.UpdatedAtName)
	syndb.RegLazy(TbSystemTotalStat, SystemTotalStatTotalGold)
	syndb.RegLazy(TbSystemTotalStat, SystemTotalStatTotalGoldConsume)
	syndb.RegLazy(TbSystemTotalStat, SystemTotalStatTotalDiamondConsume)
	syndb.RegLazy(TbSystemTotalStat, SystemTotalStatTotalRecharge)
	syndb.RegLazy(TbSystemTotalStat, SystemTotalStatTotalNormalUserRecharge)
	syndb.RegLazy(TbSystemTotalStat, SystemTotalStatTotalCoinMerchantRecharge)
	syndb.RegLazy(TbSystemTotalStat, SystemTotalStatTotalAnchorPayout)
	syndb.RegLazy(TbSystemTotalStat, SystemTotalStatTotalVirtualRecharge)
	syndb.RegLazy(TbSystemTotalStat, SystemTotalStatTotalWithdraw)
	syndb.RegLazy(TbSystemTotalStat, SystemTotalStatTotalRegisterUser)
	migrate.AutoMigrate(&SystemTotalStat{})
}

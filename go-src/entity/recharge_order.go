package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbRechargeOrder db.TbName = "recharge_orders"
)

// 充值订单状态
const (
	RechargeOrderStatusPending   uint8 = 0 // 待支付
	RechargeOrderStatusCompleted uint8 = 1 // 已完成(金币已发放)
	RechargeOrderStatusCancelled uint8 = 2 // 已取消
)

// 充值订单来源
const (
	RechargeOrderSourceApp    uint8 = 1 // App玩家发起
	RechargeOrderSourceManual uint8 = 2 // 后台手动充值
)

const (
	RechargeOrderUserId       db.TbCol = "user_id"
	RechargeOrderCfgId        db.TbCol = "cfg_id"
	RechargeOrderPrice        db.TbCol = "price"
	RechargeOrderCurrency     db.TbCol = "currency"
	RechargeOrderGold         db.TbCol = "gold"
	RechargeOrderStatus       db.TbCol = "status"
	RechargeOrderSource       db.TbCol = "source"
	RechargeOrderPayChannel   db.TbCol = "pay_channel"
	RechargeOrderThirdOrderId db.TbCol = "third_order_id"
	RechargeOrderRemark       db.TbCol = "remark"
	RechargeOrderOperatorId   db.TbCol = "operator_id"
	RechargeOrderPaidAt       db.TbCol = "paid_at"
)

// RechargeOrder 充值订单
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type RechargeOrder struct {
	migrate.OneModel
	UserId       uint64    `gorm:"index;default:0;comment:充值用户ID" json:"userId"`
	CfgId        uint64    `gorm:"index;default:0;comment:充值档位ID(0=手动输入金额)" json:"cfgId"`
	Price        float64   `gorm:"type:decimal(10,4);default:0;comment:实付金额(单位:USD)" json:"price"`
	Currency     string    `gorm:"size:8;default:'CNY';comment:货币(CNY/USD等)" json:"currency"`
	Gold         float64   `gorm:"default:0;comment:发放金币数(订单完成时增加到玩家金币)" json:"gold"`
	Status       uint8     `gorm:"index;default:0;comment:状态(0-待支付,1-已完成,2-已取消)" json:"status"`
	Source       uint8     `gorm:"default:0;comment:来源(1-App玩家发起,2-后台手动)" json:"source"`
	PayChannel   string    `gorm:"size:32;default:'';comment:支付渠道" json:"payChannel"`
	ThirdOrderId string    `gorm:"size:64;default:'';comment:第三方订单号" json:"thirdOrderId"`
	Remark       string    `gorm:"size:255;default:'';comment:备注" json:"remark"`
	OperatorId   uint64    `gorm:"default:0;comment:后台操作员ID(手动充值时记录)" json:"operatorId"`
	PaidAt       time.Time `gorm:"comment:完成支付时间" json:"paidAt"`
}

// NewRechargeOrder 构造一条新订单(ID 由 snowflake 生成,各字段通过 Setter 推入 syndb 缓冲)
func NewRechargeOrder(userId, cfgId uint64, price float64, currency string, gold float64, source uint8) *RechargeOrder {
	ret := &RechargeOrder{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetUserId(userId)
	ret.SetCfgId(cfgId)
	ret.SetPrice(price)
	ret.SetCurrency(currency)
	ret.SetGold(gold)
	ret.SetStatus(RechargeOrderStatusPending)
	ret.SetSource(source)
	return ret
}

func (r *RechargeOrder) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetCfgId(v uint64) {
	r.CfgId = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderCfgId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetPrice(v float64) {
	r.Price = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderPrice, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetCurrency(v string) {
	r.Currency = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderCurrency, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetGold(v float64) {
	r.Gold = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderGold, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetStatus(v uint8) {
	r.Status = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderStatus, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetSource(v uint8) {
	r.Source = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderSource, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetPayChannel(v string) {
	r.PayChannel = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderPayChannel, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetThirdOrderId(v string) {
	r.ThirdOrderId = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderThirdOrderId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetRemark(v string) {
	r.Remark = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderRemark, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetOperatorId(v uint64) {
	r.OperatorId = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderOperatorId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetPaidAt(v time.Time) {
	r.PaidAt = v
	syndb.AddDataToQuickChan(TbRechargeOrder, RechargeOrderPaidAt, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddDataToQuickChan(TbRechargeOrder, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *RechargeOrder) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddDataToQuickChan(TbRechargeOrder, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initRechargeOrder() {
	syndb.RegQuickWithLarge(TbRechargeOrder, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbRechargeOrder, db.UpdatedAtName)

	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderUserId)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderCfgId)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderPrice)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderCurrency)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderGold)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderStatus)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderSource)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderPayChannel)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderThirdOrderId)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderRemark)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderOperatorId)
	syndb.RegQuickWithLarge(TbRechargeOrder, RechargeOrderPaidAt)

	migrate.AutoMigrate(&RechargeOrder{})
}

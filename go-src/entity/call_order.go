package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbCallOrder db.TbName = "call_orders"
)

// 通话类型
const (
	CallOrderTypeVoice uint8 = 1 // 语音
	CallOrderTypeVideo uint8 = 2 // 视频
)

// 通话来源
const (
	CallOrderSourceLiveRoom       uint8 = 1 // 直播间
	CallOrderSourcePrivateMessage uint8 = 2 // 私信
)

const (
	CallOrderCallerId            db.TbCol = "caller_id"
	CallOrderReceiverId          db.TbCol = "receiver_id"
	CallOrderCallStartTime       db.TbCol = "call_start_time"
	CallOrderAnswerTime          db.TbCol = "answer_time"
	CallOrderCallerConfirmTime   db.TbCol = "caller_confirm_time"
	CallOrderReceiverConfirmTime db.TbCol = "receiver_confirm_time"
	CallOrderCallerHangUpTime    db.TbCol = "caller_hang_up_time"
	CallOrderReceiverHangUpTime  db.TbCol = "receiver_hang_up_time"
	CallOrderOrderEndTime        db.TbCol = "order_end_time"
	CallOrderCallDuration        db.TbCol = "call_duration"
	CallOrderCallType            db.TbCol = "call_type"
	CallOrderSource              db.TbCol = "source"
	CallOrderParams              db.TbCol = "params"
	CallOrderTicketPrice         db.TbCol = "ticket_price"
	CallOrderPricePerMinute      db.TbCol = "price_per_minute"
	CallOrderTotalCost           db.TbCol = "total_cost"
	CallOrderChargeTime          db.TbCol = "charge_time"
	CallOrderBillingDuration     db.TbCol = "billing_duration"
)

// CallOrder 通话订单
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type CallOrder struct {
	migrate.OneModel
	CallerId            uint64     `gorm:"index;default:0;comment:呼叫者ID" json:"callerId"`
	ReceiverId          uint64     `gorm:"index;default:0;comment:接收者ID" json:"receiverId"`
	CallStartTime       time.Time  `gorm:"index;comment:呼叫开始时间" json:"callStartTime"`
	AnswerTime          *time.Time `gorm:"index;comment:接听时间" json:"answerTime"`
	CallerConfirmTime   *time.Time `gorm:"index;comment:呼叫者确认时间" json:"callerConfirmTime"`
	ReceiverConfirmTime *time.Time `gorm:"index;comment:接听者确认时间" json:"receiverConfirmTime"`
	CallerHangUpTime    *time.Time `gorm:"index;comment:呼叫者挂断时间" json:"callerHangUpTime"`
	ReceiverHangUpTime  *time.Time `gorm:"index;comment:接听者挂断时间" json:"receiverHangUpTime"`
	OrderEndTime        *time.Time `gorm:"index;comment:订单结束时间" json:"orderEndTime"`
	CallDuration        uint32     `gorm:"default:0;comment:通话时长(秒)" json:"callDuration"`
	CallType            uint8      `gorm:"default:1;comment:通话类型(1-语音,2-视频)" json:"callType"`
	Source              uint8      `gorm:"default:1;comment:来源(1-直播间,2-私信)" json:"source"`
	Params              string     `gorm:"size:512;default:'';comment:扩展参数" json:"params"`
	TicketPrice         float64    `gorm:"type:decimal(10,4);default:0;comment:门票价格" json:"ticketPrice"`
	PricePerMinute      float64    `gorm:"type:decimal(10,4);default:0;comment:分钟计费价格(每分钟)" json:"pricePerMinute"`
	TotalCost           float64    `gorm:"type:decimal(10,4);default:0;comment:总费用" json:"totalCost"`
	ChargeTime          *time.Time `gorm:"index;comment:扣费时间" json:"chargeTime"`
	BillingDuration     uint32     `gorm:"default:0;comment:计费时长(分钟)" json:"billingDuration"`
}

func NewCallOrder(callerId, receiverId uint64, callType, source uint8, params string, ticketPrice, pricePerMinute float64) *CallOrder {
	ret := &CallOrder{}
	ret.ID = snowflake.GetId()
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetCallerId(callerId)
	ret.SetReceiverId(receiverId)
	ret.SetCallStartTime(now)
	ret.SetCallType(callType)
	ret.SetSource(source)
	ret.SetParams(params)
	ret.SetTicketPrice(ticketPrice)
	ret.SetPricePerMinute(pricePerMinute)
	return ret
}

func (m *CallOrder) HasEnded() bool {
	return m != nil && m.OrderEndTime != nil
}

func (m *CallOrder) IsCalling() bool {
	return m != nil && !m.HasEnded() && m.AnswerTime == nil
}

func (m *CallOrder) IsAccepted() bool {
	return m != nil && !m.HasEnded() && m.AnswerTime != nil
}

func (m *CallOrder) IsCallStarted() bool {
	return m.IsAccepted() && m.CallerConfirmTime != nil && m.ReceiverConfirmTime != nil
}

func (m *CallOrder) AnswerConfirmCount() uint32 {
	if m == nil {
		return 0
	}
	var count uint32
	if m.CallerConfirmTime != nil {
		count++
	}
	if m.ReceiverConfirmTime != nil {
		count++
	}
	return count
}

func (m *CallOrder) HasUserConfirmed(userId uint64) bool {
	if m == nil {
		return false
	}
	if userId == m.CallerId {
		return m.CallerConfirmTime != nil
	}
	if userId == m.ReceiverId {
		return m.ReceiverConfirmTime != nil
	}
	return false
}

func (m *CallOrder) HasPeerConfirmed(userId uint64) bool {
	if m == nil {
		return false
	}
	if userId == m.CallerId {
		return m.ReceiverConfirmTime != nil
	}
	if userId == m.ReceiverId {
		return m.CallerConfirmTime != nil
	}
	return false
}

func (m *CallOrder) PeerConfirmTime(userId uint64) *time.Time {
	if m == nil {
		return nil
	}
	if userId == m.CallerId {
		return m.ReceiverConfirmTime
	}
	if userId == m.ReceiverId {
		return m.CallerConfirmTime
	}
	return nil
}

func (m *CallOrder) SetUserConfirmTime(userId uint64, confirmTime time.Time) {
	if m == nil {
		return
	}
	t := confirmTime
	if userId == m.CallerId {
		m.SetCallerConfirmTime(&t)
		return
	}
	if userId == m.ReceiverId {
		m.SetReceiverConfirmTime(&t)
	}
}

func (m *CallOrder) SetCallerId(v uint64) {
	m.CallerId = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderCallerId, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetReceiverId(v uint64) {
	m.ReceiverId = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderReceiverId, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetCallStartTime(v time.Time) {
	m.CallStartTime = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderCallStartTime, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetAnswerTime(v *time.Time) {
	m.AnswerTime = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderAnswerTime, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetCallerConfirmTime(v *time.Time) {
	m.CallerConfirmTime = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderCallerConfirmTime, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetReceiverConfirmTime(v *time.Time) {
	m.ReceiverConfirmTime = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderReceiverConfirmTime, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetCallerHangUpTime(v *time.Time) {
	m.CallerHangUpTime = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderCallerHangUpTime, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetReceiverHangUpTime(v *time.Time) {
	m.ReceiverHangUpTime = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderReceiverHangUpTime, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetOrderEndTime(v *time.Time) {
	m.OrderEndTime = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderOrderEndTime, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetCallDuration(v uint32) {
	m.CallDuration = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderCallDuration, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetCallType(v uint8) {
	m.CallType = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderCallType, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetSource(v uint8) {
	m.Source = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderSource, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetParams(v string) {
	m.Params = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderParams, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetTicketPrice(v float64) {
	m.TicketPrice = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderTicketPrice, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetPricePerMinute(v float64) {
	m.PricePerMinute = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderPricePerMinute, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetTotalCost(v float64) {
	m.TotalCost = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderTotalCost, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetChargeTime(v *time.Time) {
	m.ChargeTime = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderChargeTime, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) AddBillingDuration(v uint32) {
	m.BillingDuration += v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderBillingDuration, &syndb.ColData{IdVal: m.ID, ColVal: m.BillingDuration})
}

func (m *CallOrder) SubBillingDuration(v uint32) {
	if v >= m.BillingDuration {
		m.BillingDuration = 0
	} else {
		m.BillingDuration -= v
	}
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderBillingDuration, &syndb.ColData{IdVal: m.ID, ColVal: m.BillingDuration})
}

func (m *CallOrder) SetCreatedAt(v time.Time) {
	m.CreatedAt = v
	syndb.AddDataToQuickChan(TbCallOrder, db.CreatedAtName, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallOrder) SetUpdatedAt(v time.Time) {
	m.UpdatedAt = v
	syndb.AddDataToQuickChan(TbCallOrder, db.UpdatedAtName, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func initCallOrder() {
	syndb.RegQuickWithMiddle(TbCallOrder, db.CreatedAtName)
	syndb.RegQuickWithMiddle(TbCallOrder, db.UpdatedAtName)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderCallerId)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderReceiverId)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderCallStartTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderAnswerTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderCallerConfirmTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderReceiverConfirmTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderCallerHangUpTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderReceiverHangUpTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderOrderEndTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderCallDuration)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderCallType)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderSource)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderParams)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderTicketPrice)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderPricePerMinute)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderTotalCost)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderChargeTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderBillingDuration)
	migrate.AutoMigrate(&CallOrder{})
}

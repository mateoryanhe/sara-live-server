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

// 通话状态
const (
	CallOrderStatusCalling  uint8 = 1 // 呼叫中
	CallOrderStatusInCall   uint8 = 2 // 通话中
	CallOrderStatusEnded    uint8 = 3 // 通话结束
	CallOrderStatusRejected uint8 = 4 // 拒接
	CallOrderStatusTimeout  uint8 = 5 // 呼叫超时
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
	CallOrderCallerId           db.TbCol = "caller_id"
	CallOrderReceiverId         db.TbCol = "receiver_id"
	CallOrderCallStartTime      db.TbCol = "call_start_time"
	CallOrderAnswerTime         db.TbCol = "answer_time"
	CallOrderCallerHangUpTime   db.TbCol = "caller_hang_up_time"
	CallOrderReceiverHangUpTime db.TbCol = "receiver_hang_up_time"
	CallOrderOrderEndTime       db.TbCol = "order_end_time"
	CallOrderCallDuration       db.TbCol = "call_duration"
	CallOrderCallStatus         db.TbCol = "call_status"
	CallOrderCallType           db.TbCol = "call_type"
	CallOrderSource             db.TbCol = "source"
	CallOrderParams             db.TbCol = "params"
	CallOrderTicketPrice        db.TbCol = "ticket_price"
	CallOrderPricePerMinute     db.TbCol = "price_per_minute"
	CallOrderTotalCost          db.TbCol = "total_cost"
)

// CallOrder 通话订单
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type CallOrder struct {
	migrate.OneModel
	CallerId           uint64     `gorm:"index;default:0;comment:呼叫者ID" json:"callerId"`
	ReceiverId         uint64     `gorm:"index;default:0;comment:接收者ID" json:"receiverId"`
	CallStartTime      time.Time  `gorm:"index;comment:呼叫开始时间" json:"callStartTime"`
	AnswerTime         *time.Time `gorm:"index;comment:接听时间" json:"answerTime"`
	CallerHangUpTime   *time.Time `gorm:"index;comment:呼叫者挂断时间" json:"callerHangUpTime"`
	ReceiverHangUpTime *time.Time `gorm:"index;comment:接听者挂断时间" json:"receiverHangUpTime"`
	OrderEndTime       *time.Time `gorm:"index;comment:订单结束时间" json:"orderEndTime"`
	CallDuration       uint32     `gorm:"default:0;comment:通话时长(秒)" json:"callDuration"`
	CallStatus         uint8      `gorm:"index;default:1;comment:通话状态(1-呼叫中,2-通话中,3-通话结束,4-拒接,5-呼叫超时)" json:"callStatus"`
	CallType           uint8      `gorm:"default:1;comment:通话类型(1-语音,2-视频)" json:"callType"`
	Source             uint8      `gorm:"default:1;comment:来源(1-直播间,2-私信)" json:"source"`
	Params             string     `gorm:"size:512;default:'';comment:扩展参数" json:"params"`
	TicketPrice        float64    `gorm:"type:decimal(10,4);default:0;comment:门票价格" json:"ticketPrice"`
	PricePerMinute     float64    `gorm:"type:decimal(10,4);default:0;comment:分钟计费价格(每分钟)" json:"pricePerMinute"`
	TotalCost          float64    `gorm:"type:decimal(10,4);default:0;comment:总费用" json:"totalCost"`
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
	ret.SetCallStatus(CallOrderStatusCalling)
	ret.SetCallType(callType)
	ret.SetSource(source)
	ret.SetParams(params)
	ret.SetTicketPrice(ticketPrice)
	ret.SetPricePerMinute(pricePerMinute)
	return ret
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

func (m *CallOrder) SetCallStatus(v uint8) {
	m.CallStatus = v
	syndb.AddDataToQuickChan(TbCallOrder, CallOrderCallStatus, &syndb.ColData{IdVal: m.ID, ColVal: v})
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
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderCallerHangUpTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderReceiverHangUpTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderOrderEndTime)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderCallDuration)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderCallStatus)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderCallType)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderSource)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderParams)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderTicketPrice)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderPricePerMinute)
	syndb.RegQuickWithMiddle(TbCallOrder, CallOrderTotalCost)
	migrate.AutoMigrate(&CallOrder{})
}

package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbCallUser db.TbName = "call_users"
)

const (
	CallUserCallOrderId db.TbCol = "call_order_id"
	CallUserHeartTime   db.TbCol = "heart_time"
)

// CallUser 通话用户
// 主键 ID 直接使用 userId,不重新生成
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type CallUser struct {
	ID          uint64     `gorm:"primaryKey;comment:用户ID" json:"id"`
	CallOrderId uint64     `gorm:"index;default:0;comment:通话订单ID" json:"callOrderId"`
	HeartTime   *time.Time `gorm:"index;comment:心跳时间" json:"heartTime"`
}

func NewCallUser(userId, callOrderId uint64) *CallUser {
	ret := &CallUser{}
	ret.ID = userId
	ret.SetCallOrderId(callOrderId)
	return ret
}

func (m *CallUser) SetCallOrderId(v uint64) {
	m.CallOrderId = v
	syndb.AddDataToQuickChan(TbCallUser, CallUserCallOrderId, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallUser) SetHeartTime(v *time.Time) {
	m.HeartTime = v
	syndb.AddDataToQuickChan(TbCallUser, CallUserHeartTime, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func initCallUser() {
	syndb.RegQuickWithMiddle(TbCallUser, CallUserCallOrderId)
	syndb.RegQuickWithMiddle(TbCallUser, CallUserHeartTime)
	migrate.AutoMigrate(&CallUser{})
}

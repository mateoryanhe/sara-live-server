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
	CallUserChannelName  db.TbCol = "channel_name"
	CallUserTokenGenTime db.TbCol = "token_gen_time"
	CallUserToken        db.TbCol = "token"
	CallUserCallOrderId  db.TbCol = "call_order_id"
)

// CallUser 通话用户
// 主键 ID 直接使用 userId,不重新生成
// 走 syndb 快速同步缓冲(quick chan),变更通过 Setter 推送到队列由 worker 周期性 Save 到 DB
type CallUser struct {
	ID           uint64    `gorm:"primaryKey;comment:用户ID" json:"id"`
	ChannelName  string    `gorm:"size:128;default:'';comment:频道名称" json:"channelName"`
	TokenGenTime time.Time `gorm:"index;comment:Token生成时间" json:"tokenGenTime"`
	Token        string    `gorm:"size:512;default:'';comment:Token字符串" json:"token"`
	CallOrderId  uint64    `gorm:"index;default:0;comment:通话订单ID" json:"callOrderId"`
}

func NewCallUser(userId uint64, channelName, token string, callOrderId uint64) *CallUser {
	ret := &CallUser{}
	ret.ID = userId
	ret.SetChannelName(channelName)
	ret.SetTokenGenTime(time.Now())
	ret.SetToken(token)
	ret.SetCallOrderId(callOrderId)
	return ret
}

func (m *CallUser) SetChannelName(v string) {
	m.ChannelName = v
	syndb.AddDataToQuickChan(TbCallUser, CallUserChannelName, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallUser) SetTokenGenTime(v time.Time) {
	m.TokenGenTime = v
	syndb.AddDataToQuickChan(TbCallUser, CallUserTokenGenTime, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallUser) SetToken(v string) {
	m.Token = v
	syndb.AddDataToQuickChan(TbCallUser, CallUserToken, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func (m *CallUser) SetCallOrderId(v uint64) {
	m.CallOrderId = v
	syndb.AddDataToQuickChan(TbCallUser, CallUserCallOrderId, &syndb.ColData{IdVal: m.ID, ColVal: v})
}

func initCallUser() {
	syndb.RegQuickWithMiddle(TbCallUser, CallUserChannelName)
	syndb.RegQuickWithMiddle(TbCallUser, CallUserTokenGenTime)
	syndb.RegQuickWithLarge(TbCallUser, CallUserToken)
	syndb.RegQuickWithMiddle(TbCallUser, CallUserCallOrderId)
	migrate.AutoMigrate(&CallUser{})
}

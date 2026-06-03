package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRoomOnline db.TbName = "live_room_onlines"
)

const (
	LiveRoomOnlineRoomId      db.TbCol = "room_id"
	LiveRoomOnlineUserId      db.TbCol = "user_id"
	LiveRoomOnlineStatus      db.TbCol = "status"
	TbLiveRoomOnlineHeartTime db.TbCol = "heart_time"
)

// 在线状态
const (
	LiveRoomOnlineStatusOffline uint8 = 0 // 已离开
	LiveRoomOnlineStatusOnline  uint8 = 1 // 在线中
)

// LiveRoomOnline 直播间在线玩家记录
// 主键 ID 使用复合字符串: "{userId}_{roomId}",方便唯一定位以及前端拼接
// 玩家加入/离开通过 Status 字段切换,记录长期保留
type LiveRoomOnline struct {
	ID        string     `gorm:"primaryKey;size:64;comment:复合ID(userId_roomId)" json:"id"`
	RoomId    uint64     `gorm:"index:idx_room_status,priority:1;default:0;comment:直播间ID" json:"roomId"`
	UserId    uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	Status    uint8      `gorm:"index:idx_room_status,priority:2;default:0;comment:状态(0已离开,1在线)" json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"-"`
	HeartTime *time.Time `gorm:"comment:房间心跳状态,大于5分钟，判断下播" json:"heart_time"`
}

// BuildLiveRoomOnlineId 拼接复合主键
func BuildLiveRoomOnlineId(userId, roomId uint64) string {
	return fmt.Sprintf("%d_%d", userId, roomId)
}

// NewLiveRoomOnline 构造内存对象,字段写入通过 syndb 异步入库
func NewLiveRoomOnline(userId, roomId uint64) *LiveRoomOnline {
	o := &LiveRoomOnline{}
	o.ID = BuildLiveRoomOnlineId(userId, roomId)
	now := time.Now()
	o.SetCreatedAt(now)
	o.SetUpdatedAt(now)
	o.SetUserId(userId)
	o.SetRoomId(roomId)
	o.SetStatus(LiveRoomOnlineStatusOnline)
	return o
}

func (o *LiveRoomOnline) SetRoomId(v uint64) {
	o.RoomId = v
	syndb.AddDataToQuickChan(TbLiveRoomOnline, LiveRoomOnlineRoomId, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})
}

func (o *LiveRoomOnline) SetUserId(v uint64) {
	o.UserId = v
	syndb.AddDataToQuickChan(TbLiveRoomOnline, LiveRoomOnlineUserId, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})
}

func (o *LiveRoomOnline) SetStatus(v uint8) {
	o.Status = v
	o.UpdatedAt = time.Now()
	syndb.AddDataToQuickChan(TbLiveRoomOnline, LiveRoomOnlineStatus, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})
	syndb.AddDataToQuickChan(TbLiveRoomOnline, db.UpdatedAtName, &syndb.ColData{
		IdVal: o.ID, ColVal: o.UpdatedAt,
	})
}

func (o *LiveRoomOnline) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
	syndb.AddDataToQuickChan(TbLiveRoomOnline, db.CreatedAtName, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})
}

func (o *LiveRoomOnline) SetUpdatedAt(v time.Time) {
	o.UpdatedAt = v
	syndb.AddDataToQuickChan(TbLiveRoomOnline, db.UpdatedAtName, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})
}

func (o *LiveRoomOnline) SetHeartTime(v *time.Time) {
	o.HeartTime = v
	syndb.AddDataToLazyChan(TbLiveRoomOnline, TbLiveRoomOnlineHeartTime, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})
}

func initLiveRoomOnline() {
	syndb.RegQuickWithLarge(TbLiveRoomOnline, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbLiveRoomOnline, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbLiveRoomOnline, LiveRoomOnlineRoomId)
	syndb.RegQuickWithLarge(TbLiveRoomOnline, LiveRoomOnlineUserId)
	syndb.RegQuickWithLarge(TbLiveRoomOnline, LiveRoomOnlineStatus)

	syndb.RegLazyWithLarge(TbLiveRoomOnline, TbLiveRoomOnlineHeartTime)

	migrate.AutoMigrate(&LiveRoomOnline{})
}

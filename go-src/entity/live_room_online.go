package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
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
	LiveRoomOnlineMuted       db.TbCol = "muted"
	TbLiveRoomOnlineHeartTime db.TbCol = "heart_time"
	TbLiveRoomOnlineJoinTime  db.TbCol = "join_time"
	TbLiveRoomOnlineKickTime  db.TbCol = "kick_time"
	LiveRoomOnlineTotalReward db.TbCol = "total_reward"
)

const LiveRoomKickBanDuration = 30 * time.Minute

// 在线状态
const (
	LiveRoomOnlineStatusOffline uint8 = 0 // 已离开
	LiveRoomOnlineStatusOnline  uint8 = 1 // 在线中
)

// LiveRoomOnline 直播间在线玩家记录
// 主键 ID 使用复合字符串: "{userId}_{roomId}",方便唯一定位以及前端拼接
// 玩家加入/离开通过 Status 字段切换,记录长期保留
type LiveRoomOnline struct {
	ID          string     `gorm:"primaryKey;size:64;comment:复合ID(userId_roomId)" json:"id"`
	RoomId      uint64     `gorm:"index:idx_room_status,priority:1;default:0;comment:直播间ID" json:"roomId"`
	UserId      uint64     `gorm:"index;default:0;comment:用户ID" json:"userId"`
	Status      uint8      `gorm:"index:idx_room_status,priority:2;default:0;comment:状态(0已离开,1在线)" json:"status"`
	Muted       bool       `gorm:"default:0;comment:禁言标记" json:"muted"`
	CreatedAt   time.Time  `json:"createdAt"`
	JoinTime    *time.Time `gorm:"comment:进入直播间时间" json:"join_time"`
	HeartTime   *time.Time `gorm:"comment:房间心跳状态,大于5分钟，判断下播" json:"heart_time"`
	KickTime    *time.Time `gorm:"comment:被主播踢出时间" json:"kickTime"`
	TotalReward float64    `gorm:"type:decimal(16,4);default:0;comment:打赏累计(钻石)" json:"totalReward"`
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
	syndb.AddDataToQuickChan(TbLiveRoomOnline, LiveRoomOnlineStatus, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})

}

func (o *LiveRoomOnline) SetMuted(v bool) {
	o.Muted = v

	syndb.AddDataToQuickChan(TbLiveRoomOnline, LiveRoomOnlineMuted, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})

}

func (o *LiveRoomOnline) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
	syndb.AddDataToQuickChan(TbLiveRoomOnline, db.CreatedAtName, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})
}

func (o *LiveRoomOnline) SetJoinTime(v *time.Time) {
	o.JoinTime = v
	syndb.AddDataToQuickChan(TbLiveRoomOnline, TbLiveRoomOnlineJoinTime, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})
}

func (o *LiveRoomOnline) SetHeartTime(v *time.Time) {
	o.HeartTime = v
	syndb.AddDataToLazyChan(TbLiveRoomOnline, TbLiveRoomOnlineHeartTime, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})
}

func (o *LiveRoomOnline) SetKickTime(v *time.Time) {
	o.KickTime = v
	syndb.AddDataToQuickChan(TbLiveRoomOnline, TbLiveRoomOnlineKickTime, &syndb.ColData{
		IdVal: o.ID, ColVal: v,
	})
}

func (o *LiveRoomOnline) AddTotalReward(amount float64) {
	if amount <= 0 {
		return
	}
	o.TotalReward = math.AddFloat64(o.TotalReward, amount)
	syndb.AddDataToLazyChan(TbLiveRoomOnline, LiveRoomOnlineTotalReward, &syndb.ColData{
		IdVal: o.ID, ColVal: o.TotalReward,
	})
}

// IsKickBanned 是否在踢出封禁期内
func (o *LiveRoomOnline) IsKickBanned() bool {
	if o == nil || o.KickTime == nil {
		return false
	}
	return time.Now().Before(o.KickTime.Add(LiveRoomKickBanDuration))
}

func initLiveRoomOnline() {
	syndb.RegQuickWithLarge(TbLiveRoomOnline, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbLiveRoomOnline, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbLiveRoomOnline, LiveRoomOnlineRoomId)
	syndb.RegQuickWithLarge(TbLiveRoomOnline, LiveRoomOnlineUserId)
	syndb.RegQuickWithLarge(TbLiveRoomOnline, LiveRoomOnlineStatus)
	syndb.RegQuickWithLarge(TbLiveRoomOnline, LiveRoomOnlineMuted)
	syndb.RegQuickWithMiddle(TbLiveRoomOnline, TbLiveRoomOnlineJoinTime)
	syndb.RegQuickWithLarge(TbLiveRoomOnline, TbLiveRoomOnlineKickTime)
	syndb.RegLazyWithLarge(TbLiveRoomOnline, LiveRoomOnlineTotalReward)

	syndb.RegLazyWithLarge(TbLiveRoomOnline, TbLiveRoomOnlineHeartTime)

	migrate.AutoMigrate(&LiveRoomOnline{})
}

package entity

import (
	"fmt"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveFollow db.TbName = "live_follows"
)

const (
	LiveFollowUserId   db.TbCol = "user_id"
	LiveFollowAnchorId db.TbCol = "anchor_id"
	LiveFollowStatus   db.TbCol = "status"
)

// 关注/拉黑状态
const (
	LiveFollowStatusUnfollow uint8 = 0 // 已取消关注
	LiveFollowStatusFollow   uint8 = 1 // 已关注
	LiveFollowStatusBlock    uint8 = 2 // 已拉黑
)

// LiveFollow 关注主播记录
// 主键 ID 使用复合字符串: "{userId}_{anchorId}",同一对粉丝/主播仅一行,通过 Status 切换
type LiveFollow struct {
	ID        string    `gorm:"primaryKey;size:64;comment:复合ID(userId_anchorId)" json:"id"`
	UserId    uint64    `gorm:"index:idx_user_status,priority:1;default:0;comment:粉丝用户ID" json:"userId"`
	AnchorId  uint64    `gorm:"index:idx_anchor_status,priority:1;default:0;comment:被关注主播ID" json:"anchorId"`
	Status    uint8     `gorm:"index:idx_user_status,priority:2;index:idx_anchor_status,priority:2;default:0;comment:状态(0已取消,1已关注,2已拉黑)" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BuildLiveFollowId 拼接复合主键
func BuildLiveFollowId(userId, anchorId uint64) string {
	return fmt.Sprintf("%d_%d", userId, anchorId)
}

// NewLiveFollow 构造关注记录,字段通过 syndb 异步入库
func NewLiveFollow(userId, anchorId uint64) *LiveFollow {
	return NewLiveFollowWithStatus(userId, anchorId, LiveFollowStatusFollow)
}

// NewLiveFollowWithStatus 构造指定状态的 live_follow 记录
func NewLiveFollowWithStatus(userId, anchorId uint64, status uint8) *LiveFollow {
	f := &LiveFollow{}
	f.ID = BuildLiveFollowId(userId, anchorId)
	now := time.Now()
	f.SetCreatedAt(now)
	f.SetUpdatedAt(now)
	f.SetUserId(userId)
	f.SetAnchorId(anchorId)
	f.SetStatus(status)
	return f
}

func (f *LiveFollow) SetUserId(v uint64) {
	f.UserId = v
	syndb.AddDataToQuickChan(TbLiveFollow, LiveFollowUserId, &syndb.ColData{
		IdVal: f.ID, ColVal: v,
	})
}

func (f *LiveFollow) SetAnchorId(v uint64) {
	f.AnchorId = v
	syndb.AddDataToQuickChan(TbLiveFollow, LiveFollowAnchorId, &syndb.ColData{
		IdVal: f.ID, ColVal: v,
	})
}

func (f *LiveFollow) SetStatus(v uint8) {
	f.Status = v
	f.UpdatedAt = time.Now()
	syndb.AddDataToQuickChan(TbLiveFollow, LiveFollowStatus, &syndb.ColData{
		IdVal: f.ID, ColVal: v,
	})
	syndb.AddDataToQuickChan(TbLiveFollow, db.UpdatedAtName, &syndb.ColData{
		IdVal: f.ID, ColVal: f.UpdatedAt,
	})
}

func (f *LiveFollow) SetCreatedAt(v time.Time) {
	f.CreatedAt = v
	syndb.AddDataToQuickChan(TbLiveFollow, db.CreatedAtName, &syndb.ColData{
		IdVal: f.ID, ColVal: v,
	})
}

func (f *LiveFollow) SetUpdatedAt(v time.Time) {
	f.UpdatedAt = v
	syndb.AddDataToQuickChan(TbLiveFollow, db.UpdatedAtName, &syndb.ColData{
		IdVal: f.ID, ColVal: v,
	})
}

func initLiveFollow() {
	syndb.RegQuick(TbLiveFollow, db.CreatedAtName)
	syndb.RegQuick(TbLiveFollow, db.UpdatedAtName)
	syndb.RegQuick(TbLiveFollow, LiveFollowUserId)
	syndb.RegQuick(TbLiveFollow, LiveFollowAnchorId)
	syndb.RegQuick(TbLiveFollow, LiveFollowStatus)
	migrate.AutoMigrate(&LiveFollow{})
}

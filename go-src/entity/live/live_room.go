package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRoom db.TbName = "live_rooms"
)

const (
	LiveRoomGuildId      db.TbCol = "guild_id"
	LiveRoomTitle        db.TbCol = "title"
	LiveRoomCover        db.TbCol = "cover"
	LiveRoomNotice       db.TbCol = "notice"
	LiveRoomLiveId       db.TbCol = "live_record_id"
	LiveRoomHeartTime    db.TbCol = "heart_time"
	LiveRoomBan          db.TbCol = "ban"
	LiveRoomBanApplyTime db.TbCol = "ban_apply_time"
	LiveRoomBanReason    db.TbCol = "ban_reason"
	LiveRoomStatus       db.TbCol = "status"
)

const (
	LiveRoomCategoryHot     uint8 = 1 // hot
	LiveRoomCategoryGame    uint8 = 2 // game
	LiveRoomCategoryPrivate uint8 = 3 // 私密
)

// 直播间上下架状态
const (
	LiveRoomStatusOffShelf uint8 = 0 // 下架
	LiveRoomStatusOnShelf  uint8 = 1 // 上架
)

const (
	LiveRoomPrivateInviteAll    uint8 = 1 // 接受所有人
	LiveRoomPrivateInviteReject uint8 = 3 // 拒绝所有人

	LiveRoomCallMinTotalRechargeUSD = 10 // 直播间1v1通话最低累计充值(USD)
)

// NormalizePrivateInviteType 归一化私密邀请类型
func NormalizePrivateInviteType(v, category uint8) uint8 {
	if v == 0 {
		v = DefaultPrivateInviteType(category)
	}
	if v != LiveRoomPrivateInviteAll && v != LiveRoomPrivateInviteReject {
		return LiveRoomPrivateInviteAll
	}
	return v
}

// DefaultPrivateInviteType 按直播间分类返回私密邀请类型默认值
func DefaultPrivateInviteType(category uint8) uint8 {
	if category == LiveRoomCategoryHot {
		return LiveRoomPrivateInviteReject
	}
	return LiveRoomPrivateInviteAll
}

// LiveRoom 直播间(LiveRoom.ID 与 UserInfo.ID 均为主播用户ID,每个主播仅一个直播间)
type LiveRoom struct {
	migrate.OneModel
	GuildId      uint64     `gorm:"index;default:0;comment:所属工会ID" json:"guildId"`
	Title        string     `gorm:"size:128;default:'';comment:直播间标题" json:"title"`
	Cover        string     `gorm:"size:255;default:'';comment:封面图URL" json:"cover"`
	Notice       string     `gorm:"size:512;default:'';comment:公告" json:"notice"`
	LiveRecordId uint64     `gorm:"default:0;comment:直播记录id" json:"liveRecordId"`
	HeartTime    *time.Time `gorm:"comment:房间心跳状态,大于5分钟，判断下播" json:"heart_time"`
	Ban          bool       `gorm:"default:0;comment:封禁状态" json:"ban"`
	BanApplyTime *time.Time `gorm:"comment:封禁截止时间" json:"banApplyTime"`
	BanReason    string     `gorm:"size:512;default:'';comment:封禁原因" json:"banReason"`
	Status       uint8      `gorm:"default:1;comment:状态(0-下架,1-上架)" json:"status"`
}

// NewLiveRoom 构造内存对象,字段写入通过 syndb 异步入库
// anchorId 同时作为 LiveRoom 的主键 ID
func NewLiveRoom(anchorId, guildId uint64, title, cover, notice string) *LiveRoom {
	r := &LiveRoom{}
	r.ID = anchorId
	now := time.Now()
	r.SetCreatedAt(now)
	r.SetUpdatedAt(now)
	r.SetGuildId(guildId)
	r.SetTitle(title)
	r.SetCover(cover)
	r.SetNotice(notice)
	r.SetStatus(LiveRoomStatusOnShelf)
	return r
}

func (r *LiveRoom) SetGuildId(v uint64) {
	r.GuildId = v
	syndb.AddData(TbLiveRoom, LiveRoomGuildId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetTitle(v string) {
	r.Title = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomTitle, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetHeartTime(v *time.Time) {
	r.HeartTime = v
	syndb.AddData(TbLiveRoom, LiveRoomHeartTime, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetCover(v string) {
	r.Cover = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomCover, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetLiveRecordId(v uint64) {
	r.LiveRecordId = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomLiveId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetNotice(v string) {
	r.Notice = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomNotice, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetBan(v bool) {
	r.Ban = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomBan, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetBanApplyTime(v *time.Time) {
	r.BanApplyTime = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomBanApplyTime, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetBanReason(v string) {
	r.BanReason = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomBanReason, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetStatus(v uint8) {
	if v != LiveRoomStatusOffShelf && v != LiveRoomStatusOnShelf {
		v = LiveRoomStatusOnShelf
	}
	r.Status = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomStatus, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbLiveRoom, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbLiveRoom, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	syndb.AddData(TbLiveRoom, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: r.UpdatedAt,
	})
}

func initLiveRoom() {
	syndb.RegQuick(TbLiveRoom, db.CreatedAtName)
	syndb.RegQuick(TbLiveRoom, db.UpdatedAtName)
	syndb.RegQuick(TbLiveRoom, LiveRoomGuildId)
	syndb.RegQuick(TbLiveRoom, LiveRoomTitle)
	syndb.RegQuick(TbLiveRoom, LiveRoomCover)
	syndb.RegQuick(TbLiveRoom, LiveRoomNotice)
	syndb.RegQuick(TbLiveRoom, LiveRoomLiveId)
	syndb.RegQuick(TbLiveRoom, LiveRoomBan)
	syndb.RegQuick(TbLiveRoom, LiveRoomBanApplyTime)
	syndb.RegQuick(TbLiveRoom, LiveRoomBanReason)
	syndb.RegLazy(TbLiveRoom, LiveRoomHeartTime)
	migrate.AutoMigrate(&LiveRoom{})
}

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
	LiveRoomGuildId                  db.TbCol = "guild_id"
	LiveRoomTitle                    db.TbCol = "title"
	LiveRoomCover                    db.TbCol = "cover"
	LiveRoomNotice                   db.TbCol = "notice"
	LiveRoomLiveId                   db.TbCol = "live_record_id"
	LiveRoomHeartTime                db.TbCol = "heart_time"
	LiveRoomBan                      db.TbCol = "ban"
	LiveRoomBanApplyTime             db.TbCol = "ban_apply_time"
	LiveRoomPrivateInviteType        db.TbCol = "private_invite_type"
	LiveRoomBanReason                db.TbCol = "ban_reason"
	LiveRoomCategory                 db.TbCol = "category"
	LiveRoomTagId                    db.TbCol = "tag_id"
	LiveRoomTicket                   db.TbCol = "ticket"
	LiveRoomBilling                  db.TbCol = "billing"
	LiveRoomCloudPlayerVideo         db.TbCol = "cloud_player_video"
	LiveRoomPushStream               db.TbCol = "push_stream"
	LiveRoomIsTest                   db.TbCol = "is_test"
	LiveRoomCloudPlayerId            db.TbCol = "cloud_player_id"
	LiveRoomCloudPlayerTokenExpireAt db.TbCol = "cloud_player_token_expire_at"
	LiveRoomStatus                   db.TbCol = "status"
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
	LiveRoomPrivateInviteVip    uint8 = 2 // 仅VIP
	LiveRoomPrivateInviteReject uint8 = 3 // 拒绝所有人
)

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
	GuildId                  uint64     `gorm:"index;default:0;comment:所属工会ID" json:"guildId"`
	Title                    string     `gorm:"size:128;default:'';comment:直播间标题" json:"title"`
	Cover                    string     `gorm:"size:255;default:'';comment:封面图URL" json:"cover"`
	Notice                   string     `gorm:"size:512;default:'';comment:公告" json:"notice"`
	LiveRecordId             uint64     `gorm:"default:0;comment:直播记录id" json:"liveRecordId"`
	HeartTime                *time.Time `gorm:"comment:房间心跳状态,大于5分钟，判断下播" json:"heart_time"`
	Ban                      bool       `gorm:"default:0;comment:封禁状态" json:"ban"`
	BanApplyTime             *time.Time `gorm:"comment:封禁截止时间" json:"banApplyTime"`
	PrivateInviteType        uint8      `gorm:"default:1;comment:私密邀请类型(1=接受所有人,2=仅VIP,3=拒绝所有人)" json:"privateInviteType"`
	BanReason                string     `gorm:"size:512;default:'';comment:封禁原因" json:"banReason"`
	Category                 uint8      `gorm:"default:1;comment:分类(1=hot,2=game,3=私密)" json:"category"`
	TagId                    uint64     `gorm:"default:0;comment:直播间标签ID" json:"tagId"`
	Ticket                   float64    `gorm:"type:decimal(10,4);default:0;comment:门票价格(钻石)" json:"ticket"`
	Billing                  float64    `gorm:"type:decimal(10,4);default:0;comment:计费价格(每分钟钻石)" json:"billing"`
	CloudPlayerVideo         string     `gorm:"size:512;default:'';comment:云播放器MP4视频URL/路径" json:"cloudPlayerVideo"`
	PushStream               bool       `gorm:"default:0;comment:是否推流" json:"pushStream"`
	IsTest                   bool       `gorm:"default:0;comment:是否测试机器人主播(仅下发App)" json:"isTest"`
	CloudPlayerId            string     `gorm:"size:64;default:'';comment:声网云播放器ID" json:"cloudPlayerId"`
	CloudPlayerTokenExpireAt *time.Time `gorm:"comment:云播放器RTC token过期时间" json:"cloudPlayerTokenExpireAt"`
	Status                   uint8      `gorm:"default:1;comment:状态(0-下架,1-上架)" json:"status"`
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
	r.SetCategory(LiveRoomCategoryHot)
	r.SetPrivateInviteType(DefaultPrivateInviteType(LiveRoomCategoryHot))
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

func (r *LiveRoom) SetPrivateInviteType(v uint8) {
	if v != LiveRoomPrivateInviteAll && v != LiveRoomPrivateInviteVip && v != LiveRoomPrivateInviteReject {
		v = LiveRoomPrivateInviteAll
	}
	r.PrivateInviteType = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomPrivateInviteType, &syndb.ColData{
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

func (r *LiveRoom) SetCategory(v uint8) {
	if v != LiveRoomCategoryHot && v != LiveRoomCategoryGame && v != LiveRoomCategoryPrivate {
		v = LiveRoomCategoryHot
	}
	r.Category = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomCategory, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetTagId(v uint64) {
	r.TagId = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomTagId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetTicket(v float64) {
	r.Ticket = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomTicket, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetBilling(v float64) {
	r.Billing = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomBilling, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetCloudPlayerVideo(v string) {
	r.CloudPlayerVideo = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomCloudPlayerVideo, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetPushStream(v bool) {
	r.PushStream = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomPushStream, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetIsTest(v bool) {
	r.IsTest = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomIsTest, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetCloudPlayerId(v string) {
	r.CloudPlayerId = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomCloudPlayerId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *LiveRoom) SetCloudPlayerTokenExpireAt(v *time.Time) {
	r.CloudPlayerTokenExpireAt = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoom, LiveRoomCloudPlayerTokenExpireAt, &syndb.ColData{
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
	syndb.RegQuick(TbLiveRoom, LiveRoomPrivateInviteType)
	syndb.RegQuick(TbLiveRoom, LiveRoomBanReason)
	syndb.RegQuick(TbLiveRoom, LiveRoomCategory)
	syndb.RegQuick(TbLiveRoom, LiveRoomTagId)
	syndb.RegQuick(TbLiveRoom, LiveRoomTicket)
	syndb.RegQuick(TbLiveRoom, LiveRoomBilling)
	syndb.RegQuick(TbLiveRoom, LiveRoomCloudPlayerVideo)
	syndb.RegQuick(TbLiveRoom, LiveRoomPushStream)
	syndb.RegQuick(TbLiveRoom, LiveRoomIsTest)
	syndb.RegQuick(TbLiveRoom, LiveRoomCloudPlayerId)
	syndb.RegQuick(TbLiveRoom, LiveRoomCloudPlayerTokenExpireAt)

	syndb.RegLazy(TbLiveRoom, LiveRoomHeartTime)

	migrate.AutoMigrate(&LiveRoom{})
}

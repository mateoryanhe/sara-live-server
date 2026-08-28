package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbLiveRoomCfg db.TbName = "live_room_cfgs"
)

const (
	LiveRoomCfgPrivateInviteType        db.TbCol = "private_invite_type"
	LiveRoomCfgCategory                 db.TbCol = "category"
	LiveRoomCfgTagId                    db.TbCol = "tag_id"
	LiveRoomCfgTicket                   db.TbCol = "ticket"
	LiveRoomCfgBilling                  db.TbCol = "billing"
	LiveRoomCfgCloudPlayerVideo         db.TbCol = "cloud_player_video"
	LiveRoomCfgCloudPlayerFrameRate     db.TbCol = "cloud_player_frame_rate"
	LiveRoomCfgCloudPlayerBitrate       db.TbCol = "cloud_player_bitrate"
	LiveRoomCfgCloudPlayerWidth         db.TbCol = "cloud_player_width"
	LiveRoomCfgCloudPlayerHeight        db.TbCol = "cloud_player_height"
	LiveRoomCfgPushStream               db.TbCol = "push_stream"
	LiveRoomCfgIsTest                   db.TbCol = "is_test"
	LiveRoomCfgCloudPlayerId            db.TbCol = "cloud_player_id"
	LiveRoomCfgCloudPlayerTokenExpireAt db.TbCol = "cloud_player_token_expire_at"
)

const (
	DefaultCloudPlayerFrameRate = 24
	DefaultCloudPlayerBitrate   = 400
	DefaultCloudPlayerWidth     = 960
	DefaultCloudPlayerHeight    = 540
)

// LiveRoomCfg 直播间配置(主键=房间ID=主播ID)
type LiveRoomCfg struct {
	migrate.OneModel
	PrivateInviteType        uint8      `gorm:"default:1;comment:私密邀请类型(1=接受所有人,3=拒绝所有人)" json:"privateInviteType"`
	Category                 uint8      `gorm:"default:1;comment:分类(1=hot,2=game,3=私密)" json:"category"`
	TagId                    uint64     `gorm:"default:0;comment:直播间标签ID" json:"tagId"`
	Ticket                   float64    `gorm:"type:decimal(10,4);default:0;comment:门票价格(钻石)" json:"ticket"`
	Billing                  float64    `gorm:"type:decimal(10,4);default:0;comment:计费价格(每分钟钻石)" json:"billing"`
	CloudPlayerVideo         string     `gorm:"size:512;default:'';comment:云播放器MP4视频URL/路径" json:"cloudPlayerVideo"`
	CloudPlayerFrameRate     uint8      `gorm:"default:24;comment:云播放器输出帧率(fps)" json:"cloudPlayerFrameRate"`
	CloudPlayerBitrate       int        `gorm:"default:400;comment:云播放器输出码率(Kbps)" json:"cloudPlayerBitrate"`
	CloudPlayerWidth         int        `gorm:"default:960;comment:云播放器输出宽度(px)" json:"cloudPlayerWidth"`
	CloudPlayerHeight        int        `gorm:"default:540;comment:云播放器输出高度(px)" json:"cloudPlayerHeight"`
	PushStream               bool       `gorm:"default:0;comment:是否推流" json:"pushStream"`
	IsTest                   bool       `gorm:"default:0;comment:是否测试机器人主播(仅下发App)" json:"isTest"`
	CloudPlayerId            string     `gorm:"size:64;default:'';comment:声网云播放器ID" json:"cloudPlayerId"`
	CloudPlayerTokenExpireAt *time.Time `gorm:"comment:云播放器RTC token过期时间" json:"cloudPlayerTokenExpireAt"`
}

// NewLiveRoomCfg 构造直播间配置并入库(默认 hot + 对应邀请类型)
func NewLiveRoomCfg(roomId uint64) *LiveRoomCfg {
	r := &LiveRoomCfg{}
	r.ID = roomId
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	syndb.AddData(TbLiveRoomCfg, db.CreatedAtName, &syndb.ColData{IdVal: roomId, ColVal: now})
	syndb.AddData(TbLiveRoomCfg, db.UpdatedAtName, &syndb.ColData{IdVal: roomId, ColVal: now})
	r.SetCategory(LiveRoomCategoryHot)
	r.SetPrivateInviteType(DefaultPrivateInviteType(LiveRoomCategoryHot))
	r.SetCloudPlayerFrameRate(DefaultCloudPlayerFrameRate)
	r.SetCloudPlayerBitrate(DefaultCloudPlayerBitrate)
	r.SetCloudPlayerWidth(DefaultCloudPlayerWidth)
	r.SetCloudPlayerHeight(DefaultCloudPlayerHeight)
	return r
}

func NormalizeCloudPlayerFrameRate(v uint8) uint8 {
	if v == 0 {
		return DefaultCloudPlayerFrameRate
	}
	if v < 1 {
		return 1
	}
	if v > 30 {
		return 30
	}
	return v
}

func NormalizeCloudPlayerBitrate(v int) int {
	if v <= 0 {
		return DefaultCloudPlayerBitrate
	}
	if v > 10000 {
		return 10000
	}
	return v
}

func NormalizeCloudPlayerWidth(v int) int {
	if v <= 0 {
		return DefaultCloudPlayerWidth
	}
	if v > 4096 {
		return 4096
	}
	return v
}

func NormalizeCloudPlayerHeight(v int) int {
	if v <= 0 {
		return DefaultCloudPlayerHeight
	}
	if v > 2160 {
		return 2160
	}
	return v
}

func (r *LiveRoomCfg) EffectiveCloudPlayerFrameRate() int {
	if r == nil {
		return DefaultCloudPlayerFrameRate
	}
	return int(NormalizeCloudPlayerFrameRate(r.CloudPlayerFrameRate))
}

func (r *LiveRoomCfg) EffectiveCloudPlayerBitrate() int {
	if r == nil {
		return DefaultCloudPlayerBitrate
	}
	return NormalizeCloudPlayerBitrate(r.CloudPlayerBitrate)
}

func (r *LiveRoomCfg) EffectiveCloudPlayerWidth() int {
	if r == nil {
		return DefaultCloudPlayerWidth
	}
	return NormalizeCloudPlayerWidth(r.CloudPlayerWidth)
}

func (r *LiveRoomCfg) EffectiveCloudPlayerHeight() int {
	if r == nil {
		return DefaultCloudPlayerHeight
	}
	return NormalizeCloudPlayerHeight(r.CloudPlayerHeight)
}

func (r *LiveRoomCfg) touchUpdatedAt() {
	r.UpdatedAt = time.Now()
	syndb.AddData(TbLiveRoomCfg, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: r.UpdatedAt})
}

func (r *LiveRoomCfg) SetPrivateInviteType(v uint8) {
	v = NormalizePrivateInviteType(v, r.Category)
	r.PrivateInviteType = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgPrivateInviteType, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetCategory(v uint8) {
	if v != LiveRoomCategoryHot && v != LiveRoomCategoryGame && v != LiveRoomCategoryPrivate {
		v = LiveRoomCategoryHot
	}
	r.Category = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgCategory, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetTagId(v uint64) {
	r.TagId = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgTagId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetTicket(v float64) {
	r.Ticket = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgTicket, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetBilling(v float64) {
	r.Billing = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgBilling, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetCloudPlayerVideo(v string) {
	r.CloudPlayerVideo = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgCloudPlayerVideo, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetCloudPlayerFrameRate(v uint8) {
	v = NormalizeCloudPlayerFrameRate(v)
	r.CloudPlayerFrameRate = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgCloudPlayerFrameRate, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetCloudPlayerBitrate(v int) {
	v = NormalizeCloudPlayerBitrate(v)
	r.CloudPlayerBitrate = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgCloudPlayerBitrate, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetCloudPlayerWidth(v int) {
	v = NormalizeCloudPlayerWidth(v)
	r.CloudPlayerWidth = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgCloudPlayerWidth, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetCloudPlayerHeight(v int) {
	v = NormalizeCloudPlayerHeight(v)
	r.CloudPlayerHeight = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgCloudPlayerHeight, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetPushStream(v bool) {
	r.PushStream = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgPushStream, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetIsTest(v bool) {
	r.IsTest = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgIsTest, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetCloudPlayerId(v string) {
	r.CloudPlayerId = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgCloudPlayerId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *LiveRoomCfg) SetCloudPlayerTokenExpireAt(v *time.Time) {
	r.CloudPlayerTokenExpireAt = v
	r.touchUpdatedAt()
	syndb.AddData(TbLiveRoomCfg, LiveRoomCfgCloudPlayerTokenExpireAt, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initLiveRoomCfg() {
	syndb.RegQuick(TbLiveRoomCfg, db.CreatedAtName)
	syndb.RegQuick(TbLiveRoomCfg, db.UpdatedAtName)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgPrivateInviteType)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgCategory)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgTagId)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgTicket)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgBilling)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgCloudPlayerVideo)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgCloudPlayerFrameRate)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgCloudPlayerBitrate)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgCloudPlayerWidth)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgCloudPlayerHeight)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgPushStream)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgIsTest)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgCloudPlayerId)
	syndb.RegQuick(TbLiveRoomCfg, LiveRoomCfgCloudPlayerTokenExpireAt)
	migrate.AutoMigrate(&LiveRoomCfg{})
}

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
	TbUserInfo db.TbName = "user_infos"
)

const (
	UserInfoNickname      db.TbCol = "nickname"
	UserInfoPhone         db.TbCol = "phone"
	UserInfoAvatar        db.TbCol = "avatar"
	UserInfoRemark        db.TbCol = "remark"
	UserInfoGold          db.TbCol = "gold"
	UserInfoDiamond       db.TbCol = "diamond"
	UserInfoShareCode     db.TbCol = "share_code"
	UserInfoGuildId       db.TbCol = "guild_id"
	UserInfoIsAnchor      db.TbCol = "is_anchor"
	UserInfoHasLiveRoom   db.TbCol = "has_live_room"
	UserInfoInviterId     db.TbCol = "inviter_id"
	UserInfoVipLevel      db.TbCol = "vip_level"
	UserInfoLastLoginTime db.TbCol = "last_login_time"
	UserInfoLiveRoomId    db.TbCol = "live_room_id"
	UserInfoGender        db.TbCol = "gender"
	UserInfoBirthday      db.TbCol = "birthday"
)

// UserInfo 用户基础信息
type UserInfo struct {
	migrate.OneModel
	Nickname      string     `gorm:"default:'';comment:用户昵称"`
	Phone         string     `gorm:"default:'';comment:手机号"`
	Avatar        string     `gorm:"default:'';comment:头像"`
	Remark        string     `gorm:"default:'';comment:备注"`
	Gold          float64    `gorm:"default:0;comment:金币"`
	Diamond       float64    `gorm:"default:0;comment:钻石"`
	ShareCode     string     `gorm:"uniqueIndex;default:'';comment:分享码"`
	GuildId       uint64     `gorm:"index;default:0;comment:所属工会ID(0为未加入)"`
	IsAnchor      bool       `gorm:"default:0;comment:是否主播(设为true后不可回退)"`
	HasLiveRoom   bool       `gorm:"default:0;comment:是否已创建直播间(App端完善资料后为true)"`
	InviterId     uint64     `gorm:"index;default:0;comment:邀请人用户ID(0为无)"`
	VipLevel      uint32     `gorm:"default:0;comment:VIP等级(0为无)"`
	LastLoginTime *time.Time `gorm:"comment:最后登录时间" json:"lastLoginTime"`
	LiveRoomId    uint64     `gorm:"index;default:0;comment:当前所在直播间ID(观众,0为不在直播间)" json:"liveRoomId"`
	Gender        uint8      `gorm:"default:0;comment:性别(0未知,1男,2女)" json:"gender"`
	Birthday      *time.Time `gorm:"type:date;comment:出生日期" json:"birthday"`
}

func NewUserInfo(userId uint64) *UserInfo {
	ret := &UserInfo{}
	ret.ID = userId
	ret.SetCreatedAt(time.Now())
	ret.SetShareCode(fmt.Sprintf("%d", userId))
	ret.SetUpdatedAt(time.Now())
	return ret
}

func (receiver *UserInfo) SetNickname(nickname string) {
	receiver.Nickname = nickname
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoNickname, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: nickname,
	})
}

func (receiver *UserInfo) SetPhone(phone string) {
	receiver.Phone = phone
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoPhone, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: phone,
	})
}

func (receiver *UserInfo) SetAvatar(avatar string) {
	receiver.Avatar = avatar
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoAvatar, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: avatar,
	})
}

func (receiver *UserInfo) SetRemark(remark string) {
	receiver.Remark = remark
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoRemark, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: remark,
	})
}

func (receiver *UserInfo) AddGold(gold float64) {
	receiver.Gold = math.AddFloat64(gold, receiver.Gold)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserInfo, UserInfoGold, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.Gold,
	})
}

func (receiver *UserInfo) SubGold(gold float64) {
	receiver.Gold = math.SubFloat64(receiver.Gold, gold)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserInfo, UserInfoGold, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.Gold,
	})
}

func (receiver *UserInfo) AddDiamond(diamond float64) {
	receiver.Diamond = math.AddFloat64(diamond, receiver.Diamond)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserInfo, UserInfoDiamond, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.Diamond,
	})
}

func (receiver *UserInfo) SubDiamond(diamond float64) {
	receiver.Diamond = math.SubFloat64(receiver.Diamond, diamond)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbUserInfo, UserInfoDiamond, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.Diamond,
	})
}

func (receiver *UserInfo) SetShareCode(shareCode string) {
	receiver.ShareCode = shareCode
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoShareCode, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: shareCode,
	})
}

func (receiver *UserInfo) SetGuildId(guildId uint64) {
	receiver.GuildId = guildId
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoGuildId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: guildId,
	})
}

func (receiver *UserInfo) SetIsAnchor(isAnchor bool) {
	receiver.IsAnchor = isAnchor
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoIsAnchor, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: isAnchor,
	})
}

func (receiver *UserInfo) SetHasLiveRoom(hasLiveRoom bool) {
	receiver.HasLiveRoom = hasLiveRoom
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoHasLiveRoom, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: hasLiveRoom,
	})
}

func (receiver *UserInfo) SetInviterId(inviterId uint64) {
	receiver.InviterId = inviterId
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoInviterId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: inviterId,
	})
}

func (receiver *UserInfo) SetVipLevel(vipLevel uint32) {
	receiver.VipLevel = vipLevel
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoVipLevel, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: vipLevel,
	})
}

func (receiver *UserInfo) SetLastLoginTime(val *time.Time) {
	receiver.LastLoginTime = val
	syndb.AddDataToLazyChan(TbUserInfo, UserInfoLastLoginTime, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserInfo) SetLiveRoomId(liveRoomId uint64) {
	receiver.LiveRoomId = liveRoomId
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoLiveRoomId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: liveRoomId,
	})
}

func (receiver *UserInfo) SetGender(gender uint8) {
	receiver.Gender = gender
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoGender, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: gender,
	})
}

func (receiver *UserInfo) SetBirthday(val *time.Time) {
	receiver.Birthday = val
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserInfo, UserInfoBirthday, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserInfo) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToQuickChan(TbUserInfo, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserInfo) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToLazyChan(TbUserInfo, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initUserInfo() {
	syndb.RegQuickWithLarge(TbUserInfo, db.CreatedAtName)
	syndb.RegLazyWithLarge(TbUserInfo, db.UpdatedAtName)

	syndb.RegQuickWithLarge(TbUserInfo, UserInfoNickname)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoPhone)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoAvatar)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoRemark)
	syndb.RegLazyWithLarge(TbUserInfo, UserInfoGold)
	syndb.RegLazyWithLarge(TbUserInfo, UserInfoDiamond)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoShareCode)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoGuildId)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoIsAnchor)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoHasLiveRoom)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoInviterId)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoVipLevel)
	syndb.RegLazyWithLarge(TbUserInfo, UserInfoLastLoginTime)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoLiveRoomId)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoGender)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoBirthday)

	migrate.AutoMigrate(&UserInfo{})
}

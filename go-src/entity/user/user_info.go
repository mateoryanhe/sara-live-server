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
	UserInfoNickname        db.TbCol = "nickname"
	UserInfoPhone           db.TbCol = "phone"
	UserInfoAvatar          db.TbCol = "avatar"
	UserInfoRemark          db.TbCol = "remark"
	UserInfoGold            db.TbCol = "gold"
	UserInfoDiamond         db.TbCol = "diamond"
	UserInfoShareCode       db.TbCol = "share_code"
	UserInfoUserType        db.TbCol = "user_type"
	UserInfoInviterId       db.TbCol = "inviter_id"
	UserInfoVipLevel        db.TbCol = "vip_level"
	UserInfoLastLoginTime   db.TbCol = "last_login_time"
	UserInfoLiveRoomId      db.TbCol = "live_room_id"
	UserInfoLiveRoomVer     db.TbCol = "live_room_ver"
	UserInfoGender          db.TbCol = "gender"
	UserInfoBirthday        db.TbCol = "birthday"
	UserInfoBotAnchorStatus db.TbCol = "bot_anchor_status"
)

const (
	UserTypeNormal       uint8 = 0 // 普通用户
	UserTypeAnchor       uint8 = 1 // 普通主播
	UserTypeBotAnchor    uint8 = 2 // 机器人主播
	UserTypeBotAudience  uint8 = 3 // 机器人观众(不参与系统统计)
	UserTypeTester       uint8 = 4 // 测试人员(不参与系统统计)
	UserTypeCMSAuthor    uint8 = 5 // CMS短视频作者(不参与系统统计)
	UserTypeSeniorAnchor uint8 = 7 // 高级主播
)

const (
	BotAnchorStatusDisabled uint8 = 0 // 停用
	BotAnchorStatusEnabled  uint8 = 1 // 启用
)

// UserInfo 用户基础信息
type UserInfo struct {
	migrate.OneModel
	Nickname        string     `gorm:"default:'';comment:用户昵称"`
	Phone           string     `gorm:"default:'';comment:手机号"`
	Avatar          string     `gorm:"default:'';comment:头像"`
	Remark          string     `gorm:"default:'';comment:备注"`
	Gold            float64    `gorm:"default:0;comment:金币"`
	Diamond         float64    `gorm:"default:0;comment:钻石"`
	ShareCode       string     `gorm:"uniqueIndex;default:'';comment:分享码"`
	UserType        uint8      `gorm:"default:0;comment:用户类型(0普通用户,1普通主播,2机器人主播,3机器人观众,4测试人员,5CMS短视频作者,7高级主播)" json:"userType"`
	InviterId       uint64     `gorm:"index;default:0;comment:邀请人用户ID(0为无)"`
	VipLevel        uint32     `gorm:"default:0;comment:VIP等级(0为无)"`
	LastLoginTime   *time.Time `gorm:"index;comment:最后登录时间" json:"lastLoginTime"`
	LiveRoomId      uint64     `gorm:"index;default:0;comment:当前所在直播间ID(观众,0为不在直播间)" json:"liveRoomId"`
	LiveRoomVer     uint64     `gorm:"default:0;comment:当前所在直播间版本(观众,0为无,通常为liveRecordId)" json:"liveRoomVer"`
	Gender          uint8      `gorm:"default:0;comment:性别(0未知,1男,2女)" json:"gender"`
	Birthday        *time.Time `gorm:"type:date;comment:出生日期" json:"birthday"`
	BotAnchorStatus uint8      `gorm:"default:1;comment:机器人主播状态(0停用,1启用)" json:"botAnchorStatus"`
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
	syndb.AddData(TbUserInfo, UserInfoNickname, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: nickname,
	})
}

func (receiver *UserInfo) SetPhone(phone string) {
	receiver.Phone = phone
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoPhone, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: phone,
	})
}

func (receiver *UserInfo) SetAvatar(avatar string) {
	receiver.Avatar = avatar
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoAvatar, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: avatar,
	})
}

func (receiver *UserInfo) SetRemark(remark string) {
	receiver.Remark = remark
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoRemark, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: remark,
	})
}

func (receiver *UserInfo) AddGold(gold float64) {
	receiver.Gold = math.AddFloat64(gold, receiver.Gold)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoGold, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.Gold,
	})
}

func (receiver *UserInfo) SubGold(gold float64) {
	receiver.Gold = math.SubFloat64(receiver.Gold, gold)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoGold, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.Gold,
	})
}

func (receiver *UserInfo) AddDiamond(diamond float64) {
	receiver.Diamond = math.AddFloat64(diamond, receiver.Diamond)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoDiamond, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.Diamond,
	})
}

func (receiver *UserInfo) SubDiamond(diamond float64) {
	receiver.Diamond = math.SubFloat64(receiver.Diamond, diamond)
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoDiamond, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: receiver.Diamond,
	})
}

func (receiver *UserInfo) SetShareCode(shareCode string) {
	receiver.ShareCode = shareCode
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoShareCode, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: shareCode,
	})
}

func UserTypeIsAnchor(userType uint8) bool {
	return userType == UserTypeAnchor || userType == UserTypeBotAnchor || userType == UserTypeSeniorAnchor
}

func UserTypeExcludedFromStat(userType uint8) bool {
	return userType == UserTypeBotAudience || userType == UserTypeTester || userType == UserTypeCMSAuthor
}

func (receiver *UserInfo) IsAnchor() bool {
	return UserTypeIsAnchor(receiver.UserType)
}

func (receiver *UserInfo) IsBotAnchor() bool {
	return receiver.UserType == UserTypeBotAnchor
}

func (receiver *UserInfo) SetUserType(userType uint8) {
	switch userType {
	case UserTypeNormal, UserTypeAnchor, UserTypeBotAnchor, UserTypeBotAudience, UserTypeTester, UserTypeCMSAuthor, UserTypeSeniorAnchor:
	default:
		userType = UserTypeNormal
	}
	receiver.UserType = userType
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoUserType, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: userType,
	})
}

func (receiver *UserInfo) SetInviterId(inviterId uint64) {
	receiver.InviterId = inviterId
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoInviterId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: inviterId,
	})
}

func (receiver *UserInfo) SetVipLevel(vipLevel uint32) {
	receiver.VipLevel = vipLevel
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoVipLevel, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: vipLevel,
	})
}

func (receiver *UserInfo) SetLastLoginTime(val *time.Time) {
	receiver.LastLoginTime = val
	syndb.AddData(TbUserInfo, UserInfoLastLoginTime, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserInfo) SetLiveRoomId(liveRoomId uint64) {
	receiver.LiveRoomId = liveRoomId
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoLiveRoomId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: liveRoomId,
	})
}

func (receiver *UserInfo) SetLiveRoomVer(liveRoomVer uint64) {
	receiver.LiveRoomVer = liveRoomVer
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoLiveRoomVer, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: liveRoomVer,
	})
}

func (receiver *UserInfo) SetGender(gender uint8) {
	receiver.Gender = gender
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoGender, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: gender,
	})
}

func (receiver *UserInfo) SetBirthday(val *time.Time) {
	receiver.Birthday = val
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoBirthday, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserInfo) SetBotAnchorStatus(status uint8) {
	receiver.BotAnchorStatus = status
	receiver.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserInfo, UserInfoBotAnchorStatus, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: status,
	})
}

func (receiver *UserInfo) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddData(TbUserInfo, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserInfo) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddData(TbUserInfo, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initUserInfo() {
	syndb.RegQuick(TbUserInfo, db.CreatedAtName)
	syndb.RegLazy(TbUserInfo, db.UpdatedAtName)

	syndb.RegQuick(TbUserInfo, UserInfoNickname)
	syndb.RegQuick(TbUserInfo, UserInfoPhone)
	syndb.RegQuick(TbUserInfo, UserInfoAvatar)
	syndb.RegQuick(TbUserInfo, UserInfoRemark)
	syndb.RegLazy(TbUserInfo, UserInfoGold)
	syndb.RegLazy(TbUserInfo, UserInfoDiamond)
	syndb.RegQuick(TbUserInfo, UserInfoShareCode)
	syndb.RegQuick(TbUserInfo, UserInfoUserType)
	syndb.RegQuick(TbUserInfo, UserInfoInviterId)
	syndb.RegQuick(TbUserInfo, UserInfoVipLevel)
	syndb.RegLazy(TbUserInfo, UserInfoLastLoginTime)
	syndb.RegQuick(TbUserInfo, UserInfoLiveRoomId)
	syndb.RegQuick(TbUserInfo, UserInfoLiveRoomVer)
	syndb.RegQuick(TbUserInfo, UserInfoGender)
	syndb.RegQuick(TbUserInfo, UserInfoBirthday)
	syndb.RegQuick(TbUserInfo, UserInfoBotAnchorStatus)

	migrate.AutoMigrate(&UserInfo{})

}

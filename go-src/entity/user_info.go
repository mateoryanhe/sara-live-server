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
	UserInfoNickname  db.TbCol = "nickname"
	UserInfoPhone     db.TbCol = "phone"
	UserInfoAvatar    db.TbCol = "avatar"
	UserInfoRemark    db.TbCol = "remark"
	UserInfoGold      db.TbCol = "gold"
	UserInfoDiamond   db.TbCol = "diamond"
	UserInfoShareCode db.TbCol = "share_code"
	UserInfoGuildId   db.TbCol = "guild_id"
	UserInfoIsAnchor  db.TbCol = "is_anchor"
	UserInfoInviterId db.TbCol = "inviter_id"
	UserInfoVipLevel  db.TbCol = "vip_level"
)

// UserInfo 用户基础信息
type UserInfo struct {
	migrate.OneModel
	Nickname  string  `gorm:"default:'';comment:用户昵称"`
	Phone     string  `gorm:"default:'';comment:手机号"`
	Avatar    string  `gorm:"default:'';comment:头像"`
	Remark    string  `gorm:"default:'';comment:备注"`
	Gold      float64 `gorm:"default:0;comment:金币"`
	Diamond   float64 `gorm:"default:0;comment:钻石"`
	ShareCode string  `gorm:"uniqueIndex;default:'';comment:分享码"`
	GuildId   uint64  `gorm:"index;default:0;comment:所属工会ID(0为未加入)"`
	IsAnchor  bool    `gorm:"default:0;comment:是否主播(设为true后不可回退)"`
	InviterId uint64  `gorm:"index;default:0;comment:邀请人用户ID(0为无)"`
	VipLevel  uint32  `gorm:"default:0;comment:VIP等级(0为无)"`
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
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoInviterId)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoVipLevel)

	migrate.AutoMigrate(&UserInfo{})
}

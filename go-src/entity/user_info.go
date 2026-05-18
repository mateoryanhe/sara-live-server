package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbUserInfo db.TbName = "user_infos"
)

const (
	UserInfoNickname db.TbCol = "nickname"
	UserInfoPhone    db.TbCol = "phone"
	UserInfoAvatar   db.TbCol = "avatar"
	UserInfoRemark   db.TbCol = "remark"
)

// UserInfo 用户基础信息
type UserInfo struct {
	migrate.OneModel
	Nickname string `gorm:"default:'';comment:用户昵称"`
	Phone    string `gorm:"default:'';comment:手机号"`
	Avatar   string `gorm:"default:'';comment:头像"`
	Remark   string `gorm:"default:'';comment:备注"`
}

func NewUserInfo(userId uint64) *UserInfo {
	ret := &UserInfo{}
	ret.ID = userId
	ret.SetCreatedAt(time.Now())
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

func (receiver *UserInfo) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToQuickChan(TbUserInfo, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserInfo) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToQuickChan(TbUserInfo, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initUserInfo() {
	syndb.RegQuickWithLarge(TbUserInfo, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbUserInfo, db.UpdatedAtName)

	syndb.RegQuickWithLarge(TbUserInfo, UserInfoNickname)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoPhone)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoAvatar)
	syndb.RegQuickWithLarge(TbUserInfo, UserInfoRemark)

	migrate.AutoMigrate(&UserInfo{})
}

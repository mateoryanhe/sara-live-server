package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbChannelPayUserProfile db.TbName = "channel_pay_user_profiles"
)

const (
	ChannelPayUserProfileName  db.TbCol = "name"
	ChannelPayUserProfileEmail db.TbCol = "email"
)

// ChannelPayUserProfile 渠道支付付款人资料(主键=userId，缓存模式同 user_infos)
type ChannelPayUserProfile struct {
	migrate.OneModel
	Name  string `gorm:"size:128;default:'';comment:付款人姓名(HaiPay name)" json:"name"`
	Email string `gorm:"size:256;default:'';comment:付款人邮箱" json:"email"`
}

func (ChannelPayUserProfile) TableName() string {
	return string(TbChannelPayUserProfile)
}

func NewChannelPayUserProfile(userId uint64) *ChannelPayUserProfile {
	ret := &ChannelPayUserProfile{}
	ret.ID = userId
	ret.SetCreatedAt(time.Now())
	ret.SetUpdatedAt(time.Now())
	return ret
}

func (p *ChannelPayUserProfile) SetName(name string) {
	p.Name = name
	p.SetUpdatedAt(time.Now())
	syndb.AddData(TbChannelPayUserProfile, ChannelPayUserProfileName, &syndb.ColData{
		IdVal:  p.ID,
		ColVal: name,
	})
}

func (p *ChannelPayUserProfile) SetEmail(email string) {
	p.Email = email
	p.SetUpdatedAt(time.Now())
	syndb.AddData(TbChannelPayUserProfile, ChannelPayUserProfileEmail, &syndb.ColData{
		IdVal:  p.ID,
		ColVal: email,
	})
}

func (p *ChannelPayUserProfile) SetCreatedAt(val time.Time) {
	p.CreatedAt = val
	syndb.AddData(TbChannelPayUserProfile, db.CreatedAtName, &syndb.ColData{
		IdVal:  p.ID,
		ColVal: val,
	})
}

func (p *ChannelPayUserProfile) SetUpdatedAt(val time.Time) {
	p.UpdatedAt = val
	syndb.AddData(TbChannelPayUserProfile, db.UpdatedAtName, &syndb.ColData{
		IdVal:  p.ID,
		ColVal: val,
	})
}

func initChannelPayUserProfile() {
	migrate.AutoMigrate(&ChannelPayUserProfile{})
	syndb.RegQuick(TbChannelPayUserProfile, db.CreatedAtName)
	syndb.RegLazy(TbChannelPayUserProfile, db.UpdatedAtName)
	syndb.RegQuick(TbChannelPayUserProfile, ChannelPayUserProfileName)
	syndb.RegQuick(TbChannelPayUserProfile, ChannelPayUserProfileEmail)
}

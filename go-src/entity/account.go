package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbAccount db.TbName = "accounts"
)

const (
	AccountOpenId       db.TbCol = "open_id"
	AccountChannel      db.TbCol = "channel"
	AccountIP           db.TbCol = "ip"
	AccountBan          db.TbCol = "ban"
	AccountBanTime      db.TbCol = "ban_time"
	AccountBanApplyTime db.TbCol = "ban_apply_time"
	AccountCancel       db.TbCol = "cancel"
)

type Account struct {
	migrate.OneModel
	OpenId       string     `gorm:"default:'';comment:开放id"`
	IP           string     `gorm:"default:'';comment:ip地址"`
	Channel      uint       `gorm:"default:0;comment:渠道id"`
	Ban          bool       `gorm:"default:0;comment:封号"`
	BanTime      *time.Time `gorm:"comment:封号时间"`
	BanApplyTime *time.Time `gorm:"comment:封号生效时间"`
	Cancel       bool       `gorm:"default:0;comment:注销"`
}

func NewAccount(openId string, channel uint) *Account {
	ret := &Account{}
	ret.ID = snowflake.GetId()
	ret.SetChannel(channel)
	ret.SetOpenId(openId)
	ret.SetCancel(false)
	ret.SetCreatedAt(time.Now())
	ret.SetUpdatedAt(time.Now())
	return ret
}

func (this *Account) SetOpenId(openId string) {
	this.OpenId = openId
	syndb.AddDataToQuickChan(TbAccount, AccountOpenId, &syndb.ColData{
		IdVal:  this.ID,
		ColVal: openId,
	})
}

func (receiver *Account) SetIp(ip string) {
	receiver.IP = ip
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountIP, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: ip,
	})
}

func (receiver *Account) SetChannel(channel uint) {
	receiver.Channel = channel
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountChannel, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: channel,
	})
}

func (receiver *Account) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToQuickChan(TbAccount, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *Account) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToQuickChan(TbAccount, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (this *Account) SetBan(ban bool) {
	this.Ban = ban
	this.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountBan, &syndb.ColData{
		IdVal:  this.ID,
		ColVal: ban,
	})
}

func (receiver *Account) SetBanTime(banTime *time.Time) {
	receiver.BanTime = banTime
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountBanTime, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: banTime,
	})
}

func (receiver *Account) SetBanApplyTime(banApplyTime *time.Time) {
	receiver.BanApplyTime = banApplyTime
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountBanApplyTime, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: banApplyTime,
	})
}

func (this *Account) SetCancel(cancel bool) {
	this.Cancel = cancel
	this.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountCancel, &syndb.ColData{
		IdVal:  this.ID,
		ColVal: cancel,
	})
}
func initAccount() {
	syndb.RegQuickWithLarge(TbAccount, db.CreatedAtName)
	syndb.RegQuickWithLarge(TbAccount, db.UpdatedAtName)
	syndb.RegQuickWithLarge(TbAccount, db.DeletedAtName)
	syndb.RegQuickWithLarge(TbAccount, db.IsDeletedName)

	syndb.RegQuickWithLarge(TbAccount, AccountOpenId)
	syndb.RegQuickWithLarge(TbAccount, AccountChannel)
	syndb.RegQuickWithLarge(TbAccount, AccountIP)
	syndb.RegQuickWithLarge(TbAccount, AccountBan)
	syndb.RegQuickWithLarge(TbAccount, AccountBanTime)
	syndb.RegQuickWithLarge(TbAccount, AccountBanApplyTime)
	syndb.RegQuickWithLarge(TbAccount, AccountCancel)

	migrate.AutoMigrate(&Account{})
}

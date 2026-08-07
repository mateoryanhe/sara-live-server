package entity

import (
	"strings"
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
	AccountOpenId          db.TbCol = "open_id"
	AccountChannel         db.TbCol = "channel"
	AccountPhoneAreaCode   db.TbCol = "phone_area_code"
	AccountIP              db.TbCol = "ip"
	AccountRegisterIp      db.TbCol = "register_ip"
	AccountRegisterCountry db.TbCol = "register_country"
	AccountLoginCountry    db.TbCol = "login_country"
	AccountBan             db.TbCol = "ban"
	AccountBanTime         db.TbCol = "ban_time"
	AccountBanApplyTime    db.TbCol = "ban_apply_time"
	AccountCancel          db.TbCol = "cancel"
	AccountPassword        db.TbCol = "password"
)

type Account struct {
	migrate.OneModel
	OpenId          string     `gorm:"default:'';comment:开放id"`
	PhoneAreaCode   string     `gorm:"default:'';comment:手机区号"`
	IP              string     `gorm:"default:'';comment:登录IP"`
	RegisterIp      string     `gorm:"default:'';comment:注册IP"`
	RegisterCountry string     `gorm:"default:'';comment:注册IP所在国家"`
	LoginCountry    string     `gorm:"default:'';comment:登录IP所在国家"`
	Channel         uint       `gorm:"default:0;comment:渠道id"`
	Ban             bool       `gorm:"default:0;comment:封号"`
	BanTime         *time.Time `gorm:"comment:封号时间"`
	BanApplyTime    *time.Time `gorm:"comment:封号生效时间"`
	Cancel          bool       `gorm:"default:0;comment:注销"`
	Password        string     `gorm:"default:'';comment:密码"`
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

// NormalizeOpenId 规范化 open_id
func NormalizeOpenId(openId string) string {
	return strings.TrimSpace(openId)
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

func (receiver *Account) SetRegisterIp(ip string) {
	receiver.RegisterIp = ip
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountRegisterIp, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: ip,
	})
}

func (receiver *Account) SetRegisterCountry(country string) {
	receiver.RegisterCountry = country
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountRegisterCountry, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: country,
	})
}

func (receiver *Account) SetLoginCountry(country string) {
	receiver.LoginCountry = country
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountLoginCountry, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: country,
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

func (receiver *Account) SetPhoneAreaCode(phoneAreaCode string) {
	receiver.PhoneAreaCode = phoneAreaCode
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountPhoneAreaCode, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: phoneAreaCode,
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

func (receiver *Account) SetPassword(password string) {
	receiver.Password = password
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbAccount, AccountPassword, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: password,
	})
}
func initAccount() {
	syndb.RegQuick(TbAccount, db.CreatedAtName)
	syndb.RegQuick(TbAccount, db.UpdatedAtName)
	syndb.RegQuick(TbAccount, db.DeletedAtName)
	syndb.RegQuick(TbAccount, db.IsDeletedName)

	syndb.RegQuick(TbAccount, AccountOpenId)
	syndb.RegQuick(TbAccount, AccountPhoneAreaCode)
	syndb.RegQuick(TbAccount, AccountChannel)
	syndb.RegQuick(TbAccount, AccountIP)
	syndb.RegQuick(TbAccount, AccountRegisterIp)
	syndb.RegQuick(TbAccount, AccountRegisterCountry)
	syndb.RegQuick(TbAccount, AccountLoginCountry)
	syndb.RegQuick(TbAccount, AccountBan)
	syndb.RegQuick(TbAccount, AccountBanTime)
	syndb.RegQuick(TbAccount, AccountBanApplyTime)
	syndb.RegQuick(TbAccount, AccountCancel)
	syndb.RegQuick(TbAccount, AccountPassword)

	migrate.AutoMigrate(&Account{})
}

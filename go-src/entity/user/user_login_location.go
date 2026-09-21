package entity

import (
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbUserLoginLocation db.TbName = "user_login_locations"

const (
	UserLoginLocationUserId          db.TbCol = "user_id"
	UserLoginLocationIP              db.TbCol = "ip"
	UserLoginLocationRegisterIP      db.TbCol = "register_ip"
	UserLoginLocationRegisterCountry db.TbCol = "register_country"
	UserLoginLocationLoginCountry    db.TbCol = "login_country"
)

// UserLoginLocation 用户登录 IP 与国家信息，user_id 是唯一主键。
type UserLoginLocation struct {
	UserId          uint64    `gorm:"column:user_id;primaryKey;autoIncrement:false;comment:用户ID" json:"userId,string"`
	IP              string    `gorm:"column:ip;size:64;default:'';comment:登录IP" json:"ip"`
	RegisterIp      string    `gorm:"column:register_ip;size:64;default:'';comment:注册IP" json:"registerIp"`
	RegisterCountry string    `gorm:"column:register_country;size:8;default:'';comment:注册IP所在国家简码(ISO alpha-2)" json:"registerCountry"`
	LoginCountry    string    `gorm:"column:login_country;size:8;default:'';comment:登录IP所在国家简码(ISO alpha-2)" json:"loginCountry"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

func (UserLoginLocation) TableName() string {
	return string(TbUserLoginLocation)
}

func NewUserLoginLocation(userId uint64) *UserLoginLocation {
	row := &UserLoginLocation{UserId: userId}
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	return row
}

func (row *UserLoginLocation) SetIP(ip string) {
	row.IP = ip
	row.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserLoginLocation, UserLoginLocationIP, &syndb.ColData{
		IdVal:  row.UserId,
		ColVal: ip,
	})
}

func (row *UserLoginLocation) SetRegisterIP(ip string) {
	row.RegisterIp = ip
	row.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserLoginLocation, UserLoginLocationRegisterIP, &syndb.ColData{
		IdVal:  row.UserId,
		ColVal: ip,
	})
}

func (row *UserLoginLocation) SetRegisterCountry(country string) {
	row.RegisterCountry = country
	row.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserLoginLocation, UserLoginLocationRegisterCountry, &syndb.ColData{
		IdVal:  row.UserId,
		ColVal: country,
	})
}

func (row *UserLoginLocation) SetLoginCountry(country string) {
	row.LoginCountry = country
	row.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserLoginLocation, UserLoginLocationLoginCountry, &syndb.ColData{
		IdVal:  row.UserId,
		ColVal: country,
	})
}

func (row *UserLoginLocation) SetCreatedAt(val time.Time) {
	row.CreatedAt = val
	syndb.AddData(TbUserLoginLocation, db.CreatedAtName, &syndb.ColData{
		IdVal:  row.UserId,
		ColVal: val,
	})
}

func (row *UserLoginLocation) SetUpdatedAt(val time.Time) {
	row.UpdatedAt = val
	syndb.AddData(TbUserLoginLocation, db.UpdatedAtName, &syndb.ColData{
		IdVal:  row.UserId,
		ColVal: val,
	})
}

func initUserLoginLocation() {
	columns := []db.TbCol{
		db.CreatedAtName,
		db.UpdatedAtName,
		UserLoginLocationIP,
		UserLoginLocationRegisterIP,
		UserLoginLocationRegisterCountry,
		UserLoginLocationLoginCountry,
	}
	for _, column := range columns {
		syndb.RegWithIDName(TbUserLoginLocation, column, UserLoginLocationUserId)
	}
	migrate.AutoMigrate(&UserLoginLocation{})
}

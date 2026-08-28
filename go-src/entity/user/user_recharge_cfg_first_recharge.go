package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

const (
	TbUserRechargeCfgFirstRecharge db.TbName = "user_recharge_cfg_first_recharges"
)

const (
	UserRechargeCfgFirstRechargeUserId        db.TbCol = "user_id"
	UserRechargeCfgFirstRechargeCfgId         db.TbCol = "cfg_id"
	UserRechargeCfgFirstRechargeFirstRecharge db.TbCol = "first_recharge"
)

// UserRechargeCfgFirstRecharge 用户充值档位首充状态(每用户每档位一行,唯一性由代码保证)
type UserRechargeCfgFirstRecharge struct {
	migrate.OneModel
	UserId        uint64 `gorm:"index;default:0;comment:用户ID" json:"userId"`
	CfgId         uint64 `gorm:"default:0;comment:充值档位ID" json:"cfgId"`
	FirstRecharge bool   `gorm:"default:1;comment:是否档位首充(1=未首充,0=已首充)" json:"firstRecharge"`
}

func NewUserRechargeCfgFirstRecharge(userId, cfgId uint64) *UserRechargeCfgFirstRecharge {
	row := &UserRechargeCfgFirstRecharge{}
	row.ID = snowflake.GetId()
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetUserId(userId)
	row.SetCfgId(cfgId)
	row.SetFirstRecharge(true)
	return row
}

func (r *UserRechargeCfgFirstRecharge) SetUserId(v uint64) {
	r.UserId = v
	syndb.AddData(TbUserRechargeCfgFirstRecharge, UserRechargeCfgFirstRechargeUserId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *UserRechargeCfgFirstRecharge) SetCfgId(v uint64) {
	r.CfgId = v
	syndb.AddData(TbUserRechargeCfgFirstRecharge, UserRechargeCfgFirstRechargeCfgId, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *UserRechargeCfgFirstRecharge) SetFirstRecharge(v bool) {
	r.FirstRecharge = v
	r.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserRechargeCfgFirstRecharge, UserRechargeCfgFirstRechargeFirstRecharge, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *UserRechargeCfgFirstRecharge) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbUserRechargeCfgFirstRecharge, db.CreatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func (r *UserRechargeCfgFirstRecharge) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbUserRechargeCfgFirstRecharge, db.UpdatedAtName, &syndb.ColData{
		IdVal: r.ID, ColVal: v,
	})
}

func initUserRechargeCfgFirstRecharge() {
	syndb.RegQuick(TbUserRechargeCfgFirstRecharge, db.CreatedAtName)
	syndb.RegQuick(TbUserRechargeCfgFirstRecharge, db.UpdatedAtName)
	syndb.RegQuick(TbUserRechargeCfgFirstRecharge, UserRechargeCfgFirstRechargeUserId)
	syndb.RegQuick(TbUserRechargeCfgFirstRecharge, UserRechargeCfgFirstRechargeCfgId)
	syndb.RegQuick(TbUserRechargeCfgFirstRecharge, UserRechargeCfgFirstRechargeFirstRecharge)
	migrate.AutoMigrate(&UserRechargeCfgFirstRecharge{})
}

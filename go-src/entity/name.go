package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

type DbNameType uint8

const (
	// RoleNameType 角色名称类型
	RoleNameType DbNameType = 1
)

const (
	TbName = "names"
)

const (
	NameVal    db.TbCol = "val"
	NameTypeId db.TbCol = "type_id"
)

type Name struct {
	migrate.OneModel
	Val    string `gorm:"default:'';comment:名称"`
	TypeId uint8  `gorm:"default:0;comment:名称类型"`
}

func NewName(val string, typeId DbNameType) *Name {
	ret := &Name{}
	ret.ID = snowflake.GetId()
	ret.SetVal(val)
	ret.SetTypeId(uint8(typeId))
	return ret
}

func initName() {
	migrate.AutoMigrate(&Name{})

	syndb.RegLazyWithSmall(TbName, NameVal)
	syndb.RegLazyWithSmall(TbName, NameTypeId)
	syndb.RegLazyWithSmall(TbName, db.CreatedAtName)
	syndb.RegLazyWithSmall(TbName, db.UpdatedAtName)
}

func (receiver *Name) SetVal(val string) {
	receiver.Val = val
	syndb.AddDataToLazyChan(TbName, NameVal, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}
func (receiver *Name) SetTypeId(typeId uint8) {
	receiver.TypeId = typeId
	syndb.AddDataToLazyChan(TbName, NameTypeId, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: typeId,
	})
}
func (receiver *Name) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToLazyChan(TbName, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}
func (receiver *Name) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToLazyChan(TbName, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

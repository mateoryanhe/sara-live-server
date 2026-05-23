package entity

import (
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/constants/db"
	"xr-game-server/core/math"
	"xr-game-server/core/migrate"
	"xr-game-server/core/snowflake"
	"xr-game-server/core/syndb"
)

// 定义表名
const (
	TbPlayerRank db.TbName = "player_ranks"
)

// 定义列名
const (
	PlayerRankRoleId   db.TbCol = "role_id"
	PlayerRankVal      db.TbCol = "val"
	PlayerRankTypeId   db.TbCol = "type_id"
	PlayerRankLock     db.TbCol = "lock"
	PlayerRankLockTime db.TbCol = "lock_time"
)

type PlayerRank struct {
	migrate.OneModel
	RoleId   uint64     `gorm:"default:0;comment:角色id"`
	Val      uint64     `gorm:"default:0;comment:值"`
	TypeId   uint32     `gorm:"default:0;comment:类型"`
	Lock     bool       `gorm:"default:0;comment:禁止上榜"`
	LockTime *time.Time `gorm:"default:null;comment:锁定时间"`
}

func initPlayerRankChan() {
	syndb.RegLazyWithMiddle(TbPlayerRank, PlayerRankRoleId)
	syndb.RegLazyWithMiddle(TbPlayerRank, PlayerRankVal)
	syndb.RegLazyWithMiddle(TbPlayerRank, PlayerRankTypeId)
	syndb.RegLazyWithMiddle(TbPlayerRank, PlayerRankLock)
	syndb.RegLazyWithMiddle(TbPlayerRank, PlayerRankLockTime)

	syndb.RegLazyWithMiddle(TbPlayerRank, db.CreatedAtName)
	syndb.RegLazyWithMiddle(TbPlayerRank, db.UpdatedAtName)
}

func initRank() {
	initPlayerRankChan()
	migrate.AutoMigrate(&PlayerRank{})
}

func (receiver *PlayerRank) AddVal(val uint64) {
	ret := math.Add(receiver.Val, val)
	receiver.SetVal(ret)
}
func (receiver *PlayerRank) ReduceVal(val uint64) {
	if receiver.Val < val {
		receiver.SetVal(common.Zero)
	} else {
		receiver.SetVal(receiver.Val - val)
	}
}

func (receiver *PlayerRank) SetLock(lock bool) {
	receiver.Lock = lock
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbPlayerRank, PlayerRankLock, &syndb.ColData{
		ColVal: lock,
		IdVal:  receiver.ID,
	})
}
func (receiver *PlayerRank) SetLockTime(lockTime *time.Time) {
	receiver.LockTime = lockTime
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbPlayerRank, PlayerRankLockTime, &syndb.ColData{
		ColVal: lockTime,
		IdVal:  receiver.ID,
	})
}
func (receiver *PlayerRank) SetVal(val uint64) {
	receiver.Val = val
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToLazyChan(TbPlayerRank, PlayerRankVal, &syndb.ColData{
		ColVal: val,
		IdVal:  receiver.ID,
	})
}

func (receiver *PlayerRank) SetTypeId(typeId uint32) {
	receiver.TypeId = typeId
	syndb.AddDataToLazyChan(TbPlayerRank, PlayerRankTypeId, &syndb.ColData{
		ColVal: typeId,
		IdVal:  receiver.ID,
	})
}

func (receiver *PlayerRank) SetRoleId(roleId uint64) {
	receiver.RoleId = roleId
	syndb.AddDataToLazyChan(TbPlayerRank, PlayerRankRoleId, &syndb.ColData{
		ColVal: roleId,
		IdVal:  receiver.ID,
	})
}
func (receiver *PlayerRank) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToLazyChan(TbPlayerRank, db.CreatedAtName, &syndb.ColData{
		ColVal: val,
		IdVal:  receiver.ID,
	})
}
func (receiver *PlayerRank) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToLazyChan(TbPlayerRank, db.UpdatedAtName, &syndb.ColData{
		ColVal: val,
		IdVal:  receiver.ID,
	})
}

func NewPlayerRank(roleId uint64, typeId uint32) *PlayerRank {
	ret := &PlayerRank{}
	ret.ID = snowflake.GetId()
	ret.SetRoleId(roleId)
	ret.SetTypeId(typeId)
	ret.SetCreatedAt(time.Now())
	ret.SetUpdatedAt(time.Now())
	return ret
}

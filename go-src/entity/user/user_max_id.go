package entity

import (
	"sync/atomic"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbUserMaxId db.TbName = "user_max_ids"
)

const (
	// UserMaxIdDefaultID 用户ID起点(当前最大ID初始值;下次分配为 DefaultID+1)
	UserMaxIdDefaultID uint64 = 1000000
)

const (
	UserMaxIdMaxId db.TbCol = "max_id"
)

// UserMaxId 用户最大ID记录(全局单条;字段写入通过 syndb lazy 异步入库)
type UserMaxId struct {
	ID    uint64 `gorm:"primarykey" json:"id"`
	MaxId uint64 `gorm:"default:0;comment:当前最大ID" json:"maxId"`
}

// NewUserMaxId 构造最大ID记录(id 为表主键,通常为1)
func NewUserMaxId(id uint64) *UserMaxId {
	if id == 0 {
		id = 1
	}
	return &UserMaxId{ID: id}
}

// Add 原子增加最大ID,返回增加后的值
func (u *UserMaxId) Add(val uint64) uint64 {
	if val == 0 {
		return atomic.LoadUint64(&u.MaxId)
	}
	newVal := atomic.AddUint64(&u.MaxId, val)
	syndb.AddData(TbUserMaxId, UserMaxIdMaxId, &syndb.ColData{
		IdVal:  u.ID,
		ColVal: newVal,
	})
	return newVal
}

func (u *UserMaxId) SetMaxId(v uint64) {
	atomic.StoreUint64(&u.MaxId, v)
	syndb.AddData(TbUserMaxId, UserMaxIdMaxId, &syndb.ColData{
		IdVal:  u.ID,
		ColVal: v,
	})
}

func initUserMaxId() {
	syndb.RegLazy(TbUserMaxId, UserMaxIdMaxId)
	migrate.AutoMigrate(&UserMaxId{})
}

package entity

import (
	"sync/atomic"
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbUserMaxId db.TbName = "user_max_ids"
)

const (
	UserMaxIdDefaultID uint64 = 1000000
)

const (
	UserMaxIdMaxId db.TbCol = "max_id"
)

// UserMaxId 用户最大ID记录(全局单条,默认ID=1;字段写入通过 syndb lazy 异步入库)
type UserMaxId struct {
	migrate.OneModel
	MaxId uint64 `gorm:"default:0;comment:当前最大ID" json:"maxId"`
}

// NewUserMaxId 构造最大ID记录
func NewUserMaxId(id uint64) *UserMaxId {
	if id == 0 {
		id = UserMaxIdDefaultID
	}
	ret := &UserMaxId{}
	ret.ID = id
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	return ret
}

func (u *UserMaxId) SetMaxId(v uint64) {
	u.MaxId = v
	u.touchUpdatedAt()
	syndb.AddData(TbUserMaxId, UserMaxIdMaxId, &syndb.ColData{
		IdVal:  u.ID,
		ColVal: v,
	})
}

// Add 原子增加最大ID,返回增加后的值;val=0时仅返回当前值
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

func (u *UserMaxId) SetCreatedAt(v time.Time) {
	u.CreatedAt = v
	syndb.AddData(TbUserMaxId, db.CreatedAtName, &syndb.ColData{
		IdVal:  u.ID,
		ColVal: v,
	})
}

func (u *UserMaxId) SetUpdatedAt(v time.Time) {
	u.UpdatedAt = v
	syndb.AddData(TbUserMaxId, db.UpdatedAtName, &syndb.ColData{
		IdVal:  u.ID,
		ColVal: v,
	})
}

func (u *UserMaxId) touchUpdatedAt() {
	u.UpdatedAt = time.Now()
	syndb.AddData(TbUserMaxId, db.UpdatedAtName, &syndb.ColData{
		IdVal:  u.ID,
		ColVal: u.UpdatedAt,
	})
}

func initUserMaxId() {
	syndb.RegLazy(TbUserMaxId, db.CreatedAtName)
	syndb.RegLazy(TbUserMaxId, db.UpdatedAtName)
	syndb.RegLazy(TbUserMaxId, UserMaxIdMaxId)
	migrate.AutoMigrate(&UserMaxId{})
}

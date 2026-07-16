package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbUserExt db.TbName = "user_exts"
)

const (
	UserExtCanRank db.TbCol = "can_rank"
)

// UserExt 用户扩展信息(与用户一一对应,主键ID即用户ID)
type UserExt struct {
	migrate.OneModel
	CanRank bool `gorm:"default:1;comment:是否可上排行榜" json:"canRank"`
}

func NewUserExt(userId uint64) *UserExt {
	ret := &UserExt{}
	ret.ID = userId
	now := time.Now()
	ret.SetCreatedAt(now)
	ret.SetUpdatedAt(now)
	ret.SetCanRank(true)
	return ret
}

func (receiver *UserExt) SetCanRank(canRank bool) {
	receiver.CanRank = canRank
	receiver.SetUpdatedAt(time.Now())
	syndb.AddDataToQuickChan(TbUserExt, UserExtCanRank, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: canRank,
	})
}

func (receiver *UserExt) SetCreatedAt(val time.Time) {
	receiver.CreatedAt = val
	syndb.AddDataToQuickChan(TbUserExt, db.CreatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func (receiver *UserExt) SetUpdatedAt(val time.Time) {
	receiver.UpdatedAt = val
	syndb.AddDataToQuickChan(TbUserExt, db.UpdatedAtName, &syndb.ColData{
		IdVal:  receiver.ID,
		ColVal: val,
	})
}

func initUserExt() {
	syndb.RegQuick(TbUserExt, db.CreatedAtName)
	syndb.RegQuick(TbUserExt, db.UpdatedAtName)
	syndb.RegQuick(TbUserExt, UserExtCanRank)

	migrate.AutoMigrate(&UserExt{})
}

package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbAppToken db.TbName = "app_tokens"
)

const (
	AppTokenToken    db.TbCol = "token"
	AppTokenExpireAt db.TbCol = "expire_at"
)

// AppToken App端登录Token
type AppToken struct {
	ID       uint64    `gorm:"primarykey" json:"id"`
	Token    string    `gorm:"comment:Token值" json:"token"`
	ExpireAt time.Time `gorm:"index;comment:过期时间" json:"expireAt"`
}

// NewAppToken 构造一条App Token记录,字段写入通过 syndb 异步入库
func NewAppToken(userId uint64, token string, expireAt time.Time) *AppToken {
	ret := &AppToken{}
	ret.ID = userId
	ret.SetToken(token)
	ret.SetExpireAt(expireAt)
	return ret
}

func (t *AppToken) SetToken(v string) {
	t.Token = v

	syndb.AddDataToLazyChan(TbAppToken, AppTokenToken, &syndb.ColData{
		IdVal:  t.ID,
		ColVal: v,
	})
}

func (t *AppToken) SetExpireAt(v time.Time) {
	t.ExpireAt = v

	syndb.AddDataToLazyChan(TbAppToken, AppTokenExpireAt, &syndb.ColData{
		IdVal:  t.ID,
		ColVal: v,
	})
}

func initAppToken() {

	syndb.RegLazyWithLarge(TbAppToken, AppTokenToken)
	syndb.RegLazyWithLarge(TbAppToken, AppTokenExpireAt)
	migrate.AutoMigrate(&AppToken{})
}

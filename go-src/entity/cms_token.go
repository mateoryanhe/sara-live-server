package entity

import (
	"time"
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const (
	TbCmsToken db.TbName = "cms_tokens"
)

const (
	CmsTokenToken    db.TbCol = "token"
	CmsTokenExpireAt db.TbCol = "expire_at"
)

// CmsToken CMS端登录Token
type CmsToken struct {
	ID       uint64    `gorm:"primarykey" json:"id"`
	Token    string    `gorm:"comment:Token值" json:"token"`
	ExpireAt time.Time `gorm:"index;comment:过期时间" json:"expireAt"`
}

// NewCmsToken 构造一条CMS Token记录,字段写入通过 syndb 异步入库
func NewCmsToken(userId uint64, token string, expireAt time.Time) *CmsToken {
	ret := &CmsToken{}
	ret.ID = userId
	ret.SetToken(token)
	ret.SetExpireAt(expireAt)
	return ret
}

func (t *CmsToken) SetToken(v string) {
	t.Token = v
	syndb.AddDataToLazyChan(TbCmsToken, CmsTokenToken, &syndb.ColData{
		IdVal:  t.ID,
		ColVal: v,
	})
}

func (t *CmsToken) SetExpireAt(v time.Time) {
	t.ExpireAt = v
	syndb.AddDataToLazyChan(TbCmsToken, CmsTokenExpireAt, &syndb.ColData{
		IdVal:  t.ID,
		ColVal: v,
	})
}

func initCmsToken() {
	syndb.RegLazyWithLarge(TbCmsToken, CmsTokenToken)
	syndb.RegLazyWithLarge(TbCmsToken, CmsTokenExpireAt)
	migrate.AutoMigrate(&CmsToken{})
}

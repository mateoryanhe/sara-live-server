package entity

import (
	"strings"
	"time"

	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
	"xr-game-server/core/syndb"
)

const TbUserRegisterUrlAttr db.TbName = "user_register_url_attrs"

const (
	UserRegisterUrlAttrType     db.TbCol = "type"
	UserRegisterUrlAttrCh       db.TbCol = "ch"
	UserRegisterUrlAttrTjUserId db.TbCol = "tjuserid"
)

const (
	RegisterUrlAttrTypePWA = "pwa"
	RegisterUrlAttrTypeH5  = "h5"
	RegisterUrlAttrTypeAPK = "apk"
)

// RegisterUrlAttr 前端注册上报对象(合并 type/ch/tjuserid)
type RegisterUrlAttr struct {
	Type     string `json:"type" dc:"pwa/h5/apk"`
	Ch       string `json:"ch" dc:"渠道编号(URL)"`
	TjUserId string `json:"tjuserid" dc:"推广用户ID(URL)"`
}

// Normalize 清洗 type/ch/tjuserid; 全空则 ok=false(无上报不上库)
func (a *RegisterUrlAttr) Normalize() (typ, ch, tjUserId string, ok bool) {
	if a == nil {
		return "", "", "", false
	}
	typ = normalizeRegisterUrlAttrType(a.Type)
	ch = strings.TrimSpace(a.Ch)
	tjUserId = strings.TrimSpace(a.TjUserId)
	if typ == "" && ch == "" && tjUserId == "" {
		return "", "", "", false
	}
	return typ, ch, tjUserId, true
}

func normalizeRegisterUrlAttrType(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "pwa":
		return RegisterUrlAttrTypePWA
	case "h5":
		return RegisterUrlAttrTypeH5
	case "apk":
		return RegisterUrlAttrTypeAPK
	default:
		return ""
	}
}

// UserRegisterUrlAttr 注册时 H5/PWA URL 归因(与用户 1:1,主键=userId; syndb 缓冲入库)
type UserRegisterUrlAttr struct {
	migrate.OneModel
	Type     string `gorm:"size:16;default:'';index;comment:来源类型(pwa/h5/apk)" json:"type"`
	Ch       string `gorm:"size:64;default:'';index;comment:渠道编号(URL)" json:"ch"`
	TjUserId string `gorm:"column:tjuserid;size:64;default:'';index;comment:推广用户ID(URL)" json:"tjuserid"`
}

// HasValues 是否已有有效归因字段(缓存/库行判空用)
func (r *UserRegisterUrlAttr) HasValues() bool {
	if r == nil {
		return false
	}
	return strings.TrimSpace(r.Type) != "" || strings.TrimSpace(r.Ch) != "" || strings.TrimSpace(r.TjUserId) != ""
}

func NewUserRegisterUrlAttr(userId uint64, typ, ch, tjUserId string) *UserRegisterUrlAttr {
	row := &UserRegisterUrlAttr{}
	row.ID = userId
	now := time.Now()
	row.SetCreatedAt(now)
	row.SetUpdatedAt(now)
	row.SetType(typ)
	row.SetCh(ch)
	row.SetTjUserId(tjUserId)
	return row
}

func (r *UserRegisterUrlAttr) SetType(v string) {
	r.Type = v
	r.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserRegisterUrlAttr, UserRegisterUrlAttrType, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserRegisterUrlAttr) SetCh(v string) {
	r.Ch = v
	r.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserRegisterUrlAttr, UserRegisterUrlAttrCh, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserRegisterUrlAttr) SetTjUserId(v string) {
	r.TjUserId = v
	r.SetUpdatedAt(time.Now())
	syndb.AddData(TbUserRegisterUrlAttr, UserRegisterUrlAttrTjUserId, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserRegisterUrlAttr) SetCreatedAt(v time.Time) {
	r.CreatedAt = v
	syndb.AddData(TbUserRegisterUrlAttr, db.CreatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func (r *UserRegisterUrlAttr) SetUpdatedAt(v time.Time) {
	r.UpdatedAt = v
	syndb.AddData(TbUserRegisterUrlAttr, db.UpdatedAtName, &syndb.ColData{IdVal: r.ID, ColVal: v})
}

func initUserRegisterUrlAttr() {
	syndb.RegQuick(TbUserRegisterUrlAttr, db.CreatedAtName)
	syndb.RegQuick(TbUserRegisterUrlAttr, db.UpdatedAtName)
	syndb.RegQuick(TbUserRegisterUrlAttr, UserRegisterUrlAttrType)
	syndb.RegQuick(TbUserRegisterUrlAttr, UserRegisterUrlAttrCh)
	syndb.RegQuick(TbUserRegisterUrlAttr, UserRegisterUrlAttrTjUserId)
	migrate.AutoMigrate(&UserRegisterUrlAttr{})
}

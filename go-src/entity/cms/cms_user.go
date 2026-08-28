package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbCMSUser db.TbName = "cms_users"
)

// 列名
const (
	CMSUserName db.TbCol = "name"
)

const (
	CMSAdminTypeNone   uint8 = 0 // 非管理员
	CMSAdminTypeNormal uint8 = 1 // 普通管理员
	CMSAdminTypeSuper  uint8 = 2 // 超级管理员
)

type CMSUser struct {
	migrate.OneModel
	Name      string `gorm:"comment:名称" json:"name"`
	Pwd       string `gorm:"comment:密码" json:"pwd"`
	Status    uint8  `gorm:"default:0;comment:状态" json:"status"`
	Admin     bool   `gorm:"default:0;comment:是否是管理员" json:"admin"`
	AdminType uint8  `gorm:"default:0;comment:管理员类型 0非管理员 1普通管理员 2超级管理员" json:"adminType"`
	RoleId    uint64 `gorm:"default:0;comment:角色ID" json:"roleId"`
	Remark    string `gorm:"size:512;default:'';comment:备注" json:"remark"`
}

func NormalizeCMSUserAdminType(adminType uint8) uint8 {
	switch adminType {
	case CMSAdminTypeNormal, CMSAdminTypeSuper:
		return adminType
	default:
		return CMSAdminTypeNone
	}
}

func CMSUserAdminFromType(adminType uint8) bool {
	return adminType >= CMSAdminTypeNormal
}

func CMSUserIsAdmin(u *CMSUser) bool {
	if u == nil {
		return false
	}
	if u.AdminType >= CMSAdminTypeNormal {
		return true
	}
	return u.Admin
}

func CMSUserIsSuperAdmin(u *CMSUser) bool {
	if u == nil {
		return false
	}
	if u.AdminType == CMSAdminTypeSuper {
		return true
	}
	if u.AdminType == CMSAdminTypeNone && u.Admin {
		return true
	}
	return false
}

func InitCMSUser() {
	migrate.AutoMigrate(&CMSUser{})
}

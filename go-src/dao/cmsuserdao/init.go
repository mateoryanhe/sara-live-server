package cmsuserdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

func GetCMSUser(name string) *entity.CMSUser {
	//从数据库拉取数据
	var user *entity.CMSUser
	g.Model(string(entity.TbCMSUser)).Unscoped().Where(string(entity.CMSUserName), name).Scan(&user)
	if user != nil {
		return user
	}
	return nil
}

func InitCMSUser() {
	// 初始化CMS用户相关功能
	// 这里可以进行一些初始化操作
}

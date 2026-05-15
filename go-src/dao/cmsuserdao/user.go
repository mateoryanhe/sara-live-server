package cmsuserdao

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/cmsuserdto"
	"xr-game-server/entity"
)

// GetCMSUserById 根据ID获取CMS用户
func GetCMSUserById(id uint64) *entity.CMSUser {
	var user entity.CMSUser
	err := g.DB().Model(string(entity.TbCMSUser)).Where("id = ?", id).Scan(&user)
	if err != nil {
		return nil
	}
	return &user
}

// GetCMSUserByName 根据名称获取CMS用户
func GetCMSUserByName(name string) *entity.CMSUser {
	var user entity.CMSUser
	err := g.DB().Model(string(entity.TbCMSUser)).Where("name = ?", name).Scan(&user)
	if err != nil {
		return nil
	}
	return &user
}

// CreateCMSUser 创建CMS用户
func CreateCMSUser(user *entity.CMSUser) error {
	_, err := g.DB().Model(string(entity.TbCMSUser)).Save(user)
	if err == nil {
		return nil
	}
	return err
}

// UpdateCMSUser 更新CMS用户
func UpdateCMSUser(user *entity.CMSUser) error {
	return CreateCMSUser(user)
}

// DeleteCMSUser 删除CMS用户
func DeleteCMSUser(id uint64) error {
	_, err := g.DB().Model(string(entity.TbCMSUser)).WherePri(id).Delete()
	if err == nil {
		return nil
	}
	return err
}

// GetCMSUserList 获取CMS用户列表
func GetCMSUserList(req *cmsuserdto.CMSUserListReq) (int, []*cmsuserdto.CMSUserListRes) {
	sql := `select u.id, u.name, u.pwd, u.status, u.admin, u.role_id, r.name as role_name, u.created_at, u.updated_at
            from cms_users u
            left join cms_roles r on u.role_id = r.id
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*cmsuserdto.CMSUserListRes, 0)

	if req.Name != "" {
		sql += ` and u.name LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Name))
	}

	if req.Status > 0 {
		sql += ` and u.status = ?`
		param = append(param, req.Status)
	}

	if req.Admin {
		sql += ` and u.admin = ?`
		param = append(param, req.Admin)
	}

	sql += ` order by u.created_at desc`
	//获取总数
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}

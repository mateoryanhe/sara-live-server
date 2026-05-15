package cmsuserdao

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/cmsroledto"
	"xr-game-server/entity"
)

// GetRoleById 根据ID获取角色
func GetRoleById(id uint64) *entity.CMSRole {
	var role entity.CMSRole
	err := g.DB().Model(string(entity.TbCMSRole)).Where("id = ?", id).Scan(&role)
	if err != nil {
		return nil
	}
	return &role
}

// GetRoleByName 根据名称获取角色
func GetRoleByName(name string) *entity.CMSRole {
	var role entity.CMSRole
	err := g.DB().Model(string(entity.TbCMSRole)).Where("name = ?", name).Scan(&role)
	if err != nil {
		return nil
	}
	return &role
}

// CreateRole 创建角色
func CreateRole(role *entity.CMSRole) error {
	_, err := g.DB().Model(string(entity.TbCMSRole)).Save(role)
	if err == nil {
		return nil
	}
	return err
}

// UpdateRole 更新角色
func UpdateRole(role *entity.CMSRole) error {
	return CreateRole(role)
}

// DeleteRole 删除角色
func DeleteRole(id uint64) error {
	_, err := g.DB().Model(string(entity.TbCMSRole)).WherePri(id).Delete()
	if err == nil {
		return nil
	}
	return err
}

// GetRoleList 获取角色列表
func GetRoleList(req *cmsroledto.RoleListReq) (int, []*cmsroledto.RoleListRes) {
	sql := `select r.id, r.name, r.description, r.status, r.created_at, r.updated_at
            from cms_roles r
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*cmsroledto.RoleListRes, 0)

	if req.Name != "" {
		sql += ` and r.name LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Name))
	}

	sql += ` order by r.created_at desc`
	//获取总数
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}

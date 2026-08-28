package cmsuserdao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/core/str"
	"xr-game-server/dto/cmsuserdto"
	"xr-game-server/entity/cms"
)

var (
	cmsLoginUserCacheMgr     *cache.RowCache[*entity.CMSUser]
	cmsLoginUserMissCacheMgr *cache.RowCache[bool]
)

// GetCMSUser 登录按用户名查 CMS 用户: 先遍历缓存按 name 匹配,未命中再查库并写入缓存(key=user.id)
func GetCMSUser(name string) *entity.CMSUser {
	if name == "" || cmsLoginUserCacheMgr == nil {
		return nil
	}
	if user := findCMSUserInLoginCache(name); user != nil {
		return user
	}
	if hasCMSLoginUserMissCache(name) {
		return nil
	}
	var user entity.CMSUser
	err := g.DB().Model(string(entity.TbCMSUser)).Where(string(entity.CMSUserName), name).Scan(&user)
	if err != nil || user.ID == 0 {
		if cmsLoginUserMissCacheMgr != nil {
			cmsLoginUserMissCacheMgr.PublishRow(gctx.New(), name, true)
		}
		return nil
	}
	cmsLoginUserCacheMgr.PublishRow(gctx.New(), user.ID, &user)
	return &user
}

func hasCMSLoginUserMissCache(name string) bool {
	if cmsLoginUserMissCacheMgr == nil {
		return false
	}
	return cmsLoginUserMissCacheMgr.ContainsRow(gctx.New(), name)
}

func removeCMSLoginUserMissCache(name string) {
	if name == "" || cmsLoginUserMissCacheMgr == nil {
		return
	}
	cmsLoginUserMissCacheMgr.RemoveRow(gctx.New(), name)
}

func findCMSUserInLoginCache(name string) *entity.CMSUser {
	ctx := gctx.New()
	keys, err := cmsLoginUserCacheMgr.Keys(ctx)
	if err != nil || len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		user, ok := cmsLoginUserCacheMgr.GetRowCached(ctx, key)
		if !ok || user == nil || user.Name != name {
			continue
		}
		return user
	}
	return nil
}

// refreshCMSLoginUserCacheIfExists 用户变更后,若登录缓存存在则替换为最新值
func refreshCMSLoginUserCacheIfExists(user *entity.CMSUser) {
	if cmsLoginUserCacheMgr == nil || user == nil || user.ID == 0 {
		return
	}
	if _, ok := cmsLoginUserCacheMgr.GetRowCached(gctx.New(), user.ID); ok {
		cmsLoginUserCacheMgr.PublishRow(gctx.New(), user.ID, user)
	}
}

func removeCMSLoginUserCacheIfExists(userId uint64) {
	if cmsLoginUserCacheMgr == nil || userId == 0 {
		return
	}
	if _, ok := cmsLoginUserCacheMgr.GetRowCached(gctx.New(), userId); ok {
		cmsLoginUserCacheMgr.RemoveRow(gctx.New(), userId)
	}
}

// LoadCMSUserFromDB 根据 ID 从数据库加载 CMS 用户(不走缓存,供写操作使用)
func LoadCMSUserFromDB(id uint64) *entity.CMSUser {
	if id == 0 {
		return nil
	}
	var user entity.CMSUser
	err := g.DB().Model(string(entity.TbCMSUser)).Where("id = ?", id).Scan(&user)
	if err != nil || user.ID == 0 {
		return nil
	}
	return &user
}

// RemoveCMSLoginUserMissCache 清除用户名未命中缓存
func RemoveCMSLoginUserMissCache(name string) {
	removeCMSLoginUserMissCache(name)
}

func loadCMSUserFromDB(id uint64) *entity.CMSUser {
	return LoadCMSUserFromDB(id)
}

func syncCMSLoginUserCacheFromDB(userId uint64) {
	if userId == 0 {
		return
	}
	user := loadCMSUserFromDB(userId)
	if user == nil {
		removeCMSLoginUserCacheIfExists(userId)
		return
	}
	refreshCMSLoginUserCacheIfExists(user)
	removeCMSLoginUserMissCache(user.Name)
}

// GetCMSUserById 根据ID获取CMS用户(走 cmsLoginUserCacheMgr 缓存,key=user.id)
func GetCMSUserById(id uint64) *entity.CMSUser {
	if id == 0 || cmsLoginUserCacheMgr == nil {
		return nil
	}
	v := cmsLoginUserCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.CMSUser, error) {
		var user entity.CMSUser
		err := g.DB().Model(string(entity.TbCMSUser)).Where("id = ?", id).Scan(&user)
		if err != nil || user.ID == 0 {
			return nil, nil
		}
		return &user, nil
	})
	return v
}

// GetCMSUserByName 根据名称获取CMS用户
func GetCMSUserByName(name string) *entity.CMSUser {
	if name == "" {
		return nil
	}
	var user entity.CMSUser
	err := g.DB().Model(string(entity.TbCMSUser)).Where("name = ?", name).Scan(&user)
	if err != nil || user.ID == 0 {
		return nil
	}
	return &user
}

// CreateCMSUser 创建CMS用户
func CreateCMSUser(user *entity.CMSUser) error {
	_, err := g.DB().Model(string(entity.TbCMSUser)).Save(user)
	if err != nil {
		return err
	}
	if user != nil {
		removeCMSLoginUserMissCache(user.Name)
	}
	return nil
}

// UpdateCMSUser 更新CMS用户
func UpdateCMSUser(user *entity.CMSUser) error {
	_, err := g.DB().Model(string(entity.TbCMSUser)).Save(user)
	if err != nil {
		return err
	}
	syncCMSLoginUserCacheFromDB(user.ID)
	return nil
}

// ListCMSUserIdsByRoleId 按角色ID获取CMS用户ID列表
func ListCMSUserIdsByRoleId(roleId uint64) []uint64 {
	if roleId == 0 {
		return nil
	}
	users := make([]*entity.CMSUser, 0)
	_ = g.DB().Model(string(entity.TbCMSUser)).Fields("id").Where("role_id = ?", roleId).Scan(&users)
	ids := make([]uint64, 0, len(users))
	for _, user := range users {
		if user != nil && user.ID > 0 {
			ids = append(ids, user.ID)
		}
	}
	return ids
}

// DeleteCMSUser 删除CMS用户
func DeleteCMSUser(id uint64) error {
	user := GetCMSUserById(id)
	_, err := g.DB().Model(string(entity.TbCMSUser)).WherePri(id).Delete()
	if err != nil {
		return err
	}
	if user != nil {
		removeCMSLoginUserCacheIfExists(user.ID)
	}
	return nil
}

// GetCMSUserList 获取CMS用户列表
func GetCMSUserList(req *cmsuserdto.CMSUserListReq) (int, []*cmsuserdto.CMSUserListRes) {
	sql := `select u.id, u.name, u.pwd, u.status, u.admin, u.admin_type, u.role_id, u.remark, r.name as role_name, IFNULL(r.role_type, 0) as role_type, u.created_at, u.updated_at
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

	if req.Key != "" {
		likeKey := fmt.Sprintf("%%%s%%", req.Key)
		sql += ` and (u.name LIKE ? OR CAST(u.id AS CHAR) LIKE ?)`
		param = append(param, likeKey, likeKey)
	}

	if req.RoleId != "" {
		roleId, err := strconv.ParseUint(strings.TrimSpace(req.RoleId), 10, 64)
		if err == nil && roleId > 0 {
			sql += ` and u.role_id = ?`
			param = append(param, roleId)
		}
	}

	if req.Status > 0 {
		sql += ` and u.status = ?`
		param = append(param, req.Status)
	}

	if req.Admin {
		sql += ` and u.admin = ?`
		param = append(param, req.Admin)
	}

	if req.AdminType > 0 {
		sql += ` and u.admin_type = ?`
		param = append(param, req.AdminType)
	}

	if req.NonAdmin {
		sql += ` and u.admin_type = 0`
	}

	if req.RoleType > 0 {
		sql += ` and r.role_type = ?`
		param = append(param, req.RoleType)
	}

	sql += ` order by u.created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}

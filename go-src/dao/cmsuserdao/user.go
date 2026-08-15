package cmsuserdao

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/core/str"
	"xr-game-server/dto/cmsuserdto"
	"xr-game-server/entity/cms"
)

var (
	cmsLoginUserCacheMgr     *cache.CacheMgr
	cmsLoginUserMissCacheMgr *cache.CacheMgr
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
			cmsLoginUserMissCacheMgr.FlushCache(name, true)
		}
		return nil
	}
	cmsLoginUserCacheMgr.FlushCache(user.ID, &user)
	return &user
}

func hasCMSLoginUserMissCache(name string) bool {
	if cmsLoginUserMissCacheMgr == nil {
		return false
	}
	return cmsLoginUserMissCacheMgr.GetFromCache(name) != nil
}

func removeCMSLoginUserMissCache(name string) {
	if name == "" || cmsLoginUserMissCacheMgr == nil {
		return
	}
	_, _ = cmsLoginUserMissCacheMgr.Cache.Remove(gctx.New(), name)
}

func findCMSUserInLoginCache(name string) *entity.CMSUser {
	ctx := gctx.New()
	keys, err := cmsLoginUserCacheMgr.Cache.Keys(ctx)
	if err != nil || len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		v, err := cmsLoginUserCacheMgr.Cache.Get(ctx, key)
		if err != nil || v.IsNil() {
			continue
		}
		user, ok := v.Val().(*entity.CMSUser)
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
	if cmsLoginUserCacheMgr.GetFromCache(user.ID) != nil {
		cmsLoginUserCacheMgr.FlushCache(user.ID, user)
	}
}

func removeCMSLoginUserCacheIfExists(userId uint64) {
	if cmsLoginUserCacheMgr == nil || userId == 0 {
		return
	}
	if cmsLoginUserCacheMgr.GetFromCache(userId) != nil {
		_, _ = cmsLoginUserCacheMgr.Cache.Remove(gctx.New(), userId)
	}
}

// GetCMSUserById 根据ID获取CMS用户(走 cmsLoginUserCacheMgr 缓存,key=user.id)
func GetCMSUserById(id uint64) *entity.CMSUser {
	if id == 0 || cmsLoginUserCacheMgr == nil {
		return nil
	}
	v := cmsLoginUserCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var user entity.CMSUser
		err = g.DB().Model(string(entity.TbCMSUser)).Where("id = ?", id).Scan(&user)
		if err != nil || user.ID == 0 {
			return nil, nil
		}
		return &user, nil
	})
	if v == nil {
		return nil
	}
	user, _ := v.(*entity.CMSUser)
	return user
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
	refreshCMSLoginUserCacheIfExists(user)
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
	sql := `select u.id, u.name, u.pwd, u.status, u.admin, u.role_id, u.remark, r.name as role_name, u.created_at, u.updated_at
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

	if req.Status > 0 {
		sql += ` and u.status = ?`
		param = append(param, req.Status)
	}

	if req.Admin {
		sql += ` and u.admin = ?`
		param = append(param, req.Admin)
	}

	sql += ` order by u.created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}

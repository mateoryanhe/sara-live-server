package cmsuserdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/cms"
)

var rolePermissionCacheMgr *cache.CacheMgr

func initRolePermissionCache() {
	rolePermissionCacheMgr = cache.NewCacheMgr()
}

func loadRolePermissionsFromDB(roleId uint64) []*entity.Permission {
	list := make([]*entity.Permission, 0)
	_ = g.DB().Model(string(entity.TbPermission)).Where("role_id = ?", roleId).Scan(&list)
	return list
}

// GetGetPermissionList 按角色获取权限列表(走缓存)
func GetGetPermissionList(roleId uint64) []*entity.Permission {
	if roleId == 0 || rolePermissionCacheMgr == nil {
		return []*entity.Permission{}
	}
	v := rolePermissionCacheMgr.GetData(roleId, func(ctx context.Context) (value interface{}, err error) {
		return loadRolePermissionsFromDB(roleId), nil
	})
	if v == nil {
		return []*entity.Permission{}
	}
	list, _ := v.([]*entity.Permission)
	if list == nil {
		return []*entity.Permission{}
	}
	return list
}

// CheckCmsApiPermission CMS 中间件按 URL 校验用户接口权限
func CheckCmsApiPermission(userId uint64, apiPath string) bool {
	user := GetCMSUserById(userId)
	if user == nil {
		return false
	}
	if user.Admin {
		return true
	}
	return RoleHasApiPath(user.RoleId, apiPath)
}

// RoleHasApiPath 角色是否拥有该接口路径权限
func RoleHasApiPath(roleId uint64, apiPath string) bool {
	if roleId == 0 || apiPath == "" {
		return false
	}
	if roleHasApiPathExact(roleId, apiPath) {
		return true
	}
	for _, alias := range cmsApiPermissionAliasPaths(apiPath) {
		if roleHasApiPathExact(roleId, alias) {
			return true
		}
	}
	return false
}

func roleHasApiPathExact(roleId uint64, apiPath string) bool {
	for _, p := range GetGetPermissionList(roleId) {
		if p != nil && p.ApiPath == apiPath {
			return true
		}
	}
	return false
}

func cmsApiPermissionAliasPaths(apiPath string) []string {
	switch apiPath {
	case "/account/getUserDetail":
		return []string{"/account/getUserInfo"}
	case "/account/getAnchorDailyEffectiveLiveList":
		return []string{"/account/getAnchorDetail"}
	case "/guild/getGuildDailyEffectiveLiveList":
		return []string{"/guild/getGuildDetail"}
	case "/guild/cmsGuildAnchorDailyEffectiveLiveList":
		return []string{"/guild/cmsGuildAnchorIncomeSettlementLogList", "/guild/getGuildDetail"}
	case "/guild/setGuildAnchorType":
		return []string{"/account/setPlatformAnchorType"}
	case "/cmsExport/getJob":
		return []string{"/cmsExport/submitJob"}
	case "/cmsExport/deleteExport":
		return []string{"/cmsExport/submitJob"}
	case "/account/getAnchorList":
		return []string{"/liveRecord/cmsLiveRecordList", "/account/getAnchorList"}
	default:
		return nil
	}
}

func SavePermissions(data []*entity.Permission) {
	if len(data) == 0 {
		return
	}
	roleId := data[0].RoleId
	DeleteRolePermissions(roleId)
	g.DB().Model(string(entity.TbPermission)).Save(data)
	refreshRolePermissionCacheIfExists(roleId)
}

// DeleteRolePermissions 删除角色关联的权限
func DeleteRolePermissions(roleId uint64) error {
	_, err := g.DB().Model(string(entity.TbPermission)).Delete("role_id = ?", roleId)
	if err != nil {
		return err
	}
	removeRolePermissionCacheIfExists(roleId)
	return nil
}

// refreshRolePermissionCacheIfExists 角色权限变更后,若缓存存在则替换为最新值
func refreshRolePermissionCacheIfExists(roleId uint64) {
	if rolePermissionCacheMgr == nil || roleId == 0 {
		return
	}
	if rolePermissionCacheMgr.GetFromCache(roleId) != nil {
		rolePermissionCacheMgr.FlushCache(roleId, loadRolePermissionsFromDB(roleId))
	}
}

func removeRolePermissionCacheIfExists(roleId uint64) {
	if rolePermissionCacheMgr == nil || roleId == 0 {
		return
	}
	if rolePermissionCacheMgr.GetFromCache(roleId) != nil {
		_, _ = rolePermissionCacheMgr.Cache.Remove(gctx.New(), roleId)
	}
}

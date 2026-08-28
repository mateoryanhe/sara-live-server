package cmsuserdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/cms"
)

var rolePermissionCacheMgr *cache.ListCache[*entity.Permission]

func initRolePermissionCache() {
	rolePermissionCacheMgr = cache.NewListCache[*entity.Permission]()
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
	v := rolePermissionCacheMgr.MustGetList(gctx.New(), roleId, func(ctx context.Context) ([]*entity.Permission, error) {
		return loadRolePermissionsFromDB(roleId), nil
	})
	return v
}

// CheckCmsApiPermission CMS 中间件按 URL 校验用户接口权限
func CheckCmsApiPermission(userId uint64, apiPath string) bool {
	user := GetCMSUserById(userId)
	if user == nil {
		return false
	}
	if entity.CMSUserIsAdmin(user) {
		if !entity.CMSUserIsSuperAdmin(user) && isDataSyncApiPath(apiPath) {
			return false
		}
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
	case "/gamePlatform/cmsGameStartLink":
		return []string{"/gamePlatform/gameShelfList"}
	case "/account/getUserDetail":
		return []string{"/account/getUserInfo", "/liveRecord/cmsLiveRecordList", "/liveRevenueLog/cmsLiveRevenueLogList"}
	case "/account/getAnchorDailyEffectiveLiveList":
		return []string{"/account/getAnchorDetail"}
	case "/guild/getGuildDailyEffectiveLiveList":
		return []string{"/guild/getGuildDetail"}
	case "/guild/guildList":
		return []string{"/guild/cmsGuildAnchorDailyEffectiveLiveList", "/guild/cmsMyGuildAnchorDailyEffectiveLiveList", "/liveRecord/cmsLiveRecordList", "/liveRecord/cmsDailyEffectiveLiveList", "/liveRevenueLog/cmsLiveRevenueLogList"}
	case "/guild/cmsGuildAnchorDailyEffectiveLiveList":
		return []string{"/guild/cmsGuildAnchorIncomeSettlementLogList", "/guild/getGuildDetail", "/guild/guildList", "/account/getAnchorList"}
	case "/guild/setGuildAnchorType":
		return []string{"/account/setPlatformAnchorType"}
	case "/cmsExport/getJob":
		return []string{"/cmsExport/submitJob"}
	case "/cmsExport/deleteExport":
		return []string{"/cmsExport/submitJob"}
	case "/cmsuser/createGuildCMSUser":
		return []string{"/role/roleList"}
	case "/role/roleList":
		return []string{"/cmsuser/createGuildCMSUser"}
	case "/guild/cmsMyGuildAnchorDailyEffectiveLiveList":
		return []string{"/guild/getMyGuildProfile", "/guild/getMyOwnedGuildAnchorList"}
	case "/guild/getMyOwnedGuildAnchorList":
		return []string{"/guild/getMyGuildProfile"}
	case "/guild/getMyGuildAnchorList":
		return []string{"/guild/getMyGuildProfile", "/account/getAnchorList"}
	case "/guild/getMyGuildAnchorDailyEffectiveLiveList":
		return []string{"/guild/getMyGuildProfile", "/account/getAnchorDailyEffectiveLiveList"}
	case "/account/getAnchorList":
		return []string{"/liveRecord/cmsLiveRecordList", "/liveRecord/cmsDailyEffectiveLiveList", "/liveRevenueLog/cmsLiveRevenueLogList", "/account/getAnchorList", "/guild/cmsGuildAnchorDailyEffectiveLiveList", "/guild/cmsMyGuildAnchorDailyEffectiveLiveList", "/guild/getMyOwnedGuildAnchorList"}
	case "/liveRecord/cmsDailyEffectiveLiveList":
		return []string{"/liveRecord/cmsLiveRecordList", "/account/getAnchorList", "/guild/guildList"}
	case "/liveRecord/cmsWeeklyUnsettledLiveList":
		return []string{"/liveRecord/cmsLiveRecordList", "/account/getAnchorList", "/guild/guildList"}
	case "/liveRevenueLog/cmsLiveRevenueLogList":
		return []string{"/liveRecord/cmsLiveRecordList", "/account/getAnchorList", "/guild/guildList"}
	case "/liveRecord/cmsLiveRecordList":
		return []string{"/liveRevenueLog/cmsLiveRevenueLogList"}
	case "/account/getAnchorDetail":
		return []string{"/liveRecord/cmsLiveRecordList", "/liveRevenueLog/cmsLiveRevenueLogList"}
	case "/shortVideo/shortVideoStorageStat":
		return []string{"/shortVideo/shortVideoList"}
	case "/shortVideo/onShelfShortVideo", "/shortVideo/offShelfShortVideo":
		return []string{"/shortVideo/updateShortVideo"}
	case "/rechargeCfg/rechargeCfgList":
		return []string{"/rechargeOrder/manualCreateOrder"}
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
	if _, ok := rolePermissionCacheMgr.GetListCached(gctx.New(), roleId); ok {
		rolePermissionCacheMgr.PublishList(gctx.New(), roleId, loadRolePermissionsFromDB(roleId))
	}
}

func removeRolePermissionCacheIfExists(roleId uint64) {
	if rolePermissionCacheMgr == nil || roleId == 0 {
		return
	}
	if _, ok := rolePermissionCacheMgr.GetListCached(gctx.New(), roleId); ok {
		rolePermissionCacheMgr.RemoveList(gctx.New(), roleId)
	}
}

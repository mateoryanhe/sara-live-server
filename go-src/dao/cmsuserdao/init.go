package cmsuserdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var cmsLoginUserCacheMgr *cache.CacheMgr

func InitCMSUser() {
	cmsLoginUserCacheMgr = cache.NewCacheMgr()
}

// GetCMSUser 登录按用户名查 CMS 用户(走缓存保护)
func GetCMSUser(name string) *entity.CMSUser {
	if name == "" || cmsLoginUserCacheMgr == nil {
		return nil
	}
	v := cmsLoginUserCacheMgr.GetData(name, func(ctx context.Context) (value interface{}, err error) {
		var user *entity.CMSUser
		_ = g.Model(string(entity.TbCMSUser)).Where(string(entity.CMSUserName), name).Scan(&user)
		if user == nil || user.ID == 0 {
			return nil, nil
		}
		return user, nil
	})
	if v == nil {
		return nil
	}
	user, _ := v.(*entity.CMSUser)
	return user
}

// refreshCMSLoginUserCacheIfExists 用户变更后,若登录缓存存在则替换为最新值
func refreshCMSLoginUserCacheIfExists(user *entity.CMSUser, oldName string) {
	if cmsLoginUserCacheMgr == nil || user == nil || user.ID == 0 {
		return
	}
	ctx := gctx.New()
	if oldName != "" && oldName != user.Name {
		if cmsLoginUserCacheMgr.GetFromCache(oldName) != nil {
			_, _ = cmsLoginUserCacheMgr.Cache.Remove(ctx, oldName)
		}
	}
	if cmsLoginUserCacheMgr.GetFromCache(user.Name) != nil {
		cmsLoginUserCacheMgr.FlushCache(user.Name, user)
	}
}

func removeCMSLoginUserCacheIfExists(name string) {
	if cmsLoginUserCacheMgr == nil || name == "" {
		return
	}
	if cmsLoginUserCacheMgr.GetFromCache(name) != nil {
		_, _ = cmsLoginUserCacheMgr.Cache.Remove(gctx.New(), name)
	}
}

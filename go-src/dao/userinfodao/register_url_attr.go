package userinfodao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/user"
)

var userRegisterUrlAttrCacheMgr *cache.RowCache[*entity.UserRegisterUrlAttr]

func initUserRegisterUrlAttrDao() {
	userRegisterUrlAttrCacheMgr = cache.NewRowCache[*entity.UserRegisterUrlAttr]()
}

// FindUserRegisterUrlAttrByUserId 只读已有有效归因.
// 缓存命中(含负缓存 val=nil)不打库; 库无记录时 PublishRow(nil) 占位,避免穿透.
func FindUserRegisterUrlAttrByUserId(userId uint64) *entity.UserRegisterUrlAttr {
	if userId == 0 || userRegisterUrlAttrCacheMgr == nil {
		return nil
	}
	ctx := gctx.New()
	if cached, ok := userRegisterUrlAttrCacheMgr.GetRowCached(ctx, userId); ok {
		if cached != nil && cached.HasValues() {
			return cached
		}
		// ok && (nil 或空行) = 负缓存命中,不再查库
		if cached != nil && !cached.HasValues() {
			// 历史空行统一收成 nil 负缓存
			userRegisterUrlAttrCacheMgr.PublishRow(ctx, userId, nil)
		}
		return nil
	}
	var data *entity.UserRegisterUrlAttr
	_ = g.Model(string(entity.TbUserRegisterUrlAttr)).Unscoped().Where(g.Map{
		string(db.IdName): userId,
	}).Scan(&data)
	if data == nil || data.ID == 0 || !data.HasValues() {
		// 库无有效记录: 负缓存占位,后续 Find 不再打库
		userRegisterUrlAttrCacheMgr.PublishRow(ctx, userId, nil)
		return nil
	}
	userRegisterUrlAttrCacheMgr.PublishRow(ctx, data.ID, data)
	return data
}

// PublishUserRegisterUrlAttr 原地修改后刷新缓存(会覆盖负缓存 nil)
func PublishUserRegisterUrlAttr(data *entity.UserRegisterUrlAttr) {
	if data == nil || data.ID == 0 || userRegisterUrlAttrCacheMgr == nil {
		return
	}
	userRegisterUrlAttrCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

// MarkRegisterUrlAttrAbsent 注册时未上报归因: 只写负缓存(val=nil),不上库;后续 Find 命中负缓存不打库.
func MarkRegisterUrlAttrAbsent(userId uint64) {
	if userId == 0 || userRegisterUrlAttrCacheMgr == nil {
		return
	}
	userRegisterUrlAttrCacheMgr.PublishRow(gctx.New(), userId, nil)
}

// SaveRegisterUrlAttr 仅在有有效归因数据时入库;无数据跳过;已存在有效记录不覆盖.
// 新建后 Publish 覆盖可能存在的负缓存.
func SaveRegisterUrlAttr(userId uint64, typ, ch, tjUserId string) *entity.UserRegisterUrlAttr {
	if userId == 0 {
		return nil
	}
	typ = strings.TrimSpace(typ)
	ch = strings.TrimSpace(ch)
	tjUserId = strings.TrimSpace(tjUserId)
	if typ == "" && ch == "" && tjUserId == "" {
		return nil
	}
	if existing := FindUserRegisterUrlAttrByUserId(userId); existing != nil {
		return existing
	}
	row := entity.NewUserRegisterUrlAttr(userId, typ, ch, tjUserId)
	PublishUserRegisterUrlAttr(row)
	return row
}

// SaveRegisterUrlAttrFromDevice 注册时写入 URL 归因; deviceInfo.registerAttr 为空则不上库.
func SaveRegisterUrlAttrFromDevice(userId uint64, info *entity.DeviceInfo) *entity.UserRegisterUrlAttr {
	if userId == 0 || info == nil {
		return nil
	}
	typ, ch, tjUserId, ok := info.RegisterAttr.Normalize()
	if !ok {
		return nil
	}
	return SaveRegisterUrlAttr(userId, typ, ch, tjUserId)
}

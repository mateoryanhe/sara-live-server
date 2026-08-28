package userinfodao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strings"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	userentity "xr-game-server/entity/user"
)

var userInfoCacheMgr *cache.RowCache[*userentity.UserInfo]
var shareCodeUserIdCacheMgr *cache.RowCache[uint64]

func InitUserInfoDao() {
	userInfoCacheMgr = cache.NewRowCache[*userentity.UserInfo]()
	shareCodeUserIdCacheMgr = cache.NewRowCache[uint64]()
	initUserCumulativeStatDao()
	initUserExtDao()
	initUserRechargeCfgFirstRechargeDao()
}

// GetUserIdByShareCode 根据分享码获取玩家ID,不存在则返回 0
func GetUserIdByShareCode(shareCode string) uint64 {
	if shareCode == "" {
		return 0
	}
	return shareCodeUserIdCacheMgr.MustGetRow(gctx.New(), shareCode, func(ctx context.Context) (uint64, error) {
		var userId uint64
		err := g.Model(string(userentity.TbUserInfo)).Unscoped().Where(g.Map{
			string(userentity.UserInfoShareCode): shareCode,
		}).Fields(string(db.IdName)).Scan(&userId)
		return userId, err
	})
}

// PublishUserInfo 原地修改 UserInfo 后调用,刷新缓存条目.
func PublishUserInfo(data *userentity.UserInfo) {
	if data == nil || data.ID == 0 || userInfoCacheMgr == nil {
		return
	}
	userInfoCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

// GetUserInfoByUserId 根据用户ID获取用户基础信息,命中不了缓存从数据库拉取,数据库不存在则新建
func GetUserInfoByUserId(userId uint64) *userentity.UserInfo {
	return userInfoCacheMgr.MustGetRow(gctx.New(), userId, func(ctx context.Context) (*userentity.UserInfo, error) {
		var data *userentity.UserInfo
		_ = g.Model(string(userentity.TbUserInfo)).Unscoped().Where(g.Map{
			string(db.IdName): userId,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		//数据库没有,新建一条
		newData := userentity.NewUserInfo(userId)
		return newData, nil
	})
}

// GetUserInfoFromMemory 仅从内存缓存读取用户基础信息,未命中返回 nil
func GetUserInfoFromMemory(userId uint64) *userentity.UserInfo {
	if userId == 0 || userInfoCacheMgr == nil {
		return nil
	}
	v, _ := userInfoCacheMgr.GetRowCached(gctx.New(), userId)
	return v
}

// GetNicknameMapByUserIds 批量查询用户昵称(CMS列表等场景使用)
func GetNicknameMapByUserIds(userIds []uint64) map[uint64]string {
	ret := make(map[uint64]string)
	if len(userIds) == 0 {
		return ret
	}
	uniqueIds := make([]uint64, 0, len(userIds))
	seen := make(map[uint64]struct{}, len(userIds))
	for _, id := range userIds {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}
	if len(uniqueIds) == 0 {
		return ret
	}
	rows := make([]*userentity.UserInfo, 0, len(uniqueIds))
	ctx := gctx.New()
	_ = g.Model(string(userentity.TbUserInfo)).Ctx(ctx).Unscoped().
		Fields(string(db.IdName), string(userentity.UserInfoNickname)).
		WhereIn(string(db.IdName), uniqueIds).
		Scan(&rows)
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		ret[row.ID] = row.Nickname
	}
	return ret
}

// GetUserProfileMapByUserIds 批量查询用户昵称与头像
func GetUserProfileMapByUserIds(userIds []uint64) map[uint64]*userentity.UserInfo {
	ret := make(map[uint64]*userentity.UserInfo)
	if len(userIds) == 0 {
		return ret
	}
	uniqueIds := make([]uint64, 0, len(userIds))
	seen := make(map[uint64]struct{}, len(userIds))
	for _, id := range userIds {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIds = append(uniqueIds, id)
	}
	if len(uniqueIds) == 0 {
		return ret
	}
	rows := make([]*userentity.UserInfo, 0, len(uniqueIds))
	ctx := gctx.New()
	_ = g.Model(string(userentity.TbUserInfo)).Ctx(ctx).Unscoped().
		Fields(string(db.IdName), string(userentity.UserInfoNickname), string(userentity.UserInfoAvatar)).
		WhereIn(string(db.IdName), uniqueIds).
		Scan(&rows)
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		ret[row.ID] = row
	}
	return ret
}

// GetUserIdsByNicknameKeyword 根据昵称关键字从数据库模糊匹配用户ID(CMS列表筛选等场景)
func GetUserIdsByNicknameKeyword(keyword string) map[uint64]struct{} {
	ret := make(map[uint64]struct{})
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return ret
	}
	rows := make([]*userentity.UserInfo, 0)
	ctx := gctx.New()
	_ = g.Model(string(userentity.TbUserInfo)).Ctx(ctx).Unscoped().
		Fields(string(db.IdName)).
		WhereLike(string(userentity.UserInfoNickname), "%"+keyword+"%").
		Scan(&rows)
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		ret[row.ID] = struct{}{}
	}
	return ret
}

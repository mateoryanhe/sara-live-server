package userinfodao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var userInfoCacheMgr *cache.CacheMgr
var shareCodeUserIdCacheMgr *cache.CacheMgr

func InitUserInfoDao() {
	userInfoCacheMgr = cache.NewCacheMgr()
	shareCodeUserIdCacheMgr = cache.NewCacheMgr()
	initUserCumulativeStatDao()
}

// GetUserIdByShareCode 根据分享码获取玩家ID,不存在则返回 0
func GetUserIdByShareCode(shareCode string) uint64 {
	if shareCode == "" {
		return 0
	}
	cacheData := shareCodeUserIdCacheMgr.GetData(shareCode, func(ctx context.Context) (value interface{}, err error) {
		var userId uint64
		err = g.Model(string(entity.TbUserInfo)).Unscoped().Where(g.Map{
			string(entity.UserInfoShareCode): shareCode,
		}).Fields(string(db.IdName)).Scan(&userId)
		return userId, err
	})
	if cacheData == nil {
		return 0
	}
	if userId, ok := cacheData.(uint64); ok {
		return userId
	}
	return 0
}

// GetUserInfoByUserId 根据用户ID获取用户基础信息,命中不了缓存从数据库拉取,数据库不存在则新建
func GetUserInfoByUserId(userId uint64) *entity.UserInfo {
	cacheData := userInfoCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		var data *entity.UserInfo
		err = g.Model(string(entity.TbUserInfo)).Unscoped().Where(g.Map{
			string(db.IdName): userId,
		}).Scan(&data)
		if data != nil {
			return data, nil
		}
		//数据库没有,新建一条
		newData := entity.NewUserInfo(userId)
		return newData, nil
	})
	return cacheData.(*entity.UserInfo)
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
	rows := make([]*entity.UserInfo, 0, len(uniqueIds))
	ctx := gctx.New()
	_ = g.Model(string(entity.TbUserInfo)).Ctx(ctx).Unscoped().
		Fields(string(db.IdName), string(entity.UserInfoNickname)).
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
func GetUserProfileMapByUserIds(userIds []uint64) map[uint64]*entity.UserInfo {
	ret := make(map[uint64]*entity.UserInfo)
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
	rows := make([]*entity.UserInfo, 0, len(uniqueIds))
	ctx := gctx.New()
	_ = g.Model(string(entity.TbUserInfo)).Ctx(ctx).Unscoped().
		Fields(string(db.IdName), string(entity.UserInfoNickname), string(entity.UserInfoAvatar)).
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

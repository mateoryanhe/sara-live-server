package guilddao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var (
	// memberCacheMgr 按成员ID缓存单个成员对象
	memberCacheMgr *cache.CacheMgr
	// guildMembersCacheMgr 按工会ID缓存成员列表 []*entity.LiveGuildMember
	guildMembersCacheMgr *cache.CacheMgr
)

func InitGuildMemberDao() {
	memberCacheMgr = cache.NewCacheMgr()
	guildMembersCacheMgr = cache.NewCacheMgr()
}

// GetMember 按成员ID获取(走缓存)
func GetMember(memberId uint64) *entity.LiveGuildMember {
	v := memberCacheMgr.GetData(memberId, func(ctx context.Context) (value interface{}, err error) {
		var m *entity.LiveGuildMember
		_ = g.Model(string(entity.TbLiveGuildMember)).Where("id = ?", memberId).Scan(&m)
		return m, nil
	})
	if v == nil {
		return nil
	}
	m, _ := v.(*entity.LiveGuildMember)
	return m
}

// GetMembersByGuild 按工会ID获取成员列表(走缓存,包含所有状态记录)
func GetMembersByGuild(guildId uint64) []*entity.LiveGuildMember {
	v := guildMembersCacheMgr.GetData(guildId, func(ctx context.Context) (value interface{}, err error) {
		list := make([]*entity.LiveGuildMember, 0)
		_ = g.Model(string(entity.TbLiveGuildMember)).
			Where("guild_id = ?", guildId).
			Order("created_at asc").
			Scan(&list)
		return list, nil
	})
	if v == nil {
		return make([]*entity.LiveGuildMember, 0)
	}
	list, _ := v.([]*entity.LiveGuildMember)
	if list == nil {
		return make([]*entity.LiveGuildMember, 0)
	}
	return list
}

// FindMemberInGuild 在工会成员缓存中查找指定用户(可按状态过滤; statuses 为空表示不过滤)
func FindMemberInGuild(guildId, userId uint64, statuses ...uint8) *entity.LiveGuildMember {
	list := GetMembersByGuild(guildId)
	for _, m := range list {
		if m.UserId != userId {
			continue
		}
		if len(statuses) == 0 {
			return m
		}
		for _, s := range statuses {
			if m.Status == s {
				return m
			}
		}
	}
	return nil
}

// AddMemberToCache 将新成员加入缓存(单成员缓存 + 工会成员列表缓存)
func AddMemberToCache(m *entity.LiveGuildMember) {
	memberCacheMgr.FlushCache(m.ID, m)
	list := GetMembersByGuild(m.GuildId)
	list = append(list, m)
	guildMembersCacheMgr.FlushCache(m.GuildId, list)
}

// FlushMemberCache 刷新单成员缓存(用于状态变更后)
func FlushMemberCache(m *entity.LiveGuildMember) {
	memberCacheMgr.FlushCache(m.ID, m)
}

// RemoveGuildMembersCache 移除整个工会的成员列表缓存
func RemoveGuildMembersCache(guildId uint64) {
	if guildMembersCacheMgr == nil {
		return
	}
	_, _ = guildMembersCacheMgr.Cache.Remove(gctx.New(), guildId)
}

// RemoveMemberCache 移除单成员缓存
func RemoveMemberCache(memberId uint64) {
	if memberCacheMgr == nil {
		return
	}
	_, _ = memberCacheMgr.Cache.Remove(gctx.New(), memberId)
}

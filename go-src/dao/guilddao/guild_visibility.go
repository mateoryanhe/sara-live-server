package guilddao

import (
	"github.com/gogf/gf/v2/frame/g"
	liveentity "xr-game-server/entity/live"
)

// AddGuildVisibility 写入工会可见性(同一工会同一用户已存在则忽略)
func AddGuildVisibility(guildId, cmsUserId uint64) error {
	if guildId == 0 || cmsUserId == 0 {
		return nil
	}
	if HasGuildVisibility(guildId, cmsUserId) {
		return nil
	}
	row := liveentity.NewLiveGuildVisibility(guildId, cmsUserId)
	_, err := g.DB().Model(string(liveentity.TbLiveGuildVisibility)).Save(row)
	return err
}

// HasGuildVisibility 是否已有可见性记录
func HasGuildVisibility(guildId, cmsUserId uint64) bool {
	if guildId == 0 || cmsUserId == 0 {
		return false
	}
	count, err := g.DB().Model(string(liveentity.TbLiveGuildVisibility)).
		Where(string(liveentity.LiveGuildVisibilityGuildId), guildId).
		Where(string(liveentity.LiveGuildVisibilityCmsUserId), cmsUserId).
		Count()
	return err == nil && count > 0
}

// ListVisibleGuildIds 查询 CMS 用户可见的工会 ID 列表
func ListVisibleGuildIds(cmsUserId uint64) []uint64 {
	if cmsUserId == 0 {
		return nil
	}
	type row struct {
		GuildId uint64 `json:"guildId"`
	}
	rows := make([]row, 0)
	_ = g.DB().Model(string(liveentity.TbLiveGuildVisibility)).
		Fields(string(liveentity.LiveGuildVisibilityGuildId)).
		Where(string(liveentity.LiveGuildVisibilityCmsUserId), cmsUserId).
		Scan(&rows)
	ids := make([]uint64, 0, len(rows))
	seen := make(map[uint64]struct{}, len(rows))
	for _, r := range rows {
		if r.GuildId == 0 {
			continue
		}
		if _, ok := seen[r.GuildId]; ok {
			continue
		}
		seen[r.GuildId] = struct{}{}
		ids = append(ids, r.GuildId)
	}
	return ids
}

// ListVisibilitiesByCmsUserId 按 CMS 用户查询可见性记录
func ListVisibilitiesByCmsUserId(cmsUserId uint64) []*liveentity.LiveGuildVisibility {
	if cmsUserId == 0 {
		return nil
	}
	rows := make([]*liveentity.LiveGuildVisibility, 0)
	_ = g.DB().Model(string(liveentity.TbLiveGuildVisibility)).
		Where(string(liveentity.LiveGuildVisibilityCmsUserId), cmsUserId).
		Order("created_at asc").
		Scan(&rows)
	return rows
}

// DeleteGuildVisibilities 删除工会全部可见性记录
func DeleteGuildVisibilities(guildId uint64) error {
	if guildId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(liveentity.TbLiveGuildVisibility)).
		Where(string(liveentity.LiveGuildVisibilityGuildId), guildId).
		Delete()
	return err
}

// RemoveGuildVisibility 删除单条可见性
func RemoveGuildVisibility(guildId, cmsUserId uint64) error {
	if guildId == 0 || cmsUserId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(liveentity.TbLiveGuildVisibility)).
		Where(string(liveentity.LiveGuildVisibilityGuildId), guildId).
		Where(string(liveentity.LiveGuildVisibilityCmsUserId), cmsUserId).
		Delete()
	return err
}

// ListVisibilitiesByGuildId 按工会查询可见性记录(按创建时间升序)
func ListVisibilitiesByGuildId(guildId uint64) []*liveentity.LiveGuildVisibility {
	if guildId == 0 {
		return nil
	}
	rows := make([]*liveentity.LiveGuildVisibility, 0)
	_ = g.DB().Model(string(liveentity.TbLiveGuildVisibility)).
		Where(string(liveentity.LiveGuildVisibilityGuildId), guildId).
		Order("created_at asc").
		Scan(&rows)
	return rows
}

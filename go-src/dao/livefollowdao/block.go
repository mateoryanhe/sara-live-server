package livefollowdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

// IsBlocked 查询 userId 是否已拉黑 targetId
func IsBlocked(userId, targetId uint64) bool {
	if userId == 0 || targetId == 0 || userId == targetId {
		return false
	}
	existing := GetByUserAnchor(userId, targetId)
	return existing != nil && existing.Status == entity.LiveFollowStatusBlock
}

// CountBlockedByUser 统计用户当前拉黑数量
func CountBlockedByUser(userId uint64) int {
	if userId == 0 {
		return 0
	}
	n, _ := g.Model(string(entity.TbLiveFollow)).
		Where("user_id = ? AND status = ?", userId, entity.LiveFollowStatusBlock).
		Count()
	return n
}

// GetBlockedListByUser 分页获取用户当前拉黑列表(仅 Status == Block)
func GetBlockedListByUser(userId uint64, page, pageSize int) (int, []*entity.LiveFollow) {
	list := make([]*entity.LiveFollow, 0)
	if userId == 0 {
		return 0, list
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	m := g.Model(string(entity.TbLiveFollow)).
		Where("user_id = ? AND status = ?", userId, entity.LiveFollowStatusBlock)
	total, _ := m.Clone().Count()
	_ = m.Clone().Order("updated_at desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return total, list
}

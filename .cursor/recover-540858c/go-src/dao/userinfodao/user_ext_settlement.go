package userinfodao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	userentity "xr-game-server/entity/user"
)

// ListUserExtWithShortVideoUnsettledIncome 查询有待结算短视频收益的用户(缓存+DB合并)
func ListUserExtWithShortVideoUnsettledIncome() []*userentity.UserExt {
	seen := make(map[uint64]struct{})
	list := make([]*userentity.UserExt, 0)

	if userExtCacheMgr != nil {
		if vals, err := userExtCacheMgr.Values(gctx.New()); err == nil {
			for _, ext := range vals {
				if ext == nil || ext.ID == 0 || ext.ShortVideoUnsettledIncome <= 0 {
					continue
				}
				if _, ok := seen[ext.ID]; ok {
					continue
				}
				seen[ext.ID] = struct{}{}
				list = append(list, ext)
			}
		}
	}

	rows := make([]*userentity.UserExt, 0)
	_ = g.Model(string(userentity.TbUserExt)).Unscoped().
		Where(string(userentity.UserExtShortVideoUnsettledIncome)+" > ?", 0).
		Scan(&rows)
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		if _, ok := seen[row.ID]; ok {
			continue
		}
		if row.ShortVideoUnsettledIncome <= 0 {
			continue
		}
		seen[row.ID] = struct{}{}
		list = append(list, row)
	}
	return list
}

// CollectUserExtWithShortVideoUnsettledIncomeUserIds 收集待结算用户ID
func CollectUserExtWithShortVideoUnsettledIncomeUserIds() []uint64 {
	rows := ListUserExtWithShortVideoUnsettledIncome()
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.ID > 0 {
			ids = append(ids, row.ID)
		}
	}
	return ids
}

package auth

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/user"
)

var userMaxId *entity.UserMaxId

const userMaxIdRecordID uint64 = 1

func initUserMaxId() {
	row := &entity.UserMaxId{}
	err := g.Model(string(entity.TbUserMaxId)).Order("id asc").Limit(1).Scan(row)
	if err != nil || row.ID == 0 {
		g.Log().Warningf(gctx.New(), "user_max_id load empty, create start=%d err=%v", entity.UserMaxIdDefaultID, err)
		userMaxId = entity.NewUserMaxId(userMaxIdRecordID)
		// 初始最大ID=1000000, 第一次 NextUserId 得到 1000001
		userMaxId.SetMaxId(entity.UserMaxIdDefaultID)
		return
	}
	userMaxId = row
	// 兼容错误初始化(曾写成0导致发号从1开始): 拉回到默认起点
	if userMaxId.Add(0) < entity.UserMaxIdDefaultID {
		g.Log().Warningf(gctx.New(), "user_max_id too small=%d, reset to %d", userMaxId.Add(0), entity.UserMaxIdDefaultID)
		userMaxId.SetMaxId(entity.UserMaxIdDefaultID)
	}
}

// NextUserId 原子分配下一个用户ID
func NextUserId() uint64 {
	if userMaxId == nil {
		initUserMaxId()
	}
	return userMaxId.Add(1)
}

// CurrentUserMaxId 当前已分配的最大用户ID(不自增)
func CurrentUserMaxId() uint64 {
	if userMaxId == nil {
		return 0
	}
	return userMaxId.Add(0)
}

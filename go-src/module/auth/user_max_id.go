package auth

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/user"
)

var userMaxId *entity.UserMaxId

func initUserMaxId() {
	row := &entity.UserMaxId{}
	err := g.Model(string(entity.TbUserMaxId)).
		WherePri(entity.UserMaxIdDefaultID).
		Scan(row)
	if err != nil || row.ID == 0 {
		g.Log().Warningf(gctx.New(), "user_max_id load empty, create default id=%d err=%v", entity.UserMaxIdDefaultID, err)
		userMaxId = entity.NewUserMaxId(entity.UserMaxIdDefaultID)
		userMaxId.SetMaxId(0)
		return
	}
	userMaxId = row
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

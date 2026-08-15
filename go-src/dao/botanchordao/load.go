package botanchordao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/entity/user"
)

// LoadAllBotAnchorIds 从数据库加载全部机器人主播用户ID(启动初始化用)
func LoadAllBotAnchorIds() []uint64 {
	type idRow struct {
		ID uint64 `json:"id"`
	}
	rows := make([]*idRow, 0)
	ctx := gctx.New()
	_ = g.Model(string(entity.TbUserInfo)).Ctx(ctx).Unscoped().
		Fields(string(db.IdName)).
		Where(string(entity.UserInfoUserType), entity.UserTypeBotAnchor).
		OrderDesc(string(db.IdName)).
		Scan(&rows)
	ret := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		ret = append(ret, row.ID)
	}
	return ret
}

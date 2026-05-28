package weeklyloginstatdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

// ListRecent 查询最近N周登录统计(按时间正序)
func ListRecent(limit int) []*entity.WeeklyLoginStat {
	list := make([]*entity.WeeklyLoginStat, 0)
	if limit <= 0 {
		limit = 12
	}
	ctx := gctx.New()
	_ = g.Model(string(entity.TbWeeklyLoginStat)).Ctx(ctx).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseWeeklyLoginStats(list)
	return list
}

func reverseWeeklyLoginStats(list []*entity.WeeklyLoginStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

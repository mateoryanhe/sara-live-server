package dailyloginstatdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

// ListRecent 查询最近N天登录统计(按时间正序)
func ListRecent(limit int) []*entity.DailyLoginStat {
	list := make([]*entity.DailyLoginStat, 0)
	if limit <= 0 {
		limit = 30
	}
	ctx := gctx.New()
	_ = g.Model(string(entity.TbDailyLoginStat)).Ctx(ctx).
		Order("id desc").
		Limit(limit).
		Scan(&list)
	reverseDailyLoginStats(list)
	return list
}

func reverseDailyLoginStats(list []*entity.DailyLoginStat) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

package shortvideodao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/shortvideo"
)

func PublishShortVideoStat(data *entity.ShortVideoStat) {
	if data == nil || data.ID == 0 || shortVideoStatCacheMgr == nil {
		return
	}
	shortVideoStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func PublishShortVideoAuthorStat(data *entity.ShortVideoAuthorStat) {
	if data == nil || data.ID == 0 || shortVideoAuthorStatCacheMgr == nil {
		return
	}
	shortVideoAuthorStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func PublishWatchMap(userId uint64, data *userWatchMap) {
	if userId == 0 || data == nil || watchCacheMgr == nil {
		return
	}
	watchCacheMgr.PublishRow(gctx.New(), userId, data)
}

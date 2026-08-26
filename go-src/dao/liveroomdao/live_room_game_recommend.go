package liveroomdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

var liveRoomGameRecommendCacheMgr *cache.ListCache[*entity.LiveRoomGameRecommend]

func InitLiveRoomGameRecommendDao() {
	liveRoomGameRecommendCacheMgr = cache.NewPermanentListCache[*entity.LiveRoomGameRecommend]()
}

func loadActiveLiveRoomGameRecommendsByRoomIDFromDB(liveRoomID uint64) []*entity.LiveRoomGameRecommend {
	rows := make([]*entity.LiveRoomGameRecommend, 0)
	if liveRoomID == 0 {
		return rows
	}
	_ = g.DB().Model(string(entity.TbLiveRoomGameRecommend)).
		Where(string(entity.LiveRoomGameRecommendStatus)+" = ?", entity.LiveRoomGameRecommendStatusActive).
		Where(string(entity.LiveRoomGameRecommendLiveRoomId)+" = ?", liveRoomID).
		Order(string(db.IdName) + " asc").
		Scan(&rows)
	return rows
}

func GetActiveLiveRoomGameRecommendsByRoomID(liveRoomID uint64) []*entity.LiveRoomGameRecommend {
	if liveRoomID == 0 {
		return make([]*entity.LiveRoomGameRecommend, 0)
	}
	if liveRoomGameRecommendCacheMgr == nil {
		return cloneLiveRoomGameRecommendList(loadActiveLiveRoomGameRecommendsByRoomIDFromDB(liveRoomID))
	}
	v := liveRoomGameRecommendCacheMgr.MustGetList(gctx.New(), liveRoomID, func(ctx context.Context) ([]*entity.LiveRoomGameRecommend, error) {
		return loadActiveLiveRoomGameRecommendsByRoomIDFromDB(liveRoomID), nil
	})
	if v == nil {
		return make([]*entity.LiveRoomGameRecommend, 0)
	}
	return cloneLiveRoomGameRecommendList(v)
}

func GetLiveRoomGameRecommendByID(id uint64) *entity.LiveRoomGameRecommend {
	if id == 0 {
		return nil
	}
	var row entity.LiveRoomGameRecommend
	if err := g.DB().Model(string(entity.TbLiveRoomGameRecommend)).Where(string(db.IdName)+" = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func ExistsActiveLiveRoomGameRecommend(liveRoomID uint64, gameCode string) bool {
	gameCode = strings.TrimSpace(gameCode)
	if liveRoomID == 0 || gameCode == "" {
		return false
	}
	for _, row := range GetActiveLiveRoomGameRecommendsByRoomID(liveRoomID) {
		if row != nil && row.GameCode == gameCode {
			return true
		}
	}
	return false
}

func CreateLiveRoomGameRecommend(row *entity.LiveRoomGameRecommend) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbLiveRoomGameRecommend)).Save(row)
	if err != nil {
		return err
	}
	refreshLiveRoomGameRecommendCache(row.LiveRoomId)
	return nil
}

func MarkLiveRoomGameRecommendDeleted(id uint64) error {
	row := GetLiveRoomGameRecommendByID(id)
	if row == nil {
		return nil
	}
	row.SetStatus(entity.LiveRoomGameRecommendStatusDeleted)
	refreshLiveRoomGameRecommendCache(row.LiveRoomId)
	return nil
}

func refreshLiveRoomGameRecommendCache(liveRoomID uint64) {
	if liveRoomGameRecommendCacheMgr == nil || liveRoomID == 0 {
		return
	}
	liveRoomGameRecommendCacheMgr.PublishList(gctx.New(), liveRoomID, loadActiveLiveRoomGameRecommendsByRoomIDFromDB(liveRoomID))
}

// PreloadLiveRoomGameRecommendToCache 启动时批量预热全部有效推荐游戏(按直播间分组,永不过期).
func PreloadLiveRoomGameRecommendToCache() int {
	if liveRoomGameRecommendCacheMgr == nil {
		return 0
	}
	rows := loadAllActiveLiveRoomGameRecommendsFromDB()
	if len(rows) == 0 {
		return 0
	}
	grouped := make(map[uint64][]*entity.LiveRoomGameRecommend)
	for _, row := range rows {
		if row == nil || row.LiveRoomId == 0 {
			continue
		}
		grouped[row.LiveRoomId] = append(grouped[row.LiveRoomId], row)
	}
	for liveRoomID, list := range grouped {
		liveRoomGameRecommendCacheMgr.PublishList(gctx.New(), liveRoomID, cloneLiveRoomGameRecommendList(list))
	}
	return len(grouped)
}

func loadAllActiveLiveRoomGameRecommendsFromDB() []*entity.LiveRoomGameRecommend {
	rows := make([]*entity.LiveRoomGameRecommend, 0)
	_ = g.DB().Model(string(entity.TbLiveRoomGameRecommend)).
		Where(string(entity.LiveRoomGameRecommendStatus)+" = ?", entity.LiveRoomGameRecommendStatusActive).
		Order(string(entity.LiveRoomGameRecommendLiveRoomId) + " asc, " + string(db.IdName) + " asc").
		Scan(&rows)
	return rows
}

func cloneLiveRoomGameRecommend(row *entity.LiveRoomGameRecommend) *entity.LiveRoomGameRecommend {
	if row == nil {
		return nil
	}
	return &entity.LiveRoomGameRecommend{
		OneModel:   row.OneModel,
		LiveRoomId: row.LiveRoomId,
		GameCode:   row.GameCode,
		Status:     row.Status,
	}
}

func cloneLiveRoomGameRecommendList(list []*entity.LiveRoomGameRecommend) []*entity.LiveRoomGameRecommend {
	if len(list) == 0 {
		return make([]*entity.LiveRoomGameRecommend, 0)
	}
	out := make([]*entity.LiveRoomGameRecommend, 0, len(list))
	for _, row := range list {
		out = append(out, cloneLiveRoomGameRecommend(row))
	}
	return out
}

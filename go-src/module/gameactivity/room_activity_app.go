package gameactivity

import (
	"context"
	"strings"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/activitydto"
	"xr-game-server/errercode"
	"xr-game-server/module/game"
)

// GetRoomActivity App: 上报 roomId,返回活动按钮 icon、列表 URL、以及推荐 itemId.
// itemId: 房间推荐列表第一条(须在上架列表中); 无则取上架列表第一条(值为第三方编码).
func GetRoomActivity(_ context.Context, req *activitydto.GetRoomActivityReq) (*activitydto.GetRoomActivityRes, error) {
	if req == nil || req.RoomId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	res := &activitydto.GetRoomActivityRes{}
	// App 热路径只读永久缓存(CMS 保存时 Reload); 不打库.
	if cfg := cfgdao.GetGameActivityCfgFromMemory(); cfg != nil {
		res.ButtonIcon = strings.TrimSpace(cfg.ButtonIcon)
		res.ListUrl = strings.TrimSpace(cfg.GameListUrl)
	}
	res.ItemId = resolveRoomActivityItemId(req.RoomId)
	return res, nil
}

func resolveRoomActivityItemId(roomId uint64) string {
	onShelf := cfgdao.GetGameCfgCodeSetFromMemory()
	for _, row := range liveroomdao.GetActiveLiveRoomGameRecommendsByRoomID(roomId) {
		if row == nil {
			continue
		}
		code := strings.TrimSpace(row.GameCode)
		if code == "" {
			continue
		}
		if _, ok := onShelf[code]; ok {
			return code
		}
	}
	for _, row := range game.GetAllOnShelfGamesFromMemory() {
		if row == nil {
			continue
		}
		code := strings.TrimSpace(row.GameCode)
		if code != "" {
			return code
		}
	}
	return ""
}

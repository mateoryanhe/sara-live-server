package liveroom

import (
	"context"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity/game"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/game"
)

const liveRoomGameRecommendMaxCount = 20

// GetLiveRoomGameRecommendList App 查询直播间推荐游戏列表(读永久缓存 + game_cfgs).
func GetLiveRoomGameRecommendList(_ context.Context, req *liveroomdto.GetLiveRoomGameRecommendListReq) (*liveroomdto.GetLiveRoomGameRecommendListRes, error) {
	if req == nil || req.LiveRoomId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rows := liveroomdao.GetActiveLiveRoomGameRecommendsByRoomID(req.LiveRoomId)
	gameMap := buildOnShelfGameMap()
	list := make([]*liveroomdto.LiveRoomGameRecommendItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.GameCode == "" {
			continue
		}
		gameCfg, ok := gameMap[row.GameCode]
		if !ok {
			continue
		}
		list = append(list, toLiveRoomGameRecommendItem(row, gameCfg))
	}
	return &liveroomdto.GetLiveRoomGameRecommendListRes{List: list}, nil
}

// AddLiveRoomGameRecommend 新增直播间推荐游戏(写库并更新缓存列表).
func AddLiveRoomGameRecommend(liveRoomID uint64, gameCode string) (*liveentity.LiveRoomGameRecommend, error) {
	if liveRoomID == 0 || gameCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if liveroomdao.ExistsActiveLiveRoomGameRecommend(liveRoomID, gameCode) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	now := time.Now()
	row := &liveentity.LiveRoomGameRecommend{
		LiveRoomId: liveRoomID,
		GameCode:   gameCode,
		Status:     liveentity.LiveRoomGameRecommendStatusActive,
	}
	row.CreatedAt = now
	row.UpdatedAt = now
	if err := liveroomdao.CreateLiveRoomGameRecommend(row); err != nil {
		return nil, err
	}
	return row, nil
}

// DeleteLiveRoomGameRecommend 删除直播间推荐游戏(通过 status 控制,并更新缓存列表).
func DeleteLiveRoomGameRecommend(id uint64) error {
	if id == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return liveroomdao.MarkLiveRoomGameRecommendDeleted(id)
}

func buildOnShelfGameMap() map[string]*entity.GameCfg {
	all := game.GetAllOnShelfGamesFromMemory()
	m := make(map[string]*entity.GameCfg, len(all))
	for _, row := range all {
		if row == nil || row.GameCode == "" {
			continue
		}
		m[row.GameCode] = row
	}
	return m
}

func toLiveRoomGameRecommendItem(row *liveentity.LiveRoomGameRecommend, gameCfg *entity.GameCfg) *liveroomdto.LiveRoomGameRecommendItem {
	item := &liveroomdto.LiveRoomGameRecommendItem{
		GameCode: row.GameCode,
	}
	if gameCfg == nil {
		return item
	}
	item.NameEn = game.ResolveAppGameName(gameCfg)
	item.Cover = game.ResolveAppGameCover(gameCfg)
	return item
}

// SyncLiveRoomGameRecommendList 同步直播间推荐游戏列表(仅游戏直播间使用,按上报 gameCode 全量覆盖).
func SyncLiveRoomGameRecommendList(liveRoomID uint64, gameCodes []string) error {
	if liveRoomID == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	normalized := normalizeLiveRoomGameCodes(gameCodes)
	if len(normalized) == 0 {
		return nil
	}
	if len(normalized) > liveRoomGameRecommendMaxCount {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	for _, code := range normalized {
		if !cfgdao.IsGameOnShelfFromMemory(code) {
			return errercode.CreateCode(errercode.GameCfgNonExist)
		}
	}

	current := liveroomdao.GetActiveLiveRoomGameRecommendsByRoomID(liveRoomID)
	currentMap := make(map[string]*liveentity.LiveRoomGameRecommend, len(current))
	for _, row := range current {
		if row == nil || row.GameCode == "" {
			continue
		}
		currentMap[row.GameCode] = row
	}

	targetSet := make(map[string]struct{}, len(normalized))
	for _, code := range normalized {
		targetSet[code] = struct{}{}
		if _, ok := currentMap[code]; ok {
			continue
		}
		if _, err := AddLiveRoomGameRecommend(liveRoomID, code); err != nil {
			return err
		}
	}
	for _, row := range current {
		if row == nil {
			continue
		}
		if _, ok := targetSet[row.GameCode]; ok {
			continue
		}
		if err := DeleteLiveRoomGameRecommend(row.ID); err != nil {
			return err
		}
	}
	return nil
}

func normalizeLiveRoomGameCodes(gameCodes []string) []string {
	if len(gameCodes) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(gameCodes))
	list := make([]string, 0, len(gameCodes))
	for _, raw := range gameCodes {
		code := strings.TrimSpace(raw)
		if code == "" {
			continue
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		list = append(list, code)
	}
	return list
}

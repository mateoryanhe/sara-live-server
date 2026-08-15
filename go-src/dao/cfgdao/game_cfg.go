package cfgdao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/game"
)

const gameCfgCacheKey = "all"

var gameCfgCacheMgr *cache.CacheMgr

func InitGameCfgDao() {
	gameCfgCacheMgr = cache.NewPermanentCacheMgr()
}

func loadAllGameCfgFromDB() []*entity.GameCfg {
	rows := make([]*entity.GameCfg, 0)
	_ = g.DB().Model(string(entity.TbGameCfg)).Order(string(db.IdName) + " asc").Scan(&rows)
	return rows
}

// ReloadGameCfgCache 上架游戏变更后刷新内存(永不过期).
func ReloadGameCfgCache() {
	if gameCfgCacheMgr == nil {
		return
	}
	gameCfgCacheMgr.FlushCache(gameCfgCacheKey, loadAllGameCfgFromDB())
}

// GetAllGameCfgFromMemory 获取全部上架游戏配置(仅内存).
func GetAllGameCfgFromMemory() []*entity.GameCfg {
	if gameCfgCacheMgr == nil {
		return make([]*entity.GameCfg, 0)
	}
	v := gameCfgCacheMgr.GetFromCache(gameCfgCacheKey)
	if v == nil {
		return make([]*entity.GameCfg, 0)
	}
	list, _ := v.([]*entity.GameCfg)
	if list == nil {
		return make([]*entity.GameCfg, 0)
	}
	return append(make([]*entity.GameCfg, 0, len(list)), list...)
}

func GetGameCfgCodeSetFromMemory() map[string]struct{} {
	all := GetAllGameCfgFromMemory()
	set := make(map[string]struct{}, len(all))
	for _, row := range all {
		if row == nil || row.GameCode == "" {
			continue
		}
		set[row.GameCode] = struct{}{}
	}
	return set
}

func IsGameOnShelfFromMemory(gameCode string) bool {
	if gameCode == "" {
		return false
	}
	for _, row := range GetAllGameCfgFromMemory() {
		if row != nil && row.GameCode == gameCode {
			return true
		}
	}
	return false
}

func GetGameCfgByGameCode(gameCode string) *entity.GameCfg {
	gameCode = strings.TrimSpace(gameCode)
	if gameCode == "" {
		return nil
	}
	for _, row := range GetAllGameCfgFromMemory() {
		if row != nil && row.GameCode == gameCode {
			return row
		}
	}
	var row entity.GameCfg
	if err := g.DB().Model(string(entity.TbGameCfg)).Where("game_code = ?", gameCode).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func CreateGameCfg(row *entity.GameCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbGameCfg)).Save(row)
	return err
}

func DeleteGameCfg(id uint64) error {
	if id == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbGameCfg)).WherePri(id).Delete()
	return err
}

func DeleteGameCfgByGameCode(gameCode string) error {
	gameCode = strings.TrimSpace(gameCode)
	if gameCode == "" {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbGameCfg)).Where("game_code = ?", gameCode).Delete()
	return err
}

func DeleteGameCfgByGameCodes(gameCodes []string) (int, error) {
	codes := make([]string, 0, len(gameCodes))
	for _, code := range gameCodes {
		code = strings.TrimSpace(code)
		if code != "" {
			codes = append(codes, code)
		}
	}
	if len(codes) == 0 {
		return 0, nil
	}
	result, err := g.DB().Model(string(entity.TbGameCfg)).Where("game_code IN (?)", codes).Delete()
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

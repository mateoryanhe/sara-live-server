package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const gameShelfCfgCacheKey = "all"

var gameShelfCfgCacheMgr *cache.CacheMgr

func InitGameShelfCfgDao() {
	gameShelfCfgCacheMgr = cache.NewPermanentCacheMgr()
}

func loadAllGameShelfCfgFromDB() []*entity.GameShelfCfg {
	rows := make([]*entity.GameShelfCfg, 0)
	_ = g.DB().Model(string(entity.TbGameShelfCfg)).Order(string(db.IdName) + " asc").Scan(&rows)
	return rows
}

func ReloadGameShelfCfgCache() {
	if gameShelfCfgCacheMgr == nil {
		return
	}
	gameShelfCfgCacheMgr.FlushCache(gameShelfCfgCacheKey, loadAllGameShelfCfgFromDB())
}

func GetAllGameShelfCfgFromMemory() []*entity.GameShelfCfg {
	if gameShelfCfgCacheMgr == nil {
		return make([]*entity.GameShelfCfg, 0)
	}
	v := gameShelfCfgCacheMgr.GetFromCache(gameShelfCfgCacheKey)
	if v == nil {
		return make([]*entity.GameShelfCfg, 0)
	}
	list, _ := v.([]*entity.GameShelfCfg)
	if list == nil {
		return make([]*entity.GameShelfCfg, 0)
	}
	return append(make([]*entity.GameShelfCfg, 0, len(list)), list...)
}

func GetGameShelfCodeSetFromMemory() map[string]struct{} {
	all := GetAllGameShelfCfgFromMemory()
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
	for _, row := range GetAllGameShelfCfgFromMemory() {
		if row != nil && row.GameCode == gameCode {
			return true
		}
	}
	return false
}

func GetGameShelfCfgByCode(gameCode string) *entity.GameShelfCfg {
	if gameCode == "" {
		return nil
	}
	var row entity.GameShelfCfg
	if err := g.DB().Model(string(entity.TbGameShelfCfg)).Where("game_code = ?", gameCode).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func CreateGameShelfCfg(row *entity.GameShelfCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbGameShelfCfg)).Save(row)
	return err
}

func DeleteGameShelfCfg(id uint64) error {
	if id == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbGameShelfCfg)).WherePri(id).Delete()
	return err
}

func DeleteGameShelfCfgByGameCode(gameCode string) error {
	gameCode = strings.TrimSpace(gameCode)
	if gameCode == "" {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbGameShelfCfg)).Where("game_code = ?", gameCode).Delete()
	return err
}

func DeleteGameShelfCfgByGameCodes(gameCodes []string) (int, error) {
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
	result, err := g.DB().Model(string(entity.TbGameShelfCfg)).Where("game_code IN (?)", codes).Delete()
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

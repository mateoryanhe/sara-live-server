package cfgdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/game"
)

const gameCfgCacheKey = "all"

var gameCfgCacheMgr *cache.ListCache[*entity.GameCfg]

func InitGameCfgDao() {
	gameCfgCacheMgr = cache.NewPermanentListCache[*entity.GameCfg]()
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
	gameCfgCacheMgr.PublishList(gctx.New(), gameCfgCacheKey, loadAllGameCfgFromDB())
}

// GetAllGameCfgFromMemory 获取全部上架游戏配置(仅内存).
func GetAllGameCfgFromMemory() []*entity.GameCfg {
	if gameCfgCacheMgr == nil {
		return make([]*entity.GameCfg, 0)
	}
	v, _ := gameCfgCacheMgr.GetListCached(gctx.New(), gameCfgCacheKey)
	if v == nil {
		return make([]*entity.GameCfg, 0)
	}
	return append(make([]*entity.GameCfg, 0, len(v)), v...)
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

// GetGameCfgShelfKeySetFromMemory 已上架游戏键集合(game_code + platform).
func GetGameCfgShelfKeySetFromMemory() map[string]struct{} {
	all := GetAllGameCfgFromMemory()
	set := make(map[string]struct{}, len(all))
	for _, row := range all {
		if row == nil || row.GameCode == "" {
			continue
		}
		set[GameCfgShelfKey(row.GameCode, row.Platform)] = struct{}{}
	}
	return set
}

func GameCfgShelfKey(gameCode, platform string) string {
	return strings.TrimSpace(gameCode) + "\x00" + strings.TrimSpace(platform)
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
	_, err := g.DB().Model(string(entity.TbGameCfg)).Data(g.Map{
		"game_code":       row.GameCode,
		"cover":           row.Cover,
		"name_en":         row.NameEn,
		"platform":        row.Platform,
		"live_game_name":  row.LiveGameName,
		"live_game_cover": row.LiveGameCover,
		"created_at":      row.CreatedAt,
		"updated_at":      row.UpdatedAt,
	}).Insert()
	return err
}

func SetGameCfgPlatform(gameCode, platform string) (bool, error) {
	gameCode = strings.TrimSpace(gameCode)
	platform = strings.TrimSpace(platform)
	if gameCode == "" || platform == "" {
		return false, nil
	}
	result, err := g.DB().Model(string(entity.TbGameCfg)).Where("game_code = ?", gameCode).Data(g.Map{
		"platform": platform,
	}).Update()
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
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

func UpdateGameCfgLiveDisplay(gameCode, liveGameName, liveGameCover string) (bool, error) {
	gameCode = strings.TrimSpace(gameCode)
	if gameCode == "" {
		return false, nil
	}
	result, err := g.DB().Model(string(entity.TbGameCfg)).Where("game_code = ?", gameCode).Data(g.Map{
		"live_game_name":  liveGameName,
		"live_game_cover": liveGameCover,
	}).Update()
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

package cfgdao

import (
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/dao/gamequery"
	entity "xr-game-server/entity/game"
)

// VendorGameLibQuery 游戏库表查询条件(CMS).
type VendorGameLibQuery struct {
	GameCode  string
	GameCodes []string
	Name      string
	Platform  string
	Category  string
	PageIndex int
	PageSize  int
}

// QueryVendorGameLibs 分页查询游戏库表(直接读 DB, 不走缓存).
func QueryVendorGameLibs(q *VendorGameLibQuery) (int, []*entity.VendorGameLib) {
	list := make([]*entity.VendorGameLib, 0)
	if q == nil {
		return 0, list
	}
	pageIndex := q.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := q.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	m := g.DB().Model(string(entity.TbVendorGameLib))
	m = gamequery.ApplyGameCodeFilter(m, "game_code", q.GameCode, q.GameCodes)
	if name := strings.TrimSpace(q.Name); name != "" {
		m = m.Where("(name LIKE ? OR name_en LIKE ?)", "%"+name+"%", "%"+name+"%")
	}
	if platform := strings.TrimSpace(q.Platform); platform != "" {
		m = m.Where("platform = ?", platform)
	}
	if category := strings.TrimSpace(q.Category); category != "" {
		m = m.Where("category LIKE ?", "%"+category+"%")
	}

	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Order(string(db.IdName) + " asc").
		Limit(pageSize).Offset((pageIndex - 1) * pageSize).
		Scan(&list)
	return total, list
}

// ListDistinctVendorGamePlatforms 返回游戏库中去重后的平台编码(升序,库内原文).
func ListDistinctVendorGamePlatforms() []string {
	type platformRow struct {
		Platform string `json:"platform"`
	}
	rows := make([]*platformRow, 0)
	_ = g.DB().Model(string(entity.TbVendorGameLib)).
		Fields("DISTINCT platform").
		Where("platform <> ?", "").
		Order("platform asc").
		Scan(&rows)
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		p := strings.TrimSpace(row.Platform)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

// ListDistinctVendorGamePlatformsUpper 游戏库平台去重并转大写(下注/派彩日志筛选用).
func ListDistinctVendorGamePlatformsUpper() []string {
	raw := ListDistinctVendorGamePlatforms()
	out := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, p := range raw {
		u := strings.ToUpper(strings.TrimSpace(p))
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		out = append(out, u)
	}
	return out
}

// GetVendorGameLib 按 game_code + platform 查游戏库(直接读 DB).
func GetVendorGameLib(gameCode, platform string) *entity.VendorGameLib {
	gameCode = strings.TrimSpace(gameCode)
	platform = strings.TrimSpace(platform)
	if gameCode == "" || platform == "" {
		return nil
	}
	var row entity.VendorGameLib
	if err := g.DB().Model(string(entity.TbVendorGameLib)).
		Where("game_code = ? AND platform = ?", gameCode, platform).
		Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

// ReplaceAllVendorGameLibs 全量覆盖游戏库表.
func ReplaceAllVendorGameLibs(rows []*entity.VendorGameLib) error {
	if _, err := g.DB().Model(string(entity.TbVendorGameLib)).Where("id > ?", 0).Delete(); err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbVendorGameLib)).Data(rows).Batch(len(rows)).Insert()
	return err
}

func BuildVendorGameLibRows(games []*entity.VendorGameLib, now time.Time) []*entity.VendorGameLib {
	if len(games) == 0 {
		return make([]*entity.VendorGameLib, 0)
	}
	out := make([]*entity.VendorGameLib, 0, len(games))
	seen := make(map[string]struct{}, len(games))
	for _, row := range games {
		if row == nil || strings.TrimSpace(row.GameCode) == "" {
			continue
		}
		item := &entity.VendorGameLib{
			GameCode: strings.TrimSpace(row.GameCode),
			Name:     strings.TrimSpace(row.Name),
			NameEn:   strings.TrimSpace(row.NameEn),
			Category: strings.TrimSpace(row.Category),
			Cover:    strings.TrimSpace(row.Cover),
			Platform: strings.TrimSpace(row.Platform),
		}
		key := item.GameCode + "\x00" + item.Platform
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		item.CreatedAt = now
		item.UpdatedAt = now
		out = append(out, item)
	}
	return out
}

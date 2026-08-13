package guilddao

import (
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity"
)

const guildTimeLayout = "2006-01-02 15:04:05"

type guildSnapshot struct {
	byID       map[uint64]*entity.LiveGuild
	byName     map[string]*entity.LiveGuild
	byLeaderId map[uint64][]*entity.LiveGuild
	allList    []*entity.LiveGuild
}

var guildMemoryCache atomic.Value // *guildSnapshot

// ReloadGuildMemory 从 DB 重新加载工会全量快照并替换内存缓存
func ReloadGuildMemory() {
	rows := loadAllGuildsFromDB()

	byID := make(map[uint64]*entity.LiveGuild, len(rows))
	byName := make(map[string]*entity.LiveGuild, len(rows))
	byLeaderId := make(map[uint64][]*entity.LiveGuild)
	allList := make([]*entity.LiveGuild, 0, len(rows))

	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		byID[row.ID] = row
		byName[row.Name] = row
		byLeaderId[row.LeaderId] = append(byLeaderId[row.LeaderId], row)
		allList = append(allList, row)
	}

	for leaderID := range byLeaderId {
		sort.Slice(byLeaderId[leaderID], func(i, j int) bool {
			return byLeaderId[leaderID][i].ID < byLeaderId[leaderID][j].ID
		})
	}

	sort.Slice(allList, func(i, j int) bool {
		return allList[i].CreatedAt.After(allList[j].CreatedAt)
	})

	guildMemoryCache.Store(&guildSnapshot{
		byID:       byID,
		byName:     byName,
		byLeaderId: byLeaderId,
		allList:    allList,
	})
}

func loadAllGuildsFromDB() []*entity.LiveGuild {
	rows := make([]*entity.LiveGuild, 0)
	_ = g.DB().Model(string(entity.TbLiveGuild)).Order("created_at desc").Scan(&rows)
	return rows
}

func getGuildSnapshot() *guildSnapshot {
	v := guildMemoryCache.Load()
	if v == nil {
		return &guildSnapshot{
			byID:       make(map[uint64]*entity.LiveGuild),
			byName:     make(map[string]*entity.LiveGuild),
			byLeaderId: make(map[uint64][]*entity.LiveGuild),
			allList:    make([]*entity.LiveGuild, 0),
		}
	}
	return v.(*guildSnapshot)
}

func getGuildByIdFromMemory(id uint64) *entity.LiveGuild {
	if id == 0 {
		return nil
	}
	return getGuildSnapshot().byID[id]
}

func getGuildByNameFromMemory(name string) *entity.LiveGuild {
	if name == "" {
		return nil
	}
	return getGuildSnapshot().byName[name]
}

func listGuildsByLeaderIdFromMemory(leaderId uint64) []*entity.LiveGuild {
	if leaderId == 0 {
		return nil
	}
	rows := getGuildSnapshot().byLeaderId[leaderId]
	if len(rows) == 0 {
		return nil
	}
	list := make([]*entity.LiveGuild, len(rows))
	copy(list, rows)
	return list
}

func formatGuildTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(guildTimeLayout)
}

func toGuildListRes(row *entity.LiveGuild) *guilddto.GuildListRes {
	if row == nil {
		return nil
	}
	return &guilddto.GuildListRes{
		ID:          strconv.FormatUint(row.ID, 10),
		Name:        row.Name,
		LeaderId:    strconv.FormatUint(row.LeaderId, 10),
		LeaderName:  row.LeaderName,
		Contact:     row.Contact,
		Description: row.Description,
		Status:      row.Status,
		CreatedAt:   formatGuildTime(row.CreatedAt),
		UpdatedAt:   formatGuildTime(row.UpdatedAt),
	}
}

func filterGuildListFromMemory(nameKeyword string) []*entity.LiveGuild {
	keyword := strings.ToLower(strings.TrimSpace(nameKeyword))
	rows := getGuildSnapshot().allList
	if keyword == "" {
		return rows
	}
	filtered := make([]*entity.LiveGuild, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		if strings.Contains(strings.ToLower(row.Name), keyword) {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func paginateGuildList(rows []*entity.LiveGuild, pageIndex, pageSize int) ([]*entity.LiveGuild, int) {
	total := len(rows)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	start := (pageIndex - 1) * pageSize
	if start >= total {
		return []*entity.LiveGuild{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return rows[start:end], total
}

func queryGuildListFromMemory(req *guilddto.GuildListReq) (int, []*guilddto.GuildListRes) {
	if req == nil {
		return 0, []*guilddto.GuildListRes{}
	}
	rows, total := paginateGuildList(filterGuildListFromMemory(req.Name), req.PageIndex, req.PageSize)
	list := make([]*guilddto.GuildListRes, 0, len(rows))
	for _, row := range rows {
		list = append(list, toGuildListRes(row))
	}
	return total, list
}

// AddGuildToMemory 将新工会写入内存快照(异步入库由 syndb 负责)
func AddGuildToMemory(g *entity.LiveGuild) {
	if g == nil || g.ID == 0 {
		return
	}
	replaceGuildSnapshot(func(s *guildSnapshot) {
		s.byID[g.ID] = g
		s.byName[g.Name] = g
		s.byLeaderId[g.LeaderId] = appendGuildByLeaderID(s.byLeaderId[g.LeaderId], g)
		s.allList = insertGuildSortedByCreatedAt(s.allList, g)
	})
}

// RemoveGuildFromMemory 从内存快照移除工会
func RemoveGuildFromMemory(id uint64) {
	if id == 0 {
		return
	}
	replaceGuildSnapshot(func(s *guildSnapshot) {
		g, ok := s.byID[id]
		if !ok || g == nil {
			return
		}
		delete(s.byID, id)
		if g.Name != "" {
			delete(s.byName, g.Name)
		}
		s.byLeaderId[g.LeaderId] = removeGuildFromSlice(s.byLeaderId[g.LeaderId], id)
		s.allList = removeGuildFromSlice(s.allList, id)
	})
}

// ReindexGuildInMemory 工会名称或会长变更后刷新索引
func ReindexGuildInMemory(g *entity.LiveGuild, oldName string, oldLeaderId uint64) {
	if g == nil || g.ID == 0 {
		return
	}
	if oldName == g.Name && oldLeaderId == g.LeaderId {
		return
	}
	replaceGuildSnapshot(func(s *guildSnapshot) {
		if oldName != "" && oldName != g.Name {
			delete(s.byName, oldName)
		}
		s.byName[g.Name] = g
		if oldLeaderId != g.LeaderId {
			s.byLeaderId[oldLeaderId] = removeGuildFromSlice(s.byLeaderId[oldLeaderId], g.ID)
			s.byLeaderId[g.LeaderId] = appendGuildByLeaderID(s.byLeaderId[g.LeaderId], g)
		}
	})
}

func replaceGuildSnapshot(mutator func(*guildSnapshot)) {
	next := cloneGuildSnapshot(getGuildSnapshot())
	mutator(next)
	guildMemoryCache.Store(next)
}

func cloneGuildSnapshot(src *guildSnapshot) *guildSnapshot {
	if src == nil {
		return &guildSnapshot{
			byID:       make(map[uint64]*entity.LiveGuild),
			byName:     make(map[string]*entity.LiveGuild),
			byLeaderId: make(map[uint64][]*entity.LiveGuild),
			allList:    make([]*entity.LiveGuild, 0),
		}
	}
	next := &guildSnapshot{
		byID:       make(map[uint64]*entity.LiveGuild, len(src.byID)),
		byName:     make(map[string]*entity.LiveGuild, len(src.byName)),
		byLeaderId: make(map[uint64][]*entity.LiveGuild, len(src.byLeaderId)),
		allList:    make([]*entity.LiveGuild, len(src.allList)),
	}
	for id, row := range src.byID {
		next.byID[id] = row
	}
	for name, row := range src.byName {
		next.byName[name] = row
	}
	for leaderID, rows := range src.byLeaderId {
		next.byLeaderId[leaderID] = append([]*entity.LiveGuild(nil), rows...)
	}
	copy(next.allList, src.allList)
	return next
}

func appendGuildByLeaderID(rows []*entity.LiveGuild, g *entity.LiveGuild) []*entity.LiveGuild {
	rows = append(rows, g)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].ID < rows[j].ID
	})
	return rows
}

func insertGuildSortedByCreatedAt(rows []*entity.LiveGuild, g *entity.LiveGuild) []*entity.LiveGuild {
	rows = append(rows, g)
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].CreatedAt.After(rows[j].CreatedAt)
	})
	return rows
}

func removeGuildFromSlice(rows []*entity.LiveGuild, id uint64) []*entity.LiveGuild {
	if len(rows) == 0 {
		return rows
	}
	next := make([]*entity.LiveGuild, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == id {
			continue
		}
		next = append(next, row)
	}
	return next
}

package guilddao

import (
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dto/guilddto"
	liveentity "xr-game-server/entity/live"
)

const guildTimeLayout = "2006-01-02 15:04:05"

// GetGuildByIdFromDB 根据 ID 直查数据库(含下架)
func GetGuildByIdFromDB(id uint64) *liveentity.LiveGuild {
	if id == 0 {
		return nil
	}
	var row liveentity.LiveGuild
	err := g.DB().Model(string(liveentity.TbLiveGuild)).
		WherePri(id).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetGuildById 根据 ID 从数据库获取工会(不含已下架)
func GetGuildById(id uint64) *liveentity.LiveGuild {
	if id == 0 {
		return nil
	}
	var row liveentity.LiveGuild
	err := g.DB().Model(string(liveentity.TbLiveGuild)).
		WherePri(id).
		Where(string(liveentity.LiveGuildStatus), liveentity.LiveGuildStatusOnShelf).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetNameMapByIds 批量查询工会名称(CMS列表等场景使用)
func GetNameMapByIds(ids []uint64) map[uint64]string {
	if len(ids) == 0 {
		return nil
	}
	type row struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}
	rows := make([]row, 0)
	_ = g.DB().Model(string(liveentity.TbLiveGuild)).
		Where("id IN (?)", ids).
		Scan(&rows)
	ret := make(map[uint64]string, len(rows))
	for _, r := range rows {
		ret[r.ID] = r.Name
	}
	return ret
}

// ListOnShelfGuilds 查询全部上架工会
func ListOnShelfGuilds() []*liveentity.LiveGuild {
	rows := make([]*liveentity.LiveGuild, 0)
	_ = g.DB().Model(string(liveentity.TbLiveGuild)).
		Where(string(liveentity.LiveGuildStatus), liveentity.LiveGuildStatusOnShelf).
		Scan(&rows)
	return rows
}

// GetGuildByName 根据名称从数据库获取工会(不含已软删除)
func GetGuildByName(name string) *liveentity.LiveGuild {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	var row liveentity.LiveGuild
	err := g.DB().Model(string(liveentity.TbLiveGuild)).
		Where(string(liveentity.LiveGuildName), name).
		Where(string(liveentity.LiveGuildStatus), liveentity.LiveGuildStatusOnShelf).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// ListGuildsByLeaderId 根据会长 CMS 用户 ID 从数据库获取其管理的全部工会(不含已软删除)
func ListGuildsByLeaderId(leaderId uint64) []*liveentity.LiveGuild {
	if leaderId == 0 {
		return nil
	}
	rows := make([]*liveentity.LiveGuild, 0)
	_ = g.DB().Model(string(liveentity.TbLiveGuild)).
		Where(string(liveentity.LiveGuildLeaderId), leaderId).
		Where(string(liveentity.LiveGuildStatus), liveentity.LiveGuildStatusOnShelf).
		Order("id asc").
		Scan(&rows)
	return rows
}

// GetGuildByLeaderId 根据会长 CMS 用户 ID 获取工会(兼容旧逻辑,取第一条)
func GetGuildByLeaderId(leaderId uint64) *liveentity.LiveGuild {
	rows := ListGuildsByLeaderId(leaderId)
	if len(rows) == 0 {
		return nil
	}
	return rows[0]
}

// CreateGuild 直接写库
func CreateGuild(guild *liveentity.LiveGuild) error {
	if guild == nil || guild.ID == 0 {
		return nil
	}
	if guild.Status == 0 {
		guild.Status = liveentity.LiveGuildStatusOnShelf
	}
	now := time.Now()
	if guild.CreatedAt.IsZero() {
		guild.CreatedAt = now
	}
	guild.UpdatedAt = now
	_, err := g.DB().Model(string(liveentity.TbLiveGuild)).Save(guild)
	return err
}

// UpdateGuild 直接写库
func UpdateGuild(guild *liveentity.LiveGuild) error {
	if guild == nil || guild.ID == 0 {
		return nil
	}
	guild.UpdatedAt = time.Now()
	_, err := g.DB().Model(string(liveentity.TbLiveGuild)).Save(guild)
	return err
}

// DeleteGuild 下架工会(status=0)
func DeleteGuild(id uint64) error {
	if id == 0 {
		return nil
	}
	now := time.Now()
	_, err := g.DB().Model(string(liveentity.TbLiveGuild)).
		Data(g.Map{
			string(liveentity.LiveGuildStatus): liveentity.LiveGuildStatusOffShelf,
			"updated_at":                       now,
		}).
		WherePri(id).
		Where(string(liveentity.LiveGuildStatus), liveentity.LiveGuildStatusOnShelf).
		Update()
	return err
}

// OnShelfGuild 上架工会(status=1)
func OnShelfGuild(id uint64) error {
	if id == 0 {
		return nil
	}
	now := time.Now()
	_, err := g.DB().Model(string(liveentity.TbLiveGuild)).
		Data(g.Map{
			string(liveentity.LiveGuildStatus): liveentity.LiveGuildStatusOnShelf,
			"updated_at":                       now,
		}).
		WherePri(id).
		Where(string(liveentity.LiveGuildStatus), liveentity.LiveGuildStatusOffShelf).
		Update()
	return err
}

// GetOffShelfGuildList 分页查询已下架工会(直连数据库)
func GetOffShelfGuildList(req *guilddto.OffShelfGuildListReq) (int, []*guilddto.GuildListRes) {
	if req == nil {
		return 0, []*guilddto.GuildListRes{}
	}
	pageIndex := req.PageIndex
	pageSize := req.PageSize
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	m := g.DB().Model(string(liveentity.TbLiveGuild)).
		Where(string(liveentity.LiveGuildStatus), liveentity.LiveGuildStatusOffShelf)
	if keyword := strings.TrimSpace(req.Name); keyword != "" {
		m = m.WhereLike(string(liveentity.LiveGuildName), "%"+keyword+"%")
	}

	total, err := m.Count()
	if err != nil {
		return 0, []*guilddto.GuildListRes{}
	}

	rows := make([]*liveentity.LiveGuild, 0)
	_ = m.Order("updated_at desc").
		Page(pageIndex, pageSize).
		Scan(&rows)

	list := make([]*guilddto.GuildListRes, 0, len(rows))
	for _, row := range rows {
		if item := toGuildListRes(row); item != nil {
			list = append(list, item)
		}
	}
	return total, list
}

// GetGuildList 从数据库分页查询工会列表(不含已下架)
func GetGuildList(req *guilddto.GuildListReq) (int, []*guilddto.GuildListRes) {
	if req == nil {
		return 0, []*guilddto.GuildListRes{}
	}
	pageIndex := req.PageIndex
	pageSize := req.PageSize
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	m := g.DB().Model(string(liveentity.TbLiveGuild)).
		Where(string(liveentity.LiveGuildStatus), liveentity.LiveGuildStatusOnShelf)
	if keyword := strings.TrimSpace(req.Name); keyword != "" {
		m = m.WhereLike(string(liveentity.LiveGuildName), "%"+keyword+"%")
	}

	total, err := m.Count()
	if err != nil {
		return 0, []*guilddto.GuildListRes{}
	}

	rows := make([]*liveentity.LiveGuild, 0)
	_ = m.Order("created_at desc").
		Page(pageIndex, pageSize).
		Scan(&rows)

	list := make([]*guilddto.GuildListRes, 0, len(rows))
	for _, row := range rows {
		if item := toGuildListRes(row); item != nil {
			list = append(list, item)
		}
	}
	return total, list
}

func toGuildListRes(row *liveentity.LiveGuild) *guilddto.GuildListRes {
	if row == nil {
		return nil
	}
	return &guilddto.GuildListRes{
		ID:          strconv.FormatUint(row.ID, 10),
		Name:        row.Name,
		LeaderId:    strconv.FormatUint(row.LeaderId, 10),
		LeaderName:  row.LeaderName,
		Description: row.Description,
		Status:      row.Status,
		CreatedAt:   formatGuildTime(row.CreatedAt),
		UpdatedAt:   formatGuildTime(row.UpdatedAt),
	}
}

func formatGuildTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(guildTimeLayout)
}

package guilddao

import (
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity"
)

const guildTimeLayout = "2006-01-02 15:04:05"

// GetGuildById 根据 ID 从数据库获取工会(不含已软删除)
func GetGuildById(id uint64) *entity.LiveGuild {
	if id == 0 {
		return nil
	}
	var row entity.LiveGuild
	err := g.DB().Model(string(entity.TbLiveGuild)).
		WherePri(id).
		Where(string(entity.LiveGuildStatus), entity.LiveGuildStatusNormal).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetGuildByName 根据名称从数据库获取工会(不含已软删除)
func GetGuildByName(name string) *entity.LiveGuild {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	var row entity.LiveGuild
	err := g.DB().Model(string(entity.TbLiveGuild)).
		Where(string(entity.LiveGuildName), name).
		Where(string(entity.LiveGuildStatus), entity.LiveGuildStatusNormal).
		Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// ListGuildsByLeaderId 根据会长 CMS 用户 ID 从数据库获取其管理的全部工会(不含已软删除)
func ListGuildsByLeaderId(leaderId uint64) []*entity.LiveGuild {
	if leaderId == 0 {
		return nil
	}
	rows := make([]*entity.LiveGuild, 0)
	_ = g.DB().Model(string(entity.TbLiveGuild)).
		Where(string(entity.LiveGuildLeaderId), leaderId).
		Where(string(entity.LiveGuildStatus), entity.LiveGuildStatusNormal).
		Order("id asc").
		Scan(&rows)
	return rows
}

// GetGuildByLeaderId 根据会长 CMS 用户 ID 获取工会(兼容旧逻辑,取第一条)
func GetGuildByLeaderId(leaderId uint64) *entity.LiveGuild {
	rows := ListGuildsByLeaderId(leaderId)
	if len(rows) == 0 {
		return nil
	}
	return rows[0]
}

// CreateGuild 直接写库
func CreateGuild(guild *entity.LiveGuild) error {
	if guild == nil || guild.ID == 0 {
		return nil
	}
	if guild.Status == 0 {
		guild.Status = entity.LiveGuildStatusNormal
	}
	now := time.Now()
	if guild.CreatedAt.IsZero() {
		guild.CreatedAt = now
	}
	guild.UpdatedAt = now
	_, err := g.DB().Model(string(entity.TbLiveGuild)).Save(guild)
	return err
}

// UpdateGuild 直接写库
func UpdateGuild(guild *entity.LiveGuild) error {
	if guild == nil || guild.ID == 0 {
		return nil
	}
	guild.UpdatedAt = time.Now()
	_, err := g.DB().Model(string(entity.TbLiveGuild)).Save(guild)
	return err
}

// DeleteGuild 软删除工会(status=0)
func DeleteGuild(id uint64) error {
	if id == 0 {
		return nil
	}
	now := time.Now()
	_, err := g.DB().Model(string(entity.TbLiveGuild)).
		Data(g.Map{
			string(entity.LiveGuildStatus): entity.LiveGuildStatusDeleted,
			"updated_at":                   now,
		}).
		WherePri(id).
		Where(string(entity.LiveGuildStatus), entity.LiveGuildStatusNormal).
		Update()
	return err
}

// GetGuildList 从数据库分页查询工会列表(不含已软删除)
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

	m := g.DB().Model(string(entity.TbLiveGuild)).
		Where(string(entity.LiveGuildStatus), entity.LiveGuildStatusNormal)
	if keyword := strings.TrimSpace(req.Name); keyword != "" {
		m = m.WhereLike(string(entity.LiveGuildName), "%"+keyword+"%")
	}

	total, err := m.Count()
	if err != nil {
		return 0, []*guilddto.GuildListRes{}
	}

	rows := make([]*entity.LiveGuild, 0)
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

func toGuildListRes(row *entity.LiveGuild) *guilddto.GuildListRes {
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

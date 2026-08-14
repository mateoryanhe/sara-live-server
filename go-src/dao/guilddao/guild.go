package guilddao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity"
)

// GetGuildById 根据 ID 从内存快照获取工会(不含已软删除)
func GetGuildById(id uint64) *entity.LiveGuild {
	return getGuildByIdFromMemory(id)
}

// GetGuildByName 根据名称从内存快照获取工会(不含已软删除)
func GetGuildByName(name string) *entity.LiveGuild {
	return getGuildByNameFromMemory(name)
}

// ListGuildsByLeaderId 根据会长 CMS 用户 ID 从内存快照获取其管理的全部工会(不含已软删除)
func ListGuildsByLeaderId(leaderId uint64) []*entity.LiveGuild {
	return listGuildsByLeaderIdFromMemory(leaderId)
}

// GetGuildByLeaderId 根据会长 CMS 用户 ID 获取工会(兼容旧逻辑,取第一条)
func GetGuildByLeaderId(leaderId uint64) *entity.LiveGuild {
	rows := ListGuildsByLeaderId(leaderId)
	if len(rows) == 0 {
		return nil
	}
	return rows[0]
}

// CreateGuild 直接写库并加入内存快照
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
	if _, err := g.DB().Model(string(entity.TbLiveGuild)).Save(guild); err != nil {
		return err
	}
	AddGuildToMemory(guild)
	return nil
}

// UpdateGuild 直接写库并刷新内存索引
func UpdateGuild(guild *entity.LiveGuild, oldName string, oldLeaderId uint64) error {
	if guild == nil || guild.ID == 0 {
		return nil
	}
	guild.UpdatedAt = time.Now()
	if _, err := g.DB().Model(string(entity.TbLiveGuild)).Save(guild); err != nil {
		return err
	}
	ReindexGuildInMemory(guild, oldName, oldLeaderId)
	return nil
}

// DeleteGuild 软删除工会(status=0),并从内存快照移除
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
	if err != nil {
		return err
	}
	RemoveGuildFromMemory(id)
	return nil
}

// GetGuildList 从内存快照获取工会列表(不含已软删除)
func GetGuildList(req *guilddto.GuildListReq) (int, []*guilddto.GuildListRes) {
	return queryGuildListFromMemory(req)
}

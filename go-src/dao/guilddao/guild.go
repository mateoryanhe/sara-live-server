package guilddao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity"
)

// GetGuildById 根据 ID 从内存快照获取工会
func GetGuildById(id uint64) *entity.LiveGuild {
	return getGuildByIdFromMemory(id)
}

// GetGuildByName 根据名称从内存快照获取工会
func GetGuildByName(name string) *entity.LiveGuild {
	return getGuildByNameFromMemory(name)
}

// ListGuildsByLeaderId 根据会长 CMS 用户 ID 从内存快照获取其管理的全部工会
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

// CreateGuild 将工会加入内存快照(字段异步入库由 syndb 负责)
func CreateGuild(guild *entity.LiveGuild) error {
	AddGuildToMemory(guild)
	return nil
}

// UpdateGuild 刷新工会内存索引(字段异步入库由 syndb 负责)
func UpdateGuild(guild *entity.LiveGuild, oldName string, oldLeaderId uint64) {
	ReindexGuildInMemory(guild, oldName, oldLeaderId)
}

// DeleteGuild 删除工会并刷新内存快照
func DeleteGuild(id uint64) error {
	_, err := g.DB().Model(string(entity.TbLiveGuild)).WherePri(id).Delete()
	if err != nil {
		return err
	}
	RemoveGuildFromMemory(id)
	return nil
}

// GetGuildList 从内存快照获取工会列表
func GetGuildList(req *guilddto.GuildListReq) (int, []*guilddto.GuildListRes) {
	return queryGuildListFromMemory(req)
}

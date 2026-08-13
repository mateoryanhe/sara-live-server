package guilddao

import (
	"xr-game-server/entity"
)

func InitGuildDao() {
	InitGuildMemberDao()
	ReloadGuildMemory()
}

// GetGuildByIdCached 从内存快照读取工会信息(给 App 使用)
func GetGuildByIdCached(id uint64) *entity.LiveGuild {
	return getGuildByIdFromMemory(id)
}

// RemoveGuildCache 刷新工会内存快照(兼容旧调用)
func RemoveGuildCache(_ uint64) {
	ReloadGuildMemory()
}

package liveroomdao

import (
	"github.com/gogf/gf/v2/database/gdb"
	"xr-game-server/entity/live"
)

// CMSAnchorRoomVisibilityScope CMS 列表按 live_rooms 过滤:可见工会名下主播 OR 已授权平台主播(guild_id=0)
type CMSAnchorRoomVisibilityScope struct {
	Apply             bool
	GuildIds          []uint64
	PlatformAnchorIds []uint64
}

// IsEmptyScope Restrict 开启但工会与平台主播均无可见 ID
func (s CMSAnchorRoomVisibilityScope) IsEmptyScope() bool {
	return s.Apply && len(s.GuildIds) == 0 && len(s.PlatformAnchorIds) == 0
}

// ScopeFromAnchorRoomVisibility 由 cmsvis.AnchorRoomVisibility 构造 DAO 过滤条件
func ScopeFromAnchorRoomVisibility(restrict bool, guildIds, platformAnchorIds []uint64) CMSAnchorRoomVisibilityScope {
	if !restrict {
		return CMSAnchorRoomVisibilityScope{}
	}
	return CMSAnchorRoomVisibilityScope{
		Apply:             true,
		GuildIds:          guildIds,
		PlatformAnchorIds: platformAnchorIds,
	}
}

func ApplyCMSAnchorRoomVisibility(m *gdb.Model, roomIDColumn string, scope CMSAnchorRoomVisibilityScope) *gdb.Model {
	return applyCMSAnchorRoomVisibility(m, roomIDColumn, scope)
}

func applyCMSAnchorRoomVisibility(m *gdb.Model, roomIDColumn string, scope CMSAnchorRoomVisibilityScope) *gdb.Model {
	if m == nil || !scope.Apply {
		return m
	}
	guildCol := "r." + string(entity.LiveRoomGuildId)
	roomCol := "r.id"
	m = m.InnerJoin(string(entity.TbLiveRoom)+" r", roomCol+" = "+roomIDColumn)
	if len(scope.GuildIds) > 0 && len(scope.PlatformAnchorIds) > 0 {
		return m.Where("("+guildCol+" IN (?) OR ("+guildCol+" = 0 AND "+roomCol+" IN (?)))", scope.GuildIds, scope.PlatformAnchorIds)
	}
	if len(scope.GuildIds) > 0 {
		return m.WhereIn(guildCol, scope.GuildIds)
	}
	return m.Where(guildCol+" = ?", 0).WhereIn(roomCol, scope.PlatformAnchorIds)
}

package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbLiveGuildAnchor db.TbName = "live_guild_anchors"
)

const (
	GuildAnchorGuildId db.TbCol = "guild_id"
	GuildAnchorRoomId  db.TbCol = "room_id"
)

// LiveGuildAnchor 工会名下主播列表(直播间维度,直接读写数据库)
type LiveGuildAnchor struct {
	migrate.OneModel
	GuildId uint64 `gorm:"index;uniqueIndex:uk_guild_room;default:0;comment:工会ID" json:"guildId"`
	RoomId  uint64 `gorm:"index;uniqueIndex:uk_guild_room;default:0;comment:直播间ID(live_rooms.id)" json:"roomId"`
}

func initGuildAnchor() {
	migrate.AutoMigrate(&LiveGuildAnchor{})
}

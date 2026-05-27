package liveroom

import (
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity"
)

// EnsureAnchorRoom 确保主播拥有直播间记录(CMS设为主播时预创建,App端后续可完善资料)
func EnsureAnchorRoom(anchorId, guildId uint64) *entity.LiveRoom {
	if room := liveroomdao.GetRoomByAnchor(anchorId); room != nil {
		return room
	}
	room := entity.NewLiveRoom(anchorId, guildId, "", "", "")
	liveroomdao.AddRoomToCache(room)
	return room
}
